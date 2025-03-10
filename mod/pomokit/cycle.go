package pomokit

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

	cycle := extractCycleNumber(Pesan.Message)
	if cycle == 0 {
		return "Wah kak " + Pesan.Alias_name + ", format cycle tidak valid. Contoh: 'Iteung Pomodoro Report 1 cycle'"
	}

	hostname := extractValue(Pesan.Message, "Hostname : ")
	// Perbaikan: Pastikan hostname tidak menyertakan "IP" 
	// hostname = strings.TrimSuffix(hostname, "IP")
	ip := extractIP(Pesan.Message) // Gunakan fungsi khusus IP
	screenshots := extractNumber(Pesan.Message, "Jumlah ScreenShoot : ")
	pekerjaan := extractActivities(Pesan.Message) // Update parameter
	token := extractToken(Pesan.Message)

	// 3. Verifikasi public key
	publicKey, err := getPublicKey(db)
	if err != nil {
		return "Wah kak " + Pesan.Alias_name + ", sistem gagal memuat public key: " + err.Error()
	}

	// Cek apakah token sudah pernah digunakan di koleksi pomokit
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
		URLPekerjaan: url,
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
			"Aktivitas: %s\n"+
			"🔗 Alamat URL %s\n"+
			"📅 %s",
		cycle,
		Pesan.Alias_name,
		hostname,
		ip,
		pekerjaan,
		url,
		report.CreatedAt.Format("2006-01-02 🕒15:04 WIB"), // ini dikonversi
	)
}

// Fungsi untuk memeriksa apakah token sudah digunakan menggunakan koleksi pomokit
func isTokenUsed(db *mongo.Database, token string) (bool, error) {
	// Menggunakan koleksi pomokit yang sudah ada untuk mengecek token
	count, err := db.Collection("pomokit").CountDocuments(context.Background(), bson.M{"token": token})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

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
    // 1. Cek apakah format URL whatismyipaddress sudah ada
    reURL := regexp.MustCompile(`IP\s*:\s*(https://whatismyipaddress\.com/ip/\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)
    matchURL := reURL.FindStringSubmatch(msg)
    if len(matchURL) > 1 {
        return matchURL[1] // Langsung kembalikan URL lengkap
    }

    // 2. Jika tidak ada URL, cari IP biasa dan konstruksi URL
    reIP := regexp.MustCompile(`IP\s*:\s*(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)
    matchIP := reIP.FindStringSubmatch(msg)
    if len(matchIP) > 1 {
        // Bangun URL dari IP yang ditemukan
        return "https://whatismyipaddress.com/ip/" + matchIP[1]
    }

    return "" // Jika tidak ada IP/URL yang valid
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

func extractActivities(msg string) string {
    // Regex untuk menangkap konten setelah "Yang Dikerjakan :" dan menghiraukan "|" di awal
    re := regexp.MustCompile(`Yang Dikerjakan\s*:\s*\n?\|?\s*([^#]+)`)
    match := re.FindStringSubmatch(msg)
    
    if len(match) > 1 {
        // Hilangkan karakter "|" di awal (jika ada) dan whitespace
        cleaned := strings.TrimLeft(match[1], "| ") // Hapus "|" dan spasi di awal
        cleaned = strings.TrimSpace(cleaned)        // Hilangkan spasi/newline di akhir
        return cleaned
    }
    
    return "Tidak ada detail aktivitas"
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

// HandlePomodoroStart menangani pesan permintaan untuk memulai siklus Pomodoro
func HandlePomodoroStart(Profile itmodel.Profile, Pesan itmodel.IteungMessage, db *mongo.Database) string {
	// Validasi input dasar
	if Pesan.Message == "" {
		return "Wah kak " + Pesan.Alias_name + ", pesan tidak boleh kosong"
	}

	// Ekstrak informasi dari pesan
	cycle := extractCycleNumber(Pesan.Message)
	if cycle == 0 {
		return "Wah kak " + Pesan.Alias_name + ", format cycle tidak valid. Contoh: 'Pomodoro Start 1 cycle'"
	}

	milestone := extractMilestone(Pesan.Message)
	version := extractVersion(Pesan.Message)
	hostname := extractValue(Pesan.Message, "Hostname : ")
	ip := extractIP(Pesan.Message)

	// Lokasi waktu Indonesia
	loc, _ := time.LoadLocation("Asia/Jakarta")
	currentTime := time.Now().In(loc)

	// Format respons
	return fmt.Sprintf(
		"🍅 *Pomodoro Cycle %d Dimulai!*\n"+
			"Nama: %s\n"+
			"Milestone: %s\n"+
			"Version: %s\n"+
			"Hostname: %s\n"+
			"IP: %s\n"+
			"📅 %s\n\n"+
			"Semangat kak! Waktu kerja nya dimulai 🚀",
		cycle,
		Pesan.Alias_name,
		milestone,
		version,
		hostname,
		ip,
		currentTime.Format("2006-01-02 🕒15:04 WIB"),
	)
}

// Fungsi tambahan untuk ekstraksi data dari pesan
func extractMilestone(msg string) string {
	re := regexp.MustCompile(`Milestone\s*:\s*([^$\n]+)`)
	match := re.FindStringSubmatch(msg)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return "Tidak ada milestone"
}

func extractVersion(msg string) string {
	re := regexp.MustCompile(`Version\s*:\s*([^\s\n]+)`)
	match := re.FindStringSubmatch(msg)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return "1.0.0" // Default version
}