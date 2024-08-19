package siakad

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func PanduanDosen(message itmodel.IteungMessage, db *mongo.Database) string {
	// Ekstraksi pesan yang diterima dari message
	pesan := message.Message

	// Buat filter regex untuk mencari dokumen dengan prompt yang mengandung kata kunci
	filter := bson.M{"prompt": bson.M{"$regex": pesan, "$options": "i"}}

	var prompt Prompt
	err := db.Collection("panduansiakad").FindOne(context.TODO(), filter).Decode(&prompt)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Jika tidak ditemukan dokumen yang cocok, berikan respon berikut
			return "Keyword salah. Berikut keyword yang benar: myika panduan dosen"
		}
		// Jika ada kesalahan lain, tampilkan pesan kesalahan
		log.Printf("Error finding document in MongoDB: %v", err)
		return fmt.Sprintf("Maaf, terjadi kesalahan saat mengambil panduan dosen: %v", err)
	}

	// Mengembalikan jawaban yang sesuai dengan prompt
	return prompt.Answer
}

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

func extractProdi(message string) string {
	re := regexp.MustCompile(`(?i)prodi:\s*([a-zA-Z0-9]+)`)
	matches := re.FindStringSubmatch(message)
	if len(matches) > 1 {
		return matches[1]
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

	prodi := extractProdi(message.Message)
	if prodi == "" {
		return "Prodinya di sertakan dulu dong kak " + message.Alias_name + " di akhir pesan nya dengan format 'prodi: [prodi]'"
	}

	var conf Config
	err := db.Collection("config").FindOne(context.TODO(), bson.M{"phonenumber": "62895601060000"}).Decode(&conf)
	if err != nil {
		return "Wah kak " + message.Alias_name + " mohon maaf ada kesalahan dalam pengambilan config di database " + err.Error()
	}

	fmt.Println("SiakadLoginURL:", conf.SiakadLoginURL)

	if conf.SiakadLoginURL == "" {
		return "URL untuk login tidak ditemukan dalam konfigurasi."
	}

	loginRequest := LoginRequest{
		Email:    email,
		Password: password,
		Role:     role,
		Prodi:    prodi,
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

	if resp.StatusCode == http.StatusInternalServerError {
		return "Gagal login, email atau password salah. \n_Notes : jika kamu mengakses Siakad dengan SSO Google, harap cek/ganti password terlebih dahulu dan coba lagi login di domyikado._"
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("Gagal login, status code: %d", resp.StatusCode)
	}

	noHp := message.Phone_number
	if noHp == "" {
		return "Nomor telepon tidak ditemukan dalam pesan."
	}

	loginInfo := bson.M{
		"$set": bson.M{
			"nohp":       noHp,
			"email":      email,
			"role":       role,
			"prodi":      prodi,
			"login_time": time.Now(),
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err = db.Collection("siakad").UpdateOne(context.TODO(), bson.M{"email": email}, loginInfo, opts)
	if err != nil {
		return "Berhasil login, tetapi terjadi kesalahan saat menyimpan informasi login: " + err.Error()
	}

	return "Hai kak, " + message.Alias_name + "\nBerhasil login dengan email: " + email
}

func ApproveBAP(message itmodel.IteungMessage, db *mongo.Database) string {
	email := extractEmail(message.Message)
	if email == "user@email.com" {
		return "Emailnya di ubah dulu dong kak, jadi emailnya dosen yang ingin diapprove BAP "
	} else if email == "" {
		return "Emailnya di sertakan dulu dong kak " + message.Alias_name + " di akhir pesan nya"
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
		"email_dosen": email,
	})
	if err != nil {
		return "Gagal membuat request body: " + err.Error()
	}
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", conf.ApproveBapURL, bytes.NewBuffer(requestBody))
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

	// Check the response status code
	if resp.StatusCode == http.StatusForbidden {
		return "Kamu bukan Kaprodi ya! Silahkan hubungi kaprodi untuk approve BAP"
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("Gagal approve bap, status code: %d", resp.StatusCode)
	}

	return "Terima kasih pak, BAP Dosen dengan email " + email + " berhasil di approve, hubungi dosen terkait untuk cetak BAP nya di domykado dengan format pesan berikut: \n*cetak bap periode [periode]*\n\n*_Contoh Pesan:_*\n\n*_cetak bap periode 20232_*"
}

func CekApprovalBAP(message itmodel.IteungMessage, db *mongo.Database) string {
	// Ambil nomor telepon dari pesan
	noHp := message.Phone_number
	if noHp == "" {
		return "Nomor telepon tidak ditemukan dalam pesan."
	}

	// Ambil informasi login berdasarkan nomor telepon dari koleksi siakad
	var loginInfo struct {
		Email string `bson:"email"`
		Role  string `bson:"role"`
	}
	err := db.Collection("siakad").FindOne(context.TODO(), bson.M{"nohp": noHp}).Decode(&loginInfo)
	if err != nil {
		return "Nomor telepon tidak ditemukan, silahkan login dengan klik link ini: https://wa.me/628999710040?text=login%20siakad%20email%3A%20email%20password%3A%20password%20role%3A%20dosen%20prodi%3A%20D4%20TI"
	}

	// Cek apakah role pengguna adalah dosen
	if loginInfo.Role != "dosen" {
		return "Akses ini hanya tersedia untuk dosen. Mohon maaf jika Anda bukan dosen."
	}

	// Ambil URL API dari database
	var conf Config
	err = db.Collection("config").FindOne(context.TODO(), bson.M{"phonenumber": "62895601060000"}).Decode(&conf)
	if err != nil {
		return "Wah Bapak/Ibu " + message.Alias_name + ", mohon maaf ada kesalahan dalam pengambilan config di database: " + err.Error()
	}

	// Buat HTTP client baru dengan timeout
	client := &http.Client{Timeout: 10 * time.Second}

	// Buat POST request tanpa body, hanya header
	req, err := http.NewRequest("POST", conf.CekApprovalBapURL, nil)
	if err != nil {
		return "Gagal membuat request: " + err.Error()
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("nohp", noHp)

	// Kirim request
	resp, err := client.Do(req)
	if err != nil {
		return "Gagal mengirim request: " + err.Error()
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "Akun tidak ditemukan! Silakan klik link ini: https://wa.me/628999710040?text=login%20siakad%20email%3A%20email%20password%3A%20password%20role%3A%20dosen%20prodi%3A%20D4TI"
	}

	// Periksa status kode dari respon
	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("Gagal cek approval BAP, status code: %d", resp.StatusCode)
	}

	// Parse body respon sebagai string
	var approvalStatus string
	err = json.NewDecoder(resp.Body).Decode(&approvalStatus)
	if err != nil {
		return "Gagal memproses respon dari server: " + err.Error()
	}

	// Periksa status approval dan kembalikan pesan yang sesuai
	if approvalStatus == "true" {
		return "BAP sudah di Approve! Gunakan format pesan berikut: \n*cetak bap periode [periode]*\n\n*_Contoh Pesan:_*\n\n*_cetak bap periode 20232_*"
	} else {
		// Buat URL WhatsApp dengan email
		whatsappURL := fmt.Sprintf("https://wa.me/628999710040?text=approve%%20bap%%20email:%%20%s", loginInfo.Email)
		return fmt.Sprintf("BAP belum diapprove! Silakan hubungi kaprodi untuk approve BAP dengan kirimkan url ini: %s", whatsappURL)
	}
}

func extractPeriod(message string) string {
	// Function to extract class and period from the message
	var periode string
	fmt.Sscanf(message, "cetak bap periode %s", &periode)
	return periode
}

// CetakBAP processes the request for BAP
func CetakBAP(message itmodel.IteungMessage, db *mongo.Database) string {
	// Extract information from the message
	periode := extractPeriod(message.Message)
	if periode == "" {
		return "Pesan tidak sesuai format. Gunakan format 'cetak bap periode [periode]'"
	}

	// Get the phone number from the message
	noHp := message.Phone_number
	if noHp == "" {
		return "Nomor telepon tidak ditemukan dalam pesan."
	}

	// Ambil informasi login berdasarkan nomor telepon dari koleksi siakad
	var loginInfo struct {
		Email string `bson:"email"`
		Role  string `bson:"role"`
	}
	err := db.Collection("siakad").FindOne(context.TODO(), bson.M{"nohp": noHp}).Decode(&loginInfo)
	if err != nil {
		return "Nomor telepon tidak ditemukan, silahkan login dengan klik link ini: https://wa.me/628999710040?text=login%20siakad%20email%3A%20email%20password%3A%20password%20role%3A%20dosen%20prodi%3A%20D4TI"
	}

	// Cek apakah role pengguna adalah dosen
	if loginInfo.Role != "dosen" {
		return "Akses ini hanya tersedia untuk dosen. Mohon maaf jika Anda bukan dosen."
	}

	// Ambil URL API dari database
	var conf Config
	err = db.Collection("config").FindOne(context.TODO(), bson.M{"phonenumber": "62895601060000"}).Decode(&conf)
	if err != nil {
		return "Wah Bapak/Ibu " + message.Alias_name + ", mohon maaf ada kesalahan dalam pengambilan config di database: " + err.Error()
	}

	// Prepare the request body
	requestBody, err := json.Marshal(map[string]string{
		"periode": periode,
	})
	if err != nil {
		return "Gagal membuat request body: " + err.Error()
	}

	// Create and send the HTTP request
	client := &http.Client{Timeout: 540 * time.Second}
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

	email := loginInfo.Email
	if resp.StatusCode == http.StatusForbidden {
		whatsappURL := fmt.Sprintf("https://wa.me/628999710040?text=approve%%20bap%%20email:%%20%s", email)
		return fmt.Sprintf("Gagal, BAP belum diapprove! Silakan hubungi kaprodi untuk approve BAP dengan kirimkan url ini: %s", whatsappURL)
	}

	if resp.StatusCode == http.StatusNotFound {
		return "Akun tidak ditemukan! silahkan klik link ini https://wa.me/628999710040?text=login%20siakad%20email%3A%20email%20password%3A%20password%20role%3A%20dosen%20prodi%3A%20D4TI"
	}

	if resp.StatusCode != http.StatusOK {
		return "Gagal mendapatkan BAP, kamu bukan dosen."
	}

	var responseMap []map[string]string
	err = json.NewDecoder(resp.Body).Decode(&responseMap)
	if err != nil {
		return "Gagal memproses response: " + err.Error()
	}

	// Format the response message
	responseMessage := "Berikut adalah BAP Bapak/Ibu:\n"
	for _, item := range responseMap {
		k := item["kelas"]
		u := item["url"]
		if k != "" && u != "" {
			responseMessage += fmt.Sprintf("Kelas %s: %s\n", k, u)
		}
	}

	return responseMessage
}

func extractNimandTopik(message string) (string, string) {
	var nim, topik string
	// Handle non-breaking spaces
	message = strings.ReplaceAll(message, "\u00A0", " ")

	// Regex patterns to extract NIM and topik
	nimPattern := regexp.MustCompile(`(?i)nim\s+(\d+)`)
	topikPattern := regexp.MustCompile(`(?i)topik\s+(.+)`)

	// Find matches in the message
	nimMatch := nimPattern.FindStringSubmatch(message)
	topikMatch := topikPattern.FindStringSubmatch(message)

	// Extract NIM
	if len(nimMatch) > 1 {
		nim = nimMatch[1]
	}

	// Extract Topik
	if len(topikMatch) > 1 {
		topik = strings.TrimSpace(topikMatch[1])
		// Remove the word "poin" from topik if it exists
		topik = strings.ReplaceAll(topik, "poin", "")
		topik = strings.TrimSpace(topik)
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

	// Ambil informasi login berdasarkan nomor telepon dari koleksi siakad
	var loginInfo struct {
		Email string `bson:"email"`
		Role  string `bson:"role"`
	}
	err := db.Collection("siakad").FindOne(context.TODO(), bson.M{"nohp": noHp}).Decode(&loginInfo)
	if err != nil {
		return "Nomor telepon tidak ditemukan, silahkan login dengan klik link ini: https://wa.me/628999710040?text=login%20siakad%20email%3A%20email%20password%3A%20password%20role%3A%20dosen%20prodi%3A%20D4TI"
	}

	// Cek apakah role pengguna adalah dosen
	if loginInfo.Role != "dosen" {
		return "Akses ini hanya tersedia untuk dosen. Mohon maaf jika Anda bukan dosen."
	}

	// Ambil URL API dari database
	var conf Config
	err = db.Collection("config").FindOne(context.TODO(), bson.M{"phonenumber": "62895601060000"}).Decode(&conf)
	if err != nil {
		return "Wah Bapak/Ibu " + message.Alias_name + ", mohon maaf ada kesalahan dalam pengambilan config di database: " + err.Error()
	}

	// Prepare the request body
	requestBody, err := json.Marshal(map[string]string{
		"nim":   nim,
		"topik": topik,
	})
	if err != nil {
		fmt.Printf("Error creating request body: %s\n", err.Error())
		return "Gagal membuat request body: " + err.Error()
	}

	// Create and send the HTTP request
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", conf.ApproveBimbinganURL, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Printf("Error creating HTTP request: %s\n", err.Error())
		return "Gagal membuat request: " + err.Error()
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("nohp", noHp)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending HTTP request: %s\n", err.Error())
		return "Gagal mengirim request: " + err.Error()
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse map[string]string
		_ = json.NewDecoder(resp.Body).Decode(&errorResponse)
		switch resp.StatusCode {
		case http.StatusNotFound:
			return "Token tidak ditemukan! klik link ini https://wa.me/628999710040?text=login%20siakad%20email%3A%20email%20password%3A%20password%20role%3A%20mahasiswa%20prodi%3A%20D4TI"
		case http.StatusForbidden:
			return "Gagal, Bimbingan telah disetujui!"
		default:
			return fmt.Sprintf("Gagal approve bimbingan, status code: %d, error: %s", resp.StatusCode, errorResponse["error"])
		}
	}

	var responseMap map[string]string
	err = json.NewDecoder(resp.Body).Decode(&responseMap)
	if err != nil {
		fmt.Printf("Error decoding response: %s\n", err.Error())
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

	// Ambil informasi login berdasarkan nomor telepon dari koleksi siakad
	var loginInfo struct {
		Email string `bson:"email"`
		Role  string `bson:"role"`
	}
	err := db.Collection("siakad").FindOne(context.TODO(), bson.M{"nohp": noHp}).Decode(&loginInfo)
	if err != nil {
		return "Nomor telepon tidak ditemukan, silahkan login dengan klik link ini: https://wa.me/628999710040?text=login%20siakad%20email%3A%20email%20password%3A%20password%20role%3A%20dosen%20prodi%3A%20D4TI"
	}

	// Cek apakah role pengguna adalah dosen
	if loginInfo.Role != "dosen" {
		return "Akses ini hanya tersedia untuk dosen. Mohon maaf jika Anda bukan dosen."
	}

	// Ambil URL API dari database
	var conf Config
	err = db.Collection("config").FindOne(context.TODO(), bson.M{"phonenumber": "62895601060000"}).Decode(&conf)
	if err != nil {
		return "Wah Bapak/Ibu " + message.Alias_name + ", mohon maaf ada kesalahan dalam pengambilan config di database: " + err.Error()
	}

	// Prepare the request body
	requestBody, err := json.Marshal(map[string]string{
		"nim":   nim,
		"topik": topik,
	})
	if err != nil {
		fmt.Printf("Error creating request body: %s\n", err.Error())
		return "Gagal membuat request body: " + err.Error()
	}

	// Create and send the HTTP request
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", conf.ApproveBimbinganByPoinURL, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Printf("Error creating HTTP request: %s\n", err.Error())
		return "Gagal membuat request: " + err.Error()
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("nohp", noHp)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending HTTP request: %s\n", err.Error())
		return "Gagal mengirim request: " + err.Error()
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse map[string]string
		_ = json.NewDecoder(resp.Body).Decode(&errorResponse)
		switch resp.StatusCode {
		case http.StatusNotFound:
			return "Token tidak ditemukan! klik link ini https://wa.me/628999710040?text=login%20siakad%20email%3A%20email%20password%3A%20password%20role%3A%20mahasiswa%20prodi%3A%20D4TI"
		case http.StatusForbidden:
			return "Gagal, Bimbingan telah disetujui!"
		default:
			return fmt.Sprintf("Gagal approve bimbingan, status code: %d, error: %s", resp.StatusCode, errorResponse["error"])
		}
	}

	var responseMap map[string]string
	err = json.NewDecoder(resp.Body).Decode(&responseMap)
	if err != nil {
		fmt.Printf("Error decoding response: %s\n", err.Error())
		return "Gagal memproses response: " + err.Error()
	}

	return fmt.Sprintf("Bimbingan berhasil di approve! %s", responseMap["poin_mahasiswa"])
}
