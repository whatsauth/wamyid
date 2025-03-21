package strava

import (
	"context"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// StravaInfo adalah struktur data untuk menyimpan informasi poin pengguna
type StravaInfo struct {
	Name        string  `json:"name"`
	PhoneNumber string  `json:"phone_number"`
	TotalKm     float64 `json:"total_km"`
	Poin        float64 `json:"poin"`
	Count       int     `json:"count"`
}

// TambahPoinDariAktivitas menghitung total KM dan menambahkan poin ke strava_poin
func TambahPoinDariAktivitas(db *mongo.Database, phone string) error {
	colActivity := "strava_activity"
	colPoin := "strava_poin"

	// Ambil semua aktivitas berdasarkan phone_number
	cursor, err := db.Collection(colActivity).Find(context.TODO(), bson.M{"phone_number": phone})
	if err != nil {
		return err
	}
	defer cursor.Close(context.TODO())

	totalKm := 0.0
	for cursor.Next(context.TODO()) {
		var activity struct {
			Distance string `bson:"distance"`
		}
		if err := cursor.Decode(&activity); err != nil {
			return err
		}

		// Hapus " km" dan konversi ke float
		distanceStr := strings.Replace(activity.Distance, " km", "", -1)
		distance, err := strconv.ParseFloat(distanceStr, 64)
		if err != nil {
			continue // Skip jika gagal parsing
		}

		totalKm += distance
	}

	if totalKm == 0 {
		return nil // Tidak ada aktivitas valid
	}

	// Update atau Insert ke strava_poin
	filter := bson.M{"phone_number": phone}
	update := bson.M{
		"$set": bson.M{"total_km": totalKm},
		"$inc": bson.M{
			"poin":  (totalKm / 6) * 100, // Konversi KM ke poin
			"count": 1,
		},
	}
	opts := options.Update().SetUpsert(true)

	_, err = db.Collection(colPoin).UpdateOne(context.TODO(), filter, update, opts)
	return err
}
