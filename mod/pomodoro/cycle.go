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
	ip := extractIP(Pesan.Message) // Gunakan fungsi khusus IP
	screenshots := extractNumber(Pesan.Message, "Jumlah ScreenShoot : ")
	pekerjaan := extractActivities(Pesan.Message) // Update parameter
	token := extractToken(Pesan.Message)

	// 3. Verifikasi public key
	publicKey, err := getPublicKey(db)
	if err != nil {
		return "Wah kak " + Pesan.Alias_name + ", sistem gagal memuat public key: " + err.Error()
	}

	// 4. Decode token
	payload, err := watoken.Decode(publicKey, token)
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
		pekerjaan,
	)

	// Create a comparable payload string
	payloadStr := fmt.Sprintf(
		"cycle:%d|hostname:%s|ip:%s|screenshots:%d|activities:%v",
		cycle,
		hostname,
		ip,
		screenshots,
		pekerjaan,
	)

	// Compare the expected payload with what we constructed
	if payloadStr != expectedPayload {
		return "Wah kak " + Pesan.Alias_name + ", data laporan tidak sesuai dengan token"
	}

	if hostname == "" || ip == "" || len(pekerjaan) == 0 {
		return "Wah kak " + Pesan.Alias_name + ", Format laporan tidak valid!"
	}

	// 6. Simpan ke database
	report := PomodoroReport{
		PhoneNumber: Pesan.Phone_number,
		Cycle:       cycle,
		Hostname:    hostname,
		IP:          ip,
		Screenshots: screenshots,
		Pekerjaan:   pekerjaan,
		Token:   token,
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
		strings.Join(pekerjaan, "\n- "),
		url, // Tampilkan URL dari payload
		time.Now().Format("2006-01-02 15:04"),
	)
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
    re := regexp.MustCompile(`IP\s*:\s*(\S+)`)
    match := re.FindStringSubmatch(msg)
    if len(match) > 1 {
        return match[1]
    }
    return ""
}

func extractNumber(msg, prefix string) int {
	valStr := extractValue(msg, prefix)
	num, _ := strconv.Atoi(valStr)
	return num
}

func extractActivities(msg string) []string {
    re := regexp.MustCompile(`Yang Dikerjakan\s*:\s*\|([^#]+)`)
    match := re.FindStringSubmatch(msg)
    if len(match) > 1 {
        return strings.Split(match[1], "|")
    }
    return []string{}
}

func extractToken(msg string) string {
    re := regexp.MustCompile(`#(v4\..+)`)
    match := re.FindStringSubmatch(msg)
    if len(match) > 1 {
        return match[1]
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