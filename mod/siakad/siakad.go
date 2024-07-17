package siakad

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

func RequestBAP(message itmodel.IteungMessage, db *mongo.Database) string {
	noHp := message.Phone_number
	if noHp == "" {
		return "Nomor HP tidak ditemukan."
	}

	// Pastikan data login berhasil sebelum memanggil fungsi ini
	loginResponse := LoginSiakad(message, db)
	if !strings.Contains(loginResponse, "Berhasil login") {
		return loginResponse
	}

	// Ambil periode dan kelas dari message
	var requestData struct {
		Periode string `json:"periode"`
		Kelas   string `json:"kelas"`
	}

	err := json.Unmarshal([]byte(message.Message), &requestData)
	if err != nil || requestData.Periode == "" || requestData.Kelas == "" {
		return "Periode atau kelas tidak ditemukan dalam pesan."
	}

	// Definisikan data request untuk BAP
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return "Gagal membuat request body untuk BAP: " + err.Error()
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", "conf.BapURL", bytes.NewBuffer(requestBody))
	if err != nil {
		return "Gagal membuat request untuk BAP: " + err.Error()
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("nohp", noHp)

	resp, err := client.Do(req)
	if err != nil {
		return "Gagal mengirim request untuk BAP: " + err.Error()
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("Gagal mendapatkan BAP, status code: %d", resp.StatusCode)
	}

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "Gagal membaca respons BAP: " + err.Error()
	}

	url, ok := response["url"]
	if !ok {
		return "Respons BAP tidak mengandung URL."
	}

	return "Berikut adalah URL BAP yang diminta: " + url
}
