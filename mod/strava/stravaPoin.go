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
	cursor, err := db.Collection(colActivity).Find(context.TODO(), bson.M{})
	if err != nil {
		return "Error fetching activities: " + err.Error()
	}
	defer cursor.Close(context.TODO())

	// Map untuk menyimpan total km per pengguna
	userData := make(map[string]float64)

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
		userData[activity.PhoneNumber] += distance
	}

	// Update atau insert poin ke `strava_poin`
	for phone, totalKm := range userData {
		filter := bson.M{"phone_number": phone}
		update := bson.M{
			"$set": bson.M{"total_km": totalKm},
			"$inc": bson.M{
				"poin":           (totalKm / 6) * 100, // Konversi km ke poin
				"activity_count": 1,
			},
		}
		opts := options.Update().SetUpsert(true) // Insert jika belum ada

		_, err := db.Collection(colPoin).UpdateOne(context.TODO(), filter, update, opts)
		if err != nil {
			log.Println("Error updating strava_poin:", err)
		}
	}

	return "Proses inisialisasi poin dari aktivitas lama selesai."
}
