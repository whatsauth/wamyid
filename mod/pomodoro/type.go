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