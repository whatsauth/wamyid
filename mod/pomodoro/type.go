package pomodoro

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Pomodoro struct {
    ID          primitive.ObjectID `bson:"_id,omitempty"`
    PhoneNumber string             `bson:"phonenumber"`
    Cycle       int                `bson:"cycle"`
    StartTime   time.Time          `bson:"start_time"`
    Milestone   string             `bson:"milestone,omitempty"`
}

type PomodoroReport struct {
    ID            primitive.ObjectID `bson:"_id,omitempty"`
    PhoneNumber   string             `bson:"phonenumber,omitempty"`
    Cycle         int                `bson:"cycle"`
    Hostname      string             `bson:"hostname"`
    IP            string             `bson:"ip"`
    Screenshots   int                `bson:"screenshots"`
    Aktivitas     []string           `bson:"aktivitas"`
    Signature     string             `bson:"signature"`
    CreatedAt     time.Time          `bson:"createdAt"`
}