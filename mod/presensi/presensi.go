package presensi

import (
	"fmt"

	"github.com/gocroot/helper/atdb"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func PresensiMasuk(Pesan itmodel.IteungMessage, db *mongo.Database) (reply string) {
	longitude := fmt.Sprintf("%f", Pesan.Longitude)
	latitude := fmt.Sprintf("%f", Pesan.Latitude)
	lokasiuser, err := GetLokasi(db, Pesan.Longitude, Pesan.Latitude)
	if err != nil {
		return "Mohon maaf kak, kakak belum berada di lokasi presensi, silahkan menuju lokasi presensi dahulu baru cekin masuk."
	}
	if lokasiuser.Nama == "" {
		return "Nama nya kosong kak"
	}
	dtuser := &CekinLokasi{
		PhoneNumber: Pesan.Phone_number,
		Lokasi:      lokasiuser,
	}
	_, err = atdb.InsertOneDoc(db, "presensi", dtuser)
	if err != nil {
		return "Gagal insert ke database kak"
	}

	return "Hai.. hai.. kakak atas nama:\n" + Pesan.Alias_name + "\nLongitude: " + longitude + "\nLatitude: " + latitude + "\nLokasi:" + lokasiuser.Nama + "\nberhasil absen\nmakasih"
}

func GetLokasi(mongoconn *mongo.Database, long float64, lat float64) (lokasi Lokasi, err error) {
	filter := bson.M{
		"batas": bson.M{
			"$geoIntersects": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{long, lat},
				},
			},
		},
	}

	lokasi, err = atdb.GetOneDoc[Lokasi](mongoconn, "lokasi", filter)
	if err != nil {
		return
	}
	return
}
