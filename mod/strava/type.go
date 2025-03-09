package strava

import (
	"time"
)

type StravaIdentity struct {
	AthleteId     string    `bson:"athlete_id,omitempty" json:"athlete_id,omitempty"`
	Picture       string    `bson:"picture,omitempty" json:"picture,omitempty"`
	PhoneNumber   string    `bson:"phone_number,omitempty" json:"phone_number,omitempty"`
	LinkIndentity string    `bson:"link_identity,omitempty" json:"link_identity,omitempty"`
	CreatedAt     time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at" json:"updated_at"`
}

type StravaActivity struct {
	ActivityId   string    `bson:"activity_id,omitempty" json:"activity_id,omitempty"`
	Picture      string    `bson:"picture,omitempty" json:"picture,omitempty"`
	Name         string    `bson:"name,omitempty" json:"name,omitempty"`
	Title        string    `bson:"title,omitempty" json:"title,omitempty"`
	DateTime     string    `bson:"date_time,omitempty" json:"date_time,omitempty"`
	TypeSport    string    `bson:"type_sport,omitempty" json:"type_sport,omitempty"`
	Distance     string    `bson:"distance,omitempty" json:"distance,omitempty"`
	MovingTime   string    `bson:"moving_time,omitempty" json:"moving_time,omitempty"`
	Elevation    string    `bson:"elevation,omitempty" json:"elevation,omitempty"`
	LinkActivity string    `bson:"link_activity,omitempty" json:"link_activity,omitempty"`
	CreatedAt    time.Time `bson:"created_at" json:"created_at"`
}
