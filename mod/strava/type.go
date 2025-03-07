package strava

import "go.mongodb.org/mongo-driver/bson/primitive"

type StravaActivity struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ActivityId string             `bson:"activity_id,omitempty" json:"activity_id,omitempty"`
	Name       string             `bson:"name,omitempty" json:"name,omitempty"`
	Title      string             `bson:"title,omitempty" json:"title,omitempty"`
	DateTime   string             `bson:"date_time,omitempty" json:"date_time,omitempty"`
	TypeSport  string             `bson:"type_sport,omitempty" json:"type_sport,omitempty"`
	Distance   string             `bson:"distance,omitempty" json:"distance,omitempty"`
	TimePeriod string             `bson:"timePeriod,omitempty" json:"time_period,omitempty"`
	Elevation  string             `bson:"elevation,omitempty" json:"elevation,omitempty"`
}
