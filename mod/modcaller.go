package mod

import (
	"github.com/gocroot/mod/idgrup"
	"github.com/gocroot/mod/presensi"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/mongo"
)

func Caller(Profile itmodel.Profile, Modulename string, Pesan itmodel.IteungMessage, db *mongo.Database) (reply string) {
	switch Modulename {
	case "idgrup":
		reply = idgrup.IDGroup(Pesan)
	case "presensi-masuk":
		reply = presensi.PresensiMasuk(Pesan, db)
	case "presensi-pulang":
		reply = presensi.PresensiPulang(Pesan, db)
	case "selfie-masuk":
		reply = presensi.CekSelfieMasuk(Profile, Pesan, db)
	case "selfie-pulang":
		reply = presensi.CekSelfiePulang(Pesan, db)
	}
	return
}
