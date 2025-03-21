package strava

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// StravaInfo untuk menyimpan poin dan total km
type StravaInfo struct {
	Name        string  `json:"name"`
	PhoneNumber string  `json:"phone_number"`
	TotalKm     float64 `json:"total_km"`
	Poin        float64 `json:"poin"`
	Count       int     `json:"count"`
}

// Fungsi untuk membaca aktivitas yang sudah tersimpan dan menambah poin
func TambahPoinDariAktivitas(db *mongo.Database, phone string) error {
	colActivity := "strava_activity"
	colPoin := "strava_poin"

	// 1. Ambil semua aktivitas berdasarkan nomor telepon
	cursor, err := db.Collection(colActivity).Find(context.TODO(), bson.M{"phone_number": phone})
	if err != nil {
		return err
	}
	defer cursor.Close(context.TODO())

	var totalKm float64

	for cursor.Next(context.TODO()) {
		var activity struct {
			Distance string `bson:"distance"`
		}
		if err := cursor.Decode(&activity); err != nil {
			return err
		}

		// Hapus " km" jika ada dan konversi ke float64
		distanceStr := strings.Replace(activity.Distance, " km", "", -1)
		distance, err := strconv.ParseFloat(distanceStr, 64)
		if err != nil {
			fmt.Println("Gagal mengonversi jarak:", activity.Distance)
			continue
		}
		totalKm += distance
	}

	// Jika tidak ada aktivitas, tidak perlu update
	if totalKm == 0 {
		return nil
	}

	// 2. Perbarui atau buat data di strava_poin
	filter := bson.M{"phone_number": phone}
	update := bson.M{
		"$set": bson.M{
			"total_km": totalKm,
			"poin":     (totalKm / 6) * 100, // Konversi ke poin
		},
		"$inc": bson.M{
			"count": 1, // Jumlah update
		},
	}
	opts := options.Update().SetUpsert(true) // Buat dokumen baru jika belum ada

	_, err = db.Collection(colPoin).UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}
