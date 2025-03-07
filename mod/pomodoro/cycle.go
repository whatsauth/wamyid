package pomodoro

import (
    "fmt"
    "strings"
    "time"
    
	"github.com/gocroot/helper/atdb"
    "github.com/whatsauth/itmodel"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

func PomodoroHandler(Pesan itmodel.IteungMessage, db *mongo.Database) string {
    switch {
    case strings.HasPrefix(Pesan.Message, "Pomodoro Start"):
        handleStart(Pesan, db)
        return "" // Tidak ada respons saat mulai
    case strings.HasPrefix(Pesan.Message, "Pomodoro Report"):
        return handleReport(Pesan, db)
    default:
        return "Command tidak valid. Gunakan:\n- Pomodoro Start\n- Pomodoro Report"
    }
}

func handleStart(Pesan itmodel.IteungMessage, db *mongo.Database) {
    milestone := extractMilestone(Pesan.Message)
    
    // Auto-increment cycle
    lastCycle := getLastCycle(db, Pesan.Phone_number)
    newCycle := lastCycle + 1

    // Simpan ke database tanpa memberikan respons
    pomodoro := Pomodoro{
        PhoneNumber: Pesan.Phone_number,
        Cycle:       newCycle,
        StartTime:   time.Now(),
        Milestone:   milestone,
    }
    
    atdb.InsertOneDoc(db, "pomodoros", pomodoro)
}

func handleReport(Pesan itmodel.IteungMessage, db *mongo.Database) string {
    lastEntry, err := getLastEntry(db, Pesan.Phone_number)
    if err != nil {
        return "Wah kak " + Pesan.Alias_name + ", belum ada cycle yang dimulai"
    }

    duration := time.Since(lastEntry.StartTime).Round(time.Minute)
    
    return fmt.Sprintf(
        "✅ *Laporan Cycle %d*\n"+
        "Nama: %s\n"+
        "Milestone: %s\n"+
        "Durasi: %s\n"+
        "Mulai: %s\n"+
        "Selesai: %s",
        lastEntry.Cycle,
        Pesan.Alias_name,
        lastEntry.Milestone,
        duration,
        lastEntry.StartTime.Format("15:04"),
        time.Now().Format("15:04"),
    )
}

// Helper functions
func extractMilestone(msg string) string {
    parts := strings.SplitN(msg, "Milestone : ", 2)
    if len(parts) > 1 {
        return strings.TrimSpace(parts[1])
    }
    return "Tidak ada milestone"
}

func getLastEntry(db *mongo.Database, phone string) (Pomodoro, error) {
    filter := bson.M{"phonenumber": phone}
    // Tambahkan type parameter Pomodoro dalam kurung siku
    result, err := atdb.GetOneLatestDoc[Pomodoro](db, "pomodoro", filter)
    if err != nil {
        return Pomodoro{}, err
    }
    return result, nil
}

func getLastCycle(db *mongo.Database, phone string) int {
    // Tambahkan type parameter Pomodoro
    result, _ := atdb.GetOneLatestDoc[Pomodoro](db, "pomodoro", bson.M{"phonenumber": phone})
    return result.Cycle
}