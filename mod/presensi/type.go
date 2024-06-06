package presensi

import "go.mongodb.org/mongo-driver/bson/primitive"

type Lokasi struct { //lokasi yang bisa melakukan presensi
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Nama     string             `bson:"nama,omitempty"`
	Batas    Geometry           `bson:"batas,omitempty"`
	Kategori string             `bson:"kategori,omitempty"`
}

type Geometry struct { //data geometry untuk lokasi presensi
	Type        string      `json:"type" bson:"type"`
	Coordinates interface{} `json:"coordinates" bson:"coordinates"`
}

type PresensiLokasi struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	PhoneNumber string             `bson:"phonenumber,omitempty"`
	Lokasi      Lokasi             `bson:"lokasi,omitempty"`
	Selfie      bool               `bson:"selfie,omitempty"`
	IsDatang    bool               `bson:"isdatang,omitempty"`
}

type PresensiSelfie struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	CekInLokasi PresensiLokasi     `bson:"cekinlokasi,omitempty"`
	IsDatang    bool               `bson:"isdatang,omitempty"`
}
