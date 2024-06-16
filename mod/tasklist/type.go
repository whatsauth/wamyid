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

type Config struct {
	PhoneNumber          string `json:"phonenumber,omitempty" bson:"phonenumber,omitempty"`
	LeaflyURL            string `json:"leaflyurl,omitempty" bson:"leaflyurl,omitempty"`
	LeaflySecret         string `json:"leaflysecret,omitempty" bson:"leaflysecret,omitempty"`
	DomyikadoPresensiURL string `json:"domyikadopresensiurl,omitempty" bson:"domyikadopresensiurl,omitempty"`
	DomyikadoTaskListURL string `json:"domyikadotasklisturl,omitempty" bson:"domyikadotasklisturl,omitempty"`
	DomyikadoSecret      string `json:"domyikadosecret,omitempty" bson:"domyikadosecret,omitempty"`
}
