package siakad

import (
	"time"
)

type Config struct {
	PhoneNumber               string `json:"phonenumber,omitempty" bson:"phonenumber,omitempty"`
	SiakadLoginURL            string `json:"siakadloginurl,omitempty" bson:"siakadloginurl,omitempty"`
	BapURL                    string `json:"bapurl,omitempty" bson:"bapurl,omitempty"`
	ApproveBapURL             string `json:"approvebapurl,omitempty" bson:"approvebapurl,omitempty"`
	ApproveBimbinganURL       string `json:"approvebimbinganurl,omitempty" bson:"approvebimbinganurl,omitempty"`
	ApproveBimbinganByPoinURL string `json:"approvebimbinganbypoinurl,omitempty" bson:"approvebimbinganbypoinurl,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type LoginInfo struct {
	NoHp      string    `json:"nohp" bson:"nohp"`
	Email     string    `json:"email" bson:"email"`
	LoginTime time.Time `json:"login_time" bson:"login_time"`
}

type Prompt struct {
	Prompt string `bson:"prompt"`
	Answer string `bson:"answer"`
}
