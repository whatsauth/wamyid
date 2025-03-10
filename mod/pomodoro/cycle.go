package pomodoro

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocroot/helper/atdb"
	"github.com/whatsauth/itmodel"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func HandlePomodoroReport(Profile itmodel.Profile, Pesan itmodel.IteungMessage, db *mongo.Database) string {
	// 1. Validasi input dasar
	if Pesan.Message == "" {
		return "Wah kak " + Pesan.Alias_name + ", pesan tidak boleh kosong"
	}

	// 2. Ekstrak data dari pesan
	cycle := extractCycleNumber(Pesan.Message)
	if cycle == 0 {
		return "Wah kak " + Pesan.Alias_name + ", format cycle tidak valid. Contoh: 'Iteung Pomodoro Report 1 cycle'"
	}

	hostname := extractValue(Pesan.Message, "Hostname : ")
	// Perbaikan: Pastikan hostname tidak menyertakan "IP" 
	if strings.HasSuffix(hostname, "IP") {
		hostname = strings.TrimSuffix(hostname, "IP")
	}
	ip := extractIP(Pesan.Message) // Gunakan fungsi khusus IP
	screenshots := extractNumber(Pesan.Message, "Jumlah ScreenShoot : ")
	pekerjaan := extractActivities(Pesan.Message) // Update parameter
	token := extractToken(Pesan.Message)

	// 3. Verifikasi public key
	publicKey, err := getPublicKey(db)
	if err != nil {
		return "Wah kak " + Pesan.Alias_name + ", sistem gagal memuat public key: " + err.Error()
	}

	// Cek apakah token sudah pernah digunakan
	isUsed, err := isTokenUsed(db, token)
	if err != nil {
		return "Wah kak " + Pesan.Alias_name + ", sistem gagal memeriksa token: " + err.Error()
	}

	if isUsed {
		return "Wah kak " + Pesan.Alias_name + ", token ini sudah pernah digunakan sebelumnya"
	}

	// 4. Decode token
	decode, err := watoken.Decode(publicKey, token)
	if err != nil {
		errorMsg := "Token tidak valid"

		// Deteksi jenis error
		if strings.Contains(err.Error(), "expired") {
			errorMsg = "Token sudah kedaluwarsa"
		} else if strings.Contains(err.Error(), "invalid") {
			errorMsg = "Format token tidak valid"
		} else if strings.Contains(err.Error(), "hex") {
			errorMsg = "Format public key tidak valid"
		}

		return fmt.Sprintf("Wah kak %s, %s: %v",
			Pesan.Alias_name,
			errorMsg,
			strings.Split(err.Error(), ":")[0], // Ambil pesan error utama
		)
	}

	// Tandai token sebagai telah digunakan
	err = markTokenAsUsed(db, token, Pesan.Phone_number)
	if err != nil {
		return "Wah kak " + Pesan.Alias_name + ", sistem gagal menandai token: " + err.Error()
	}

	// 5. Validasi payload dan ekstrak URL
	var url string
	payloadStr := fmt.Sprintf("%v", decode)
	// Ekstrak URL dari string
	urlRegex := regexp.MustCompile(`\{(https://[^\s]+)`)
	urlMatch := urlRegex.FindStringSubmatch(payloadStr)
	if len(urlMatch) > 1 {
		url = urlMatch[1]
	}

	// 6. Simpan ke database
	loc, _ := time.LoadLocation("Asia/Jakarta")
	report := PomodoroReport{
		PhoneNumber: Pesan.Phone_number,
		Cycle:       cycle,
		Hostname:    hostname,
		IP:          ip,
		Screenshots: screenshots,
		Pekerjaan:   pekerjaan,
		Token:       token,
		CreatedAt:   time.Now().In(loc),
	}

	_, err = atdb.InsertOneDoc(db, "pomokit", report)
	if err != nil {
		return "Wah kak " + Pesan.Alias_name + ", gagal menyimpan laporan: " + err.Error()
	}

	// 7. Generate response
	return fmt.Sprintf(
		"✅ *Laporan Cycle %d Berhasil!*\n"+
			"Nama: %s\n"+
			"Hostname: %s\n"+
			"IP: %s\n"+
			"Aktivitas:\n- %s\n"+
			"🔗 Alamat URL %s\n"+
			"🕒 %s",
		cycle,
		Pesan.Alias_name,
		hostname,
		ip,
		strings.Join(pekerjaan, "\n- "),
		url, // Tampilkan URL dari payload
		time.Now().Format("2006-01-02 15:04"),
	)
}

