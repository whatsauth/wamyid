package siakad

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func extractEmail(message string) string {
	re := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	return re.FindString(message)

}

func extractPassword(message string) string {
	re := regexp.MustCompile(`password: (\S+)`)
	matches := re.FindStringSubmatch(message)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func extractRole(message string) string {
	if strings.Contains(strings.ToLower(message), "dosen") {
		return "dosen"
	} else if strings.Contains(strings.ToLower(message), "mhs") || strings.Contains(strings.ToLower(message), "mahasiswa") {
		return "mhs"
	}
	return ""
}

func LoginSiakad(message itmodel.IteungMessage, db *mongo.Database) string {
	email := extractEmail(message.Message)
	if email == "user@email.com" {
		return "Emailnya di ubah dulu dong kak, jadi emailnya kak " + message.Alias_name
	} else if email == "" {
		return "Emailnya di sertakan dulu dong kak " + message.Alias_name + " di akhir pesan nya"
	}

	password := extractPassword(message.Message)
	if password == "" {
		return "Passwordnya di sertakan dulu dong kak " + message.Alias_name + " di akhir pesan nya dengan format 'password: [password]'"
	}

	role := extractRole(message.Message)
	if role == "" {
		return "Rolenya di sertakan dulu dong kak " + message.Alias_name + " di akhir pesan nya dengan format 'role: [dosen/mhs]'"
	}

	var conf Config
	err := db.Collection("config").FindOne(context.TODO(), bson.M{"phonenumber": "62895601060000"}).Decode(&conf)
	if err != nil {
		return "Wah kak " + message.Alias_name + " mohon maaf ada kesalahan dalam pengambilan config di database " + err.Error()
	}

	// Logging the SiakadLoginURL
	fmt.Println("SiakadLoginURL:", conf.SiakadLoginURL)

	if conf.SiakadLoginURL == "" {
		return "URL untuk login tidak ditemukan dalam konfigurasi."
	}

	loginRequest := LoginRequest{
		Email:    email,
		Password: password,
		Role:     role,
	}

	loginRequestBody, err := json.Marshal(loginRequest)
	if err != nil {
		return "Gagal membuat request body: " + err.Error()
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", conf.SiakadLoginURL, bytes.NewBuffer(loginRequestBody))
	if err != nil {
		return "Gagal membuat request: " + err.Error()
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "Gagal mengirim request: " + err.Error()
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("Gagal login, status code: %d", resp.StatusCode)
	}

	return "Hai kak, " + message.Alias_name + "\nBerhasil login dengan email:" + email
}

func extractClassAndPeriod(message string) (string, string) {
	// Function to extract class and period from the message
	var kelas, periode string
	fmt.Sscanf(message, "minta bap kelas %s periode %s", &kelas, &periode)
	return kelas, periode
}

func MintaBAP(message itmodel.IteungMessage, db *mongo.Database) string {
	// Extract information from the message
	kelas, periode := extractClassAndPeriod(message.Message)
	if kelas == "" || periode == "" {
		return "Pesan tidak sesuai format. Gunakan format 'minta bap kelas [kelas] periode [periode]'"
	}

	// Get the phone number from the message
	noHp := message.Phone_number
	if noHp == "" {
		return "Nomor telepon tidak ditemukan dalam pesan."
	}

	// Get the API URL from the database
	var conf Config
	err := db.Collection("config").FindOne(context.TODO(), bson.M{"phonenumber": "62895601060000"}).Decode(&conf)
	if err != nil {
		return "Wah kak " + message.Alias_name + " mohon maaf ada kesalahan dalam pengambilan config di database " + err.Error()
	}

	// Prepare the request body
	requestBody, err := json.Marshal(map[string]string{
		"periode": periode,
		"kelas":   kelas,
	})
	if err != nil {
		return "Gagal membuat request body: " + err.Error()
	}

	// Create and send the HTTP request
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", conf.BapURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "Gagal membuat request: " + err.Error()
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("nohp", noHp)

	resp, err := client.Do(req)
	if err != nil {
		return "Gagal mengirim request: " + err.Error()
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("Gagal mendapatkan BAP, status code: %d", resp.StatusCode)
	}

	var responseMap map[string]string
	err = json.NewDecoder(resp.Body).Decode(&responseMap)
	if err != nil {
		return "Gagal memproses response: " + err.Error()
	}

	return "Berikut adalah URL BAP yang diminta: " + responseMap["url"]
}

func extractNimandTopik(message string) (string, string) {
	var nim, topik string
	nimPattern := regexp.MustCompile(`(?i)nim\s+(\d+)`)
	topikPattern := regexp.MustCompile(`(?i)topik\s+(.+?)(?:\s+poin|$)`)

	nimMatch := nimPattern.FindStringSubmatch(message)
	topikMatch := topikPattern.FindStringSubmatch(message)

	if len(nimMatch) > 1 {
		nim = nimMatch[1]
	}
	if len(topikMatch) > 1 {
		topik = strings.TrimSpace(topikMatch[1])
	}

	fmt.Printf("Extracted NIM: %s, Topik: %s\n", nim, topik)
	return nim, topik
}

func ApproveBimbingan(message itmodel.IteungMessage, db *mongo.Database) string {
	// Extract information from the message
	nim, topik := extractNimandTopik(message.Message)
	if nim == "" || topik == "" {
		return "Pesan tidak sesuai format. Gunakan format 'approve bimbingan nim [nim] topik [topik]'"
	}

	// Get the phone number from the message
	noHp := message.Phone_number
	if noHp == "" {
		return "Nomor telepon tidak ditemukan dalam pesan."
	}

	// Get the API URL from the database
	var conf Config
	err := db.Collection("config").FindOne(context.TODO(), bson.M{"phonenumber": "62895601060000"}).Decode(&conf)
	if err != nil {
		return "Wah kak " + message.Alias_name + " mohon maaf ada kesalahan dalam pengambilan config di database: " + err.Error()
	}

	// Prepare the request body
	requestBody, err := json.Marshal(map[string]string{
		"nim":   nim,
		"topik": topik,
	})
	if err != nil {
		return "Gagal membuat request body: " + err.Error()
	}

	// Create and send the HTTP request
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", conf.ApproveBimbinganURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "Gagal membuat request: " + err.Error()
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("nohp", noHp)

	resp, err := client.Do(req)
	if err != nil {
		return "Gagal mengirim request: " + err.Error()
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("Gagal approve bimbingan, status code: %d", resp.StatusCode)
	}

	var responseMap map[string]string
	err = json.NewDecoder(resp.Body).Decode(&responseMap)
	if err != nil {
		return "Gagal memproses response: " + err.Error()
	}

	return responseMap["message"]
}

func ApproveBimbinganbyPoin(message itmodel.IteungMessage, db *mongo.Database) string {
	// Extract information from the message
	nim, topik := extractNimandTopik(message.Message)
	if nim == "" || topik == "" {
		return "Pesan tidak sesuai format. Gunakan format 'approve bimbingan nim [nim] topik [topik]'"
	}

	// Get the phone number from the message
	noHp := message.Phone_number
	if noHp == "" {
		return "Nomor telepon tidak ditemukan dalam pesan."
	}

	// Get the API URL from the database
	var conf Config
	err := db.Collection("config").FindOne(context.TODO(), bson.M{"phonenumber": "62895601060000"}).Decode(&conf)
	if err != nil {
		return "Wah kak " + message.Alias_name + " mohon maaf ada kesalahan dalam pengambilan config di database: " + err.Error()
	}

	// Prepare the request body
	requestBody, err := json.Marshal(map[string]string{
		"nim":   nim,
		"topik": topik,
	})
	if err != nil {
		return "Gagal membuat request body: " + err.Error()
	}

	// Create and send the HTTP request
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", conf.ApproveBimbinganByPoinURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "Gagal membuat request: " + err.Error()
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("nohp", noHp)

	resp, err := client.Do(req)
	if err != nil {
		return "Gagal mengirim request: " + err.Error()
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Status Code: %d, Response Body: %s\n", resp.StatusCode, string(body))

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("Gagal approve bimbingan, status code: %d", resp.StatusCode)
	}

	var responseMap map[string]string
	err = json.NewDecoder(bytes.NewBuffer(body)).Decode(&responseMap)
	if err != nil {
		return "Gagal memproses response: " + err.Error()
	}

	return fmt.Sprintf("Bimbingan berhasil di approve! Poin mahasiswa telah berkurang menjadi: %s", responseMap["poin_mahasiswa"])
}
