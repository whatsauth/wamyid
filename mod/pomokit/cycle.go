package pomokit

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocroot/mod/daftar"
	"github.com/gocroot/helper/atapi"
	"github.com/gocroot/helper/atdb"
	"github.com/whatsauth/itmodel"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func HandlePomodoroReport(Profile itmodel.Profile, Pesan itmodel.IteungMessage, db *mongo.Database) string {
	// 1. Validasi input dasar

	// Tambahkan pengambilan data user untuk mendapatkan nama
	userData, err := GetUserData(Profile, Pesan, db)
	var userName string
	if err != nil {
		// Jika gagal, gunakan alias_name sebagai fallback
		userName = Pesan.Alias_name
		fmt.Printf("Warning: Failed to get user data: %v. Using alias_name instead.\n", err)
	} else {
		// Jika berhasil, gunakan nama dari data user
		userName = userData.Name
	}
	if Pesan.Message == "" {
		return "Wah kak " + userName + ", pesan tidak boleh kosong"
	}

	cycle := extractCycleNumber(Pesan.Message)
	if cycle == 0 {
		return "Wah kak " + userName + ", format cycle tidak valid. Contoh: 'Iteung Pomodoro Report 1 cycle'"
	}

	hostname := extractValue(Pesan.Message, "Hostname : ")
	ip := extractIP(Pesan.Message)
	screenshots := extractNumber(Pesan.Message, "Jumlah ScreenShoot : ")
	pekerjaan := extractActivities(Pesan.Message)
	token := extractToken(Pesan.Message)

	// 3. Verifikasi public key
	publicKey, err := getPublicKey(db)
	if err != nil {
		return "Wah kak " + userName + ", sistem gagal memuat public key: " + err.Error()
	}

	// Cek apakah token sudah pernah digunakan di koleksi pomokit
	isUsed, err := isTokenUsed(db, token)
	if err != nil {
		return "Wah kak " + userName + ", sistem gagal memeriksa token: " + err.Error()
	}

	if isUsed {
		return "Wah kak " + userName + ", token ini sudah pernah digunakan sebelumnya"
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
			userName,
			errorMsg,
			strings.Split(err.Error(), ":")[0],
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
		Name:        userName,
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
		return "Wah kak " + userName + ", gagal menyimpan laporan: " + err.Error()
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
		userName,
		hostname,
		ip,
		pekerjaan,
		url,
		report.CreatedAt.Format("2006-01-02 🕒15:04 WIB"),
	)
}

// Tambahkan fungsi untuk mengambil data user
func GetUserData(Profile itmodel.Profile, Pesan itmodel.IteungMessage, db *mongo.Database) (daftar.Userdomyikado, error) {
    var result daftar.Userdomyikado
    conf, err := atdb.GetOneDoc[Config](db, "config", bson.M{"phonenumber": Profile.Phonenumber})
    if err != nil {
        return result, fmt.Errorf("gagal mengambil config: %v", err)
    }
    
    // Mendapatkan semua data user
    statusCode, allUsers, err := atapi.GetWithToken[[]daftar.Userdomyikado]("login", Profile.Token, conf.DomyikadoAllUserURL)
    
    if err != nil {
        return result, err
    }
    
    if statusCode != 200 {
        return result, fmt.Errorf("failed to get user data: status code %d", statusCode)
    }
    
    // Gunakan phoneNumber dari parameter (Pesan.Phone_number di pemanggilan fungsi)
    targetPhoneNumber := Pesan.Phone_number
    
    // Mencari user dengan nomor telepon yang sesuai
    for _, user := range allUsers {
        if user.PhoneNumber == targetPhoneNumber {
            return user, nil
        }
    }
    
    return result, fmt.Errorf("user dengan nomor telepon %s tidak ditemukan", targetPhoneNumber)
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

	// Tambahkan pengambilan data user untuk mendapatkan nama
	userData, err := GetUserData(Profile, Pesan, db)
	var userName string
	if err != nil {
		// Jika gagal, gunakan alias_name sebagai fallback
		userName = Pesan.Alias_name
		fmt.Printf("Warning: Failed to get user data: %v. Using alias_name instead.\n", err)
	} else {
		// Jika berhasil, gunakan nama dari data user
		userName = userData.Name
	}

	if Pesan.Message == "" {
		return "Wah kak " + userName + ", pesan tidak boleh kosong"
	}

	// Pisahkan pesan menjadi baris-baris
	lines := strings.Split(Pesan.Message, "\n")
	
	// Bersihkan setiap baris dari spasi berlebih
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
	
	// Ekstrak cycle dari baris pertama atau dari seluruh pesan jika tidak ditemukan
	cycle := 0
	if strings.Contains(lines[0], "Start") && strings.Contains(lines[0], "cycle") {
		cycle = extractStartCycleNumber(lines[0])
	} else {
		cycle = extractStartCycleNumber(Pesan.Message)
	}
	
	// Validasi cycle
	if cycle == 0 {
		return "Wah kak " + userName + ", format cycle tidak valid. Contoh: 'Pomodoro Start 1 cycle'"
	}

	// Ekstrak nilai-nilai menggunakan regex yang lebih fleksibel
	milestone := extractWithRegex(lines, `Milestone\s*:\s*(.+)`)
	version := extractWithRegex(lines, `Version\s*:\s*(.+)`)
	hostname := extractWithRegex(lines, `Hostname\s*:\s*(.+)`)
	ipRaw := extractWithRegex(lines, `IP\s*:\s*(.+)`)
	
	// Format IP jika perlu
	ip := ipRaw
	if !strings.HasPrefix(ipRaw, "https://whatismyipaddress.com") && ipRaw != "" {
		// Cek apakah ini adalah alamat IP
		ipRegex := regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)
		ipMatch := ipRegex.FindStringSubmatch(ipRaw)
		if len(ipMatch) > 1 {
			ip = "https://whatismyipaddress.com/ip/" + ipMatch[1]
		}
	}
	
	// Set nilai default jika kosong
	if version == "" {
		version = "1.0.0"
	}
	
	if milestone == "" {
		milestone = "Tidak ada milestone"
	}

	// Lokasi waktu Indonesia
	loc, _ := time.LoadLocation("Asia/Jakarta")
	currentTime := time.Now().In(loc)

	// Format respons dengan baris baru yang jelas antara tiap bagian
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
		userName,
		milestone,
		version,
		hostname,
		ip,
		currentTime.Format("2006-01-02 🕒15:04 WIB"),
	)
}

// Fungsi untuk ekstraksi cycle dari pesan Start
func extractStartCycleNumber(msg string) int {
	re := regexp.MustCompile(`Start\s+(\d+)\s+cycle`)
	matches := re.FindStringSubmatch(msg)
	if len(matches) > 1 {
		cycle, _ := strconv.Atoi(matches[1])
		return cycle
	}
	return 0
}

// Fungsi untuk mengekstrak nilai dengan regex yang lebih fleksibel dan menghindari kontaminasi nilai
func extractWithRegex(lines []string, pattern string) string {
	re := regexp.MustCompile(pattern)
	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) > 1 {
			// Ekstrak bagian yang sesuai dan hapus teks pattern lain yang mungkin terbawa
			value := matches[1]
			// Bersihkan dari pattern field lain yang mungkin tercampur
			cleanValue := strings.Split(value, "Version")[0]
			cleanValue = strings.Split(cleanValue, "Hostname")[0]
			cleanValue = strings.Split(cleanValue, "IP")[0]
			return strings.TrimSpace(cleanValue)
		}
	}
	return ""
}