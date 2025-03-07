package pomodoro

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PomodoroReport struct {
    ID            primitive.ObjectID `bson:"_id,omitempty"`
    PhoneNumber   string             `bson:"phonenumber,omitempty"`
    Cycle         int                `bson:"cycle"`
    Hostname      string             `bson:"hostname"`
    IP            string             `bson:"ip"`
    Screenshots   int                `bson:"screenshots"`
    Pekerjaan     []string           `bson:"pekerjaan"`
    Token         string             `bson:"token"`
    CreatedAt     time.Time          `bson:"createdAt"`
}

type Config struct {
    PublicKeyPomokit    string `json:"publickeypomokit,omitempty" bson:"publickeypomokit,omitempty"`
}