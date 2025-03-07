package pomodoro

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gocroot/helper/atdb"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/mongo"
)

func HandlePomodoroReport(Pesan itmodel.IteungMessage, db *mongo.Database) string {
    // Validasi format pesan
    if !strings.HasPrefix(Pesan.Message, "Iteung Pomodoro Report") {
        return "Format laporan tidak valid"
    }

    // Parse data dari pesan
    lines := strings.Split(Pesan.Message, "\n")
    
    report := PomodoroReport{
        PhoneNumber: Pesan.Phone_number,
        Cycle:       extractCycle(lines[0]),
        Hostname:    extractValue(lines[1], "Hostname : "),
        IP:          extractValue(lines[2], "IP : "),
        Screenshots: extractNumber(lines[3], "Jumlah ScreenShoot : "),
        Aktivitas:   extractActivities(lines[4:]),
        Signature:   extractSignature(Pesan.Message),
        CreatedAt:   time.Now(),
    }

    // Simpan ke database
    _, err := atdb.InsertOneDoc(db, "pomodoro", report)
    if err != nil {
        return "Gagal menyimpan laporan: " + err.Error()
    }

    return generateResponse(report, Pesan.Alias_name)
}

func extractCycle(line string) int {
    parts := strings.Split(line, " ")
    if len(parts) < 4 {
        return 0
    }
    cycle, _ := strconv.Atoi(parts[3])
    return cycle
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

func generateResponse(report PomodoroReport, name string) string {
    return fmt.Sprintf(
        "✅ *Laporan Pomodoro Cycle %d*\n"+
        "Nama: %s\n"+
        "Hostname: %s\n"+
        "Durasi: %s\n"+
        "Screenshots: %d\n"+
        "Aktivitas:\n- %s",
        report.Cycle,
        name,
        report.Hostname,
        time.Since(report.CreatedAt).Round(time.Minute).String(),
        report.Screenshots,
        strings.Join(report.Aktivitas, "\n- "),
    )
}