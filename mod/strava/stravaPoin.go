package strava

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StravaInfo struct {
	PhoneNumber   string  `bson:"phone_number"`
	TotalKm       float64 `bson:"total_km"`
	Poin          float64 `bson:"poin"`
	ActivityCount int     `bson:"activity_count"`
}

// Proses semua aktivitas lama ke dalam strava_poin
func InisialisasiPoinDariAktivitasLama(Pesan itmodel.IteungMessage, db *mongo.Database) string {
	colActivity := "strava_activity"
	colPoin := "strava_poin"

	// Ambil semua aktivitas dari strava_activity
	filter := bson.M{"status": "Valid"}
	cursor, err := db.Collection(colActivity).Find(context.TODO(), filter)
	if err != nil {
		return "Error fetching activities: " + err.Error()
	}
	defer cursor.Close(context.TODO())

	// Map untuk menyimpan total km per pengguna
	userData := make(map[string]struct {
		TotalKm       float64
		ActivityCount int
	})

	// Loop melalui semua aktivitas
	for cursor.Next(context.TODO()) {
		var activity struct {
			PhoneNumber string `bson:"phone_number"`
			Distance    string `bson:"distance"`
		}
		if err := cursor.Decode(&activity); err != nil {
			log.Println("Error decoding activity:", err)
			continue
		}

		// Konversi jarak dari string ke float64 (misal: "28.6 km" → 28.6)
		distanceStr := strings.Replace(activity.Distance, " km", "", -1)
		distance, err := strconv.ParseFloat(distanceStr, 64)
		if err != nil {
			log.Println("Error converting distance:", err)
			continue
		}

		// Tambahkan ke map berdasarkan nomor telepon
		userData[activity.PhoneNumber] = struct {
			TotalKm       float64
			ActivityCount int
		}{
			TotalKm:       userData[activity.PhoneNumber].TotalKm + distance,
			ActivityCount: userData[activity.PhoneNumber].ActivityCount + 1,
		}
	}

	// Update atau insert poin ke `strava_poin`
	for phone, data := range userData {
		filter := bson.M{"phone_number": phone}

		// Ambil count sebelumnya dari `strava_poin`
		var existing struct {
			ActivityCount int `bson:"activity_count"`
		}
		err := db.Collection(colPoin).FindOne(context.TODO(), filter).Decode(&existing)
		if err != nil && err != mongo.ErrNoDocuments {
			log.Println("Error fetching existing count:", err)
			continue
		}

		// Jumlah total aktivitas valid dari sebelumnya + yang baru dihitung
		newCount := existing.ActivityCount + data.ActivityCount

		update := bson.M{
			"$set": bson.M{
				"total_km": data.TotalKm,
				"count":    newCount, // Update total count
			},
			"$inc": bson.M{
				"poin": (data.TotalKm / 6) * 100, // Konversi km ke poin
			},
		}
		opts := options.Update().SetUpsert(true) // Insert jika belum ada

		_, err = db.Collection(colPoin).UpdateOne(context.TODO(), filter, update, opts)
		if err != nil {
			log.Println("Error updating strava_poin:", err)
		}
	}

	return "Proses inisialisasi poin dari aktivitas lama selesai."
}