// Fungsi untuk memeriksa apakah token sudah digunakan
func isTokenUsed(db *mongo.Database, token string) (bool, error) {
	count, err := db.Collection("used_tokens").CountDocuments(context.Background(), bson.M{"token": token})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Fungsi untuk menandai token sebagai telah digunakan
func markTokenAsUsed(db *mongo.Database, token string, phoneNumber string) error {
	tokenData := bson.M{
		"token":        token,
		"phone_number": phoneNumber,
		"used_at":      time.Now(),
	}

	_, err := db.Collection("used_tokens").InsertOne(context.Background(), tokenData)
	return err
}

// Helper functions
func extractCycleNumber(msg string) int {
	re := regexp.MustCompile(`Report\s+(\d+)\s+cycle`)
	matches := re.FindStringSubmatch(msg)
	if len(matches) > 1 {
		cycle, _ := strconv.Atoi(matches[1])
		return cycle
	}
	return 0
}

func extractValue(msg, prefix string) string {
	re := regexp.MustCompile(regexp.QuoteMeta(prefix) + `(\S+)(?:\s|$)`)
	match := re.FindStringSubmatch(msg)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return ""
}

func extractIP(msg string) string {
	// Coba pola IP langsung
	re := regexp.MustCompile(`IP\s*:\s*(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)
	match := re.FindStringSubmatch(msg)
	if len(match) > 1 {
		return match[1]
	}

	// Jika tidak ditemukan, coba ekstrak dari URL
	reURL := regexp.MustCompile(`IP\s*:\s*https://whatismyipaddress\.com/ip/(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)
	matchURL := reURL.FindStringSubmatch(msg)
	if len(matchURL) > 1 {
		return matchURL[1]
	}

	return ""
}

func extractNumber(msg, prefix string) int {
	re := regexp.MustCompile(regexp.QuoteMeta(prefix) + `(\d+)`)
	match := re.FindStringSubmatch(msg)
	if len(match) > 1 {
		num, _ := strconv.Atoi(match[1])
		return num
	}
	return 0
}

func extractActivities(msg string) []string {
	// Perbaikan regex untuk menangkap semua konten setelah "Yang Dikerjakan :"
	// hingga tanda # atau akhir pesan, termasuk konten dengan format |text
	re := regexp.MustCompile(`Yang Dikerjakan\s*:\s*\n?([\s\S]+?)(?:\n\s*\#|$)`)
	match := re.FindStringSubmatch(msg)
	if len(match) > 1 {
		// Bersihkan teks
		text := strings.TrimSpace(match[1])
		
		// Handle format dengan pipe (|)
		if strings.HasPrefix(text, "|") {
			text = strings.TrimPrefix(text, "|")
			return []string{strings.TrimSpace(text)}
		}
		
		// Handle format dengan baris baru
		if strings.Contains(text, "\n") {
			lines := strings.Split(text, "\n")
			var activities []string
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" {
					// Handle line dengan prefix - atau |
					if strings.HasPrefix(line, "- ") {
						activities = append(activities, strings.TrimPrefix(line, "- "))
					} else if strings.HasPrefix(line, "|") {
						activities = append(activities, strings.TrimPrefix(line, "|"))
					} else {
						activities = append(activities, line)
					}
				}
			}
			return activities
		}
		
		return []string{text}
	}
	
	// Jika tidak match dengan pattern di atas, coba pattern alternatif
	altRe := regexp.MustCompile(`Yang Dikerjakan\s*:(.+?)(?:\n\s*\#|$)`)
	altMatch := altRe.FindStringSubmatch(msg)
	if len(altMatch) > 1 {
		text := strings.TrimSpace(altMatch[1])
		if text != "" {
			return []string{text}
		}
	}
	
	return []string{"Tidak ada detail aktivitas"}
}

func extractToken(msg string) string {
	re := regexp.MustCompile(`#(v4\..+)`)
	match := re.FindStringSubmatch(msg)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

// func extractURLFromPayload(payload any) string {
// 	// Cek jika payload adalah string
// 	if urlStr, ok := payload.(string); ok && strings.HasPrefix(urlStr, "http") {
// 		return urlStr
// 	}

// 	// Jika payload adalah map
// 	if payloadMap, ok := payload.(map[string]interface{}); ok {
// 		// Coba cari key yang berisi URL
// 		for _, v := range payloadMap { // Hapus variabel k yang tidak digunakan
// 			if urlStr, ok := v.(string); ok && strings.HasPrefix(urlStr, "http") {
// 				return urlStr
// 			}
// 		}
// 	}

// 	// Cek jika payload adalah struct
// 	payloadStr := fmt.Sprintf("%v", payload)
// 	// Ekstrak URL dari string representasi payload
// 	re := regexp.MustCompile(`\{(https://[^\s]+)`)
// 	match := re.FindStringSubmatch(payloadStr)
// 	if len(match) > 1 {
// 		return match[1]
// 	}

// 	return ""
// }

func getPublicKey(db *mongo.Database) (string, error) {
	conf, err := atdb.GetOneDoc[Config](db, "config", bson.M{"publickeypomokit": bson.M{"$exists": true}})
	if err != nil {
		return "", fmt.Errorf("konfigurasi tidak ditemukan")
	}
	return conf.PublicKeyPomokit, nil
}