package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Author struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" query:"id" url:"_id,omitempty" reqHeader:"id"`
	Profile   string             `json:"profile,omitempty" bson:"profile,omitempty"`
	Nama      string             `json:"nama,omitempty" bson:"nama,omitempty"`
	NIK       string             `json:"nik,omitempty" bson:"nik,omitempty"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty"`
	Phone     string             `json:"phone,omitempty" bson:"phone,omitempty"`
	Pekerjaan string             `json:"pekerjaan,omitempty" bson:"pekerjaan,omitempty"`
	Alamat    string             `json:"alamat,omitempty" bson:"alamat,omitempty"`
	Bio       string             `json:"bio,omitempty" bson:"bio,omitempty"`
	Photo     string             `json:"photo,omitempty" bson:"photo,omitempty"`
}

type Login struct {
	Login string `json:"login,omitempty" bson:"login,omitempty" query:"login" url:"login,omitempty" reqHeader:"login"`
}
