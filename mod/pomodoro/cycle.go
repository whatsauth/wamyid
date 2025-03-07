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

func HandlePomodoroStart(Pesan itmodel.IteungMessage, db *mongo.Database) string {
    // Ambil cycle terakhir
    lastCycle := getLastCycle(db, Pesan.Phone_number)
    
    // Buat cycle baru
    newCycle := lastCycle + 1
    
    // Simpan ke database
    pomodoro := Pomodoro{
        PhoneNumber: Pesan.Phone_number,
        Cycle:       newCycle,
        StartTime:   time.Now(),
        Milestone:   extractMilestone(Pesan.Message),
    }
    
    _, err := atdb.InsertOneDoc(db, "pomodoros", pomodoro)
    if err != nil {
        return "Gagal memulai cycle 😥"
    }
    
    return fmt.Sprintf(
        "🎯 *Mulai Cycle %d*\n"+
        "Hai %s\n"+
        "Waktu: %s\n"+
        "Milestone: %s",
        newCycle,
        Pesan.Alias_name,
        time.Now().Format("15:04"),
        pomodoro.Milestone,
    )
}

func HandlePomodoroReport(Pesan itmodel.IteungMessage, db *mongo.Database) string {
    // Ambil cycle terakhir
    lastCycle := getLastCycle(db, Pesan.Phone_number)
    
    return fmt.Sprintf(
        "✅ *Cycle %d Selesai*\n"+
        "Nama: %s\n"+
        "Durasi: 25 menit\n"+
        "Milestone: %s",
        lastCycle,
        Pesan.Alias_name,
        getLastMilestone(db, Pesan.Phone_number),
    )
}

// Helper functions
func getLastCycle(db *mongo.Database, phone string) int {
    filter := bson.M{"phonenumber": phone}
    result, _ := atdb.GetOneLatestDoc[Pomodoro](db, "pomodoro", filter)
    return result.Cycle
}

func extractMilestone(msg string) string {
    parts := strings.Split(msg, "Milestone : ")
    if len(parts) > 1 {
        return strings.TrimSpace(parts[1])
    }
    return "Tidak ada milestone"
}

func getLastMilestone(db *mongo.Database, phone string) string {
    filter := bson.M{"phonenumber": phone}
    result, _ := atdb.GetOneLatestDoc[Pomodoro](db, "pomodoro", filter)
    return result.Milestone
}