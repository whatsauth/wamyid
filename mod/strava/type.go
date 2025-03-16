package strava

import (
	"time"
)

type Config struct {
	PhoneNumber      string `json:"phonenumber,omitempty" bson:"phonenumber,omitempty"`
	LeaflyURL        string `json:"leaflyurl,omitempty" bson:"leaflyurl,omitempty"`
	LeaflySecret     string `json:"leaflysecret,omitempty" bson:"leaflysecret,omitempty"`
	DomyikadoUserURL string `json:"domyikadouserurl,omitempty" bson:"domyikadouserurl,omitempty"`
	DomyikadoSecret  string `json:"domyikadosecret,omitempty" bson:"domyikadosecret,omitempty"`
	KimseokgisURL    string `json:"urlkimseokgis" bson:"urlkimseokgis"`
}

type StravaIdentity struct {
	AthleteId     string    `bson:"athlete_id" json:"athlete_id"`
	Picture       string    `bson:"picture" json:"picture"`
	Name          string    `bson:"name" json:"name"`
	PhoneNumber   string    `bson:"phone_number" json:"phone_number"`
	LinkIndentity string    `bson:"link_identity" json:"link_identity"`
	CreatedAt     time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at" json:"updated_at"`
}

type StravaActivity struct {
	ActivityId   string    `bson:"activity_id" json:"activity_id"`
	Picture      string    `bson:"picture" json:"picture"`
	Name         string    `bson:"name" json:"name"`
	Title        string    `bson:"title" json:"title"`
	DateTime     string    `bson:"date_time" json:"date_time"`
	TypeSport    string    `bson:"type_sport" json:"type_sport"`
	Distance     string    `bson:"distance" json:"distance"`
	MovingTime   string    `bson:"moving_time" json:"moving_time"`
	Elevation    string    `bson:"elevation" json:"elevation"`
	LinkActivity string    `bson:"link_activity" json:"link_activity"`
	Status       string    `bson:"status" json:"status"`
	CreatedAt    time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at" json:"updated_at"`
	// Location     []GeoPoint `bson:"location" json:"created_at"`
}

// Struktur untuk data lokasi dari Strava
// type StravaData struct {
// 	Props struct {
// 		PageProps struct {
// 			Activity struct {
// 				Streams struct {
// 					Location []struct {
// 						Lat float64 `json:"lat"`
// 						Lng float64 `json:"lng"`
// 					} `json:"location"`
// 				} `json:"streams"`
// 			} `json:"activity"`
// 		} `json:"pageProps"`
// 	} `json:"props"`
// }

// type GeoPoint struct {
// 	Type        string    `bson:"type"`
// 	Coordinates []float64 `bson:"coordinates"`
// }
