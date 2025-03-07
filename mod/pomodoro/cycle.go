package pomodoro

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocroot/helper/atdb"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/mongo"
)

func HandlePomodoroReport(Pesan itmodel.IteungMessage, db *mongo.Database) string {
    // Ekstrak angka dari pesan
    cycle := extractCycleNumber(Pesan.Message)
    if cycle == 0 {
        return "Wah kak " + Pesan.Alias_name + ", format cycle tidak valid"
    }

    // Parse data lainnya
    report := PomodoroReport{
        PhoneNumber: Pesan.Phone_number,
        Cycle:       cycle,
        Hostname:    extractValue(Pesan.Message, "Hostname : "),
        IP:          extractValue(Pesan.Message, "IP : "),
        Screenshots: extractNumber(Pesan.Message, "Jumlah ScreenShoot : "),
        Aktivitas:   extractActivities(strings.Split(Pesan.Message, "\n")),
        Signature:   extractSignature(Pesan.Message),
        CreatedAt:   time.Now(),
    }

    // Simpan ke collection pomodoro-cyclez sesuai struktur database
    _, err := atdb.InsertOneDoc(db, "pomodoro-cyclez", report)
    if err != nil {
        return "Gagal menyimpan laporan cycle: " + err.Error()
    }

    return generatePomodoroResponse(report, Pesan.Alias_name)
}

func extractCycleNumber(msg string) int {
    re := regexp.MustCompile(`\d+`)
    matches := re.FindStringSubmatch(msg)
    if len(matches) > 0 {
        cycle, _ := strconv.Atoi(matches[0])
        return cycle
    }
    return 0
}

func extractValue(line, prefix string) string {
    return strings.TrimSpace(strings.TrimPrefix(line, prefix))
}

func extractNumber(line, prefix string) int {
    valStr := extractValue(line, prefix)
    num, _ := strconv.Atoi(valStr)
    return num
}

func extractActivities(lines []string) []string {
    var activities []string
    for _, line := range lines {
        if strings.HasPrefix(line, "|") {
            activities = append(activities, strings.TrimPrefix(line, "|"))
        }
    }
    return activities
}

func extractSignature(msg string) string {
    parts := strings.Split(msg, "#")
    if len(parts) > 1 {
        return parts[1]
    }
    return ""
}

func generatePomodoroResponse(report PomodoroReport, name string) string {
    return fmt.Sprintf(
        "🍅 *Laporan Cycle %d Berhasil!*\n"+
        "Nama: %s\n"+
        "Hostname: %s\n"+
        "IP: %s\n"+
        "Aktivitas:\n- %s\n"+
        "Timestamp: %s",
        report.Cycle,
        name,
        report.Hostname,
        report.IP,
        strings.Join(report.Aktivitas, "\n- "),
        report.CreatedAt.Format("2006-01-02 15:04:05"),
    )
}