package presensi

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gocroot/helper/atapi"
	"github.com/gocroot/helper/atdb"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CekSelfiePulang(Pesan itmodel.IteungMessage, db *mongo.Database) (reply string) {
	if Pesan.Filedata == "" {
		return "Kirim pap nya dulu dong kak.."
	}
	dt := FaceDetect{
		IDUser:    Pesan.Phone_number,
		Base64Str: Pesan.Filedata,
	}
	filter := bson.M{"_id": atdb.TodayFilter(), "phonenumber": Pesan.Phone_number, "ismasuk": false}
	pstoday, err := atdb.GetOneDoc[PresensiLokasi](db, "presensi", filter)
	if err != nil {
		return "Wah kak mohon maaf kakak belum cekin share live location hari ini " + err.Error()
	}
	conf, err := atdb.GetOneDoc[Config](db, "config", bson.M{"phonenumber": "62895601060000"})
	if err != nil {
		return "Wah kak mohon maaf ada kesalahan dalam pengambilan config di database " + err.Error()
	}
	statuscode, faceinfo, err := atapi.PostStructWithToken[FaceInfo]("secret", conf.LeaflySecret, dt, conf.LeaflyURL)
	if err != nil {
		return "Wah kak mohon maaf ada kesalahan pemanggilan API leafly " + err.Error()
	}
	if statuscode != http.StatusOK {
		return "Wah kak mohon maaf: " + faceinfo.Error
	}
	pselfie := PresensiSelfie{
		CekInLokasi: pstoday,
		IsMasuk:     true,
		IDUser:      faceinfo.PhoneNumber,
		Commit:      faceinfo.Commit,
		Filehash:    faceinfo.FileHash,
		Remaining:   faceinfo.Remaining,
	}
	_, err = atdb.InsertOneDoc(db, "selfie", pselfie)
	if err != nil {
		return "Wah kak mohon maaf ada kesalahan input ke database " + err.Error()
	}
	return

}

func CekSelfieMasuk(Pesan itmodel.IteungMessage, db *mongo.Database) (reply string) {
	if Pesan.Filedata == "" {
		return "Kirim pap nya dulu dong kak.."
	}
	dt := FaceDetect{
		IDUser:    Pesan.Phone_number,
		Base64Str: Pesan.Filedata,
	}
	filter := bson.M{"_id": atdb.TodayFilter(), "phonenumber": Pesan.Phone_number, "ismasuk": true}
	pstoday, err := atdb.GetOneDoc[PresensiLokasi](db, "presensi", filter)
	if err != nil {
		return "Wah kak mohon maaf kakak belum cekin share live location hari ini " + err.Error()
	}
	conf, err := atdb.GetOneDoc[Config](db, "config", bson.M{"phonenumber": "62895601060000"})
	if err != nil {
		return "Wah kak mohon maaf ada kesalahan dalam pengambilan config di database " + err.Error()
	}
	statuscode, faceinfo, err := atapi.PostStructWithToken[FaceInfo]("secret", conf.LeaflySecret, dt, conf.LeaflyURL)
	if err != nil {
		return "Wah kak mohon maaf ada kesalahan pemanggilan API leafly " + err.Error()
	}
	if statuscode != http.StatusOK {
		return "Wah kak mohon maaf: " + faceinfo.Error
	}
	pselfie := PresensiSelfie{
		CekInLokasi: pstoday,
		IsMasuk:     true,
		IDUser:      faceinfo.PhoneNumber,
		Commit:      faceinfo.Commit,
		Filehash:    faceinfo.FileHash,
		Remaining:   faceinfo.Remaining,
	}
	_, err = atdb.InsertOneDoc(db, "selfie", pselfie)
	if err != nil {
		return "Wah kak mohon maaf ada kesalahan input ke database " + err.Error()
	}
	return

}

func PresensiMasuk(Pesan itmodel.IteungMessage, db *mongo.Database) (reply string) {
	if !Pesan.LiveLoc {
		return "Minimal share live location dulu lah kak."
	}
	longitude := fmt.Sprintf("%f", Pesan.Longitude)
	latitude := fmt.Sprintf("%f", Pesan.Latitude)
	lokasiuser, err := GetLokasi(db, Pesan.Longitude, Pesan.Latitude)
	if err != nil {
		return "Mohon maaf kak, kakak belum berada di lokasi presensi, silahkan menuju lokasi presensi dahulu baru cekin masuk."
	}
	if lokasiuser.Nama == "" {
		return "Nama nya kosong kak"
	}
	dtuser := &PresensiLokasi{
		PhoneNumber: Pesan.Phone_number,
		Lokasi:      lokasiuser,
		IsMasuk:     true,
		CreatedAt:   time.Now(),
	}
	_, err = atdb.InsertOneDoc(db, "presensi", dtuser)
	if err != nil {
		return "Gagal insert ke database kak"
	}

	return "Hai.. hai.. kakak atas nama:\n" + Pesan.Alias_name + "\nLongitude: " + longitude + "\nLatitude: " + latitude + "\nLokasi:" + lokasiuser.Nama + "\nberhasil absen\nmakasih"
}

func PresensiPulang(Pesan itmodel.IteungMessage, db *mongo.Database) (reply string) {
	if !Pesan.LiveLoc {
		return "Minimal share live location dulu lah kak."
	}
	longitude := fmt.Sprintf("%f", Pesan.Longitude)
	latitude := fmt.Sprintf("%f", Pesan.Latitude)
	lokasiuser, err := GetLokasi(db, Pesan.Longitude, Pesan.Latitude)
	if err != nil {
		return "Mohon maaf kak, kakak belum berada di lokasi presensi, silahkan menuju lokasi presensi dahulu baru cekin masuk."
	}
	if lokasiuser.Nama == "" {
		return "Nama nya kosong kak"
	}
	dtuser := &PresensiLokasi{
		PhoneNumber: Pesan.Phone_number,
		Lokasi:      lokasiuser,
		IsMasuk:     false,
		CreatedAt:   time.Now(),
	}
	filter := bson.M{"_id": atdb.TodayFilter(), "iduser": Pesan.Phone_number, "ismasuk": true}
	docselfie, err := atdb.GetOneLatestDoc[PresensiSelfie](db, "selfie", filter)
	if err != nil {
		return "Kakak belum selfie masuk ini " + err.Error()
	}
	if docselfie.CekInLokasi.Lokasi.ID != lokasiuser.ID {
		return "Lokasi pulang nya harus sama dengan lokasi masuknya kak: " + lokasiuser.Nama
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
