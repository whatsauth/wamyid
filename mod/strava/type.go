package strava

import (
	"time"
)

type StravaIdentity struct {
	AthleteId     string    `bson:"athlete_id" json:"athlete_id"`
	Picture       string    `bson:"picture" json:"picture"`
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
	CreatedAt    time.Time `bson:"created_at" json:"created_at"`
}
