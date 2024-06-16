package tasklist

import "go.mongodb.org/mongo-driver/bson/primitive"

type TaskList struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	LaporanID   primitive.ObjectID `json:"laporanid,omitempty" bson:"laporanid,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	PhoneNumber string             `json:"phonenumber,omitempty" bson:"phonenumber,omitempty"`
	Task        string             `json:"task,omitempty" bson:"task,omitempty"`
	IsDone      bool               `json:"isdone,omitempty" bson:"isdone,omitempty"`
}
