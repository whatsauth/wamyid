package strava

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StravaInfo struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `json:"name"`
	PhoneNumber string             `json:"phone_number"`
	TotalKm     float64            `json:"total_km"`
	Poin        float64            `json:"poin"`
	Count       int                `json:"count"`
}

func TambahPoinDariAktivitas(db *mongo.Database, phone string) error {
	colActivity := "strava_activity"
	colPoin := "strava_poin"

	// Ambil total km dari koleksi strava_activity
	match := bson.M{"phone_number": phone}
	group := bson.M{"_id": nil, "totalKm": bson.M{"$sum": bson.M{"$toDouble": "$distance"}}}

	cursor, err := db.Collection(colActivity).Aggregate(context.TODO(), mongo.Pipeline{
		{{Key: "$match", Value: match}},
		{{Key: "$group", Value: group}},
	})
	if err != nil {
		return err
	}
	defer cursor.Close(context.TODO())

	var result struct {
		TotalKm float64 `bson:"totalKm"`
	}
	if cursor.Next(context.TODO()) {
		if err := cursor.Decode(&result); err != nil {
			return err
		}
	} else {
		// Jika tidak ada aktivitas, tidak perlu update poin
		return nil
	}

	// Update atau insert data ke strava_poin
	filter := bson.M{"phone_number": phone}
	update := bson.M{
		"$set": bson.M{"total_km": result.TotalKm},
		"$inc": bson.M{
			"poin":  (result.TotalKm / 6) * 100, // Konversi km ke poin
			"count": 1,
		},
	}
	opts := options.Update().SetUpsert(true) // Upsert otomatis insert jika belum ada

	_, err = db.Collection(colPoin).UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}
