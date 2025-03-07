package pomodoro

import (
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
	ip := extractValue(Pesan.Message, "IP : ")
	screenshots := extractNumber(Pesan.Message, "Jumlah ScreenShoot : ")
	activities := extractActivities(strings.Split(Pesan.Message, "\n"))
	signature := extractSignature(Pesan.Message)

	// 3. Verifikasi signature
	publicKey, err := getPublicKey(db)
	if err != nil {
		return "Wah kak " + Pesan.Alias_name + ", sistem gagal memuat public key: " + err.Error()
	}

	// 4. Verifikasi token dan payload
	payload, err := watoken.Decode(publicKey, signature)
	if err != nil {
		errorMsg := "Signature tidak valid"
		
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

	// 5. Validasi payload dan ekstrak URL
	var url string
	if payloadData, ok := payload.Data.(map[string]interface{}); ok {
		if u, exists := payloadData["url"].(string); exists {
			url = u
		}
	}

	expectedPayload := fmt.Sprintf(
		"cycle:%d|hostname:%s|ip:%s|screenshots:%d|activities:%v",
		cycle,
		hostname,
		ip,
		screenshots,
		activities,
	)

	// Create a comparable payload string
	payloadStr := fmt.Sprintf(
		"cycle:%d|hostname:%s|ip:%s|screenshots:%d|activities:%v",
		cycle,
		hostname,
		ip,
		screenshots,
		activities,
	)

	// Compare the expected payload with what we constructed
	if payloadStr != expectedPayload {
		return "Wah kak " + Pesan.Alias_name + ", data laporan tidak sesuai dengan signature"
	}

	// 6. Simpan ke database
	report := PomodoroReport{
		PhoneNumber: Pesan.Phone_number,
		Cycle:       cycle,
		Hostname:    hostname,
		IP:          ip,
		Screenshots: screenshots,
		Aktivitas:   activities,
		Signature:   signature,
		CreatedAt:   time.Now(),
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
		strings.Join(activities, "\n- "),
		url, // Tampilkan URL dari payload
		time.Now().Format("2006-01-02 15:04"),
	)
}

// Helper functions
func extractCycleNumber(msg string) int {
	re := regexp.MustCompile(`Report (\d+) cycle`)
	matches := re.FindStringSubmatch(msg)
	if len(matches) > 1 {
		cycle, _ := strconv.Atoi(matches[1])
		return cycle
	}
	return 0
}

func extractValue(msg, prefix string) string {
	for _, line := range strings.Split(msg, "\n") {
		if strings.Contains(line, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(line, prefix))
		}
	}
	return ""
}

func extractNumber(msg, prefix string) int {
	valStr := extractValue(msg, prefix)
	num, _ := strconv.Atoi(valStr)
	return num
}

func extractActivities(lines []string) []string {
	var activities []string
	for _, line := range lines {
		if strings.HasPrefix(line, "|") {
			activities = append(activities, strings.TrimPrefix(line, "| "))
		}
	}
	return activities
}

func extractSignature(msg string) string {
	parts := strings.Split(msg, "#")
	if len(parts) > 1 {
		return strings.TrimSpace(parts[len(parts)-1])
	}
	return ""
}

func getPublicKey(db *mongo.Database) (string, error) {
    conf, err := atdb.GetOneDoc[Config](db, "config", bson.M{"publickeypomokit": bson.M{"$exists": true}})
    if err != nil {
        return "", fmt.Errorf("konfigurasi tidak ditemukan")
    }
    return conf.PublicKeyPomokit, nil
}