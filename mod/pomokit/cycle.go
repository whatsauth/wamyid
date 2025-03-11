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
    if Pesan.Message == "" {
        return "Wah kak " + Pesan.Alias_name + ", pesan tidak boleh kosong"
    }

    // Normalisasi line endings dan split pesan menjadi baris-baris
    normalizedMsg := strings.ReplaceAll(strings.ReplaceAll(Pesan.Message, "\r\n", "\n"), "\r", "\n")
    lines := strings.Split(normalizedMsg, "\n")
    
    // Variabel untuk menyimpan nilai
    cycle := 0
    milestone := ""
    version := ""
    hostname := ""
    ip := ""
    
    // Ekstrak cycle dari baris pertama
    if len(lines) > 0 {
        firstLine := strings.ToLower(lines[0])
        if strings.Contains(firstLine, "start") && strings.Contains(firstLine, "cycle") {
            // Regex simple untuk menemukan angka
            re := regexp.MustCompile(`(\d+)`)
            matches := re.FindStringSubmatch(firstLine)
            if len(matches) > 0 {
                cycleNum, err := strconv.Atoi(matches[0])
                if err == nil {
                    cycle = cycleNum
                }
            }
        }
    }
    
    if cycle == 0 {
        return "Wah kak " + Pesan.Alias_name + ", format cycle tidak valid. Contoh: 'Pomodoro Start 1 cycle'"
    }
    
    // Parsing nilai-nilai lain dari setiap baris (satu baris satu field)
    for _, line := range lines[1:] {
        line = strings.TrimSpace(line)
        if line == "" {
            continue // Lewati baris kosong
        }
        
        // Cari pemisah ":" dalam baris
        parts := strings.SplitN(line, ":", 2)
        if len(parts) != 2 {
            continue // Lewati baris yang tidak memiliki ":"
        }
        
        // Ekstrak key dan value
        key := strings.ToLower(strings.TrimSpace(parts[0]))
        value := strings.TrimSpace(parts[1])
        
        // Lewati jika value kosong
        if value == "" {
            continue
        }
        
        // Tetapkan nilai berdasarkan key
        switch key {
        case "milestone":
            milestone = value
        case "version":
            version = value
        case "hostname":
            hostname = value
        case "ip":
            ip = value
            // Format IP jika diperlukan
            if !strings.HasPrefix(ip, "https://") && strings.Contains(ip, ".") {
                ipRegex := regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)
                if match := ipRegex.FindStringSubmatch(ip); len(match) > 1 {
                    ip = "https://whatismyipaddress.com/ip/" + match[1]
                }
            }
        }
    }
    
    // Periksa field yang kosong
    var missingFields []string
    if milestone == "" {
        missingFields = append(missingFields, "milestone")
    }
    if version == "" {
        missingFields = append(missingFields, "version")
    }
    if hostname == "" {
        missingFields = append(missingFields, "hostname")
    }
    if ip == "" {
        missingFields = append(missingFields, "ip")
    }
    
    // Jika ada field yang tidak ditemukan, kirim pesan error
    if len(missingFields) > 0 {
        return fmt.Sprintf("Wah kak %s, beberapa informasi penting belum diisi: %s. Mohon lengkapi ya!", 
            Pesan.Alias_name, 
            strings.Join(missingFields, ", "))
    }
    
    // Format waktu
    loc, _ := time.LoadLocation("Asia/Jakarta")
    currentTime := time.Now().In(loc)
    
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