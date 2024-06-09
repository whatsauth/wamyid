package mod

import (
	"github.com/gocroot/mod/idgrup"
	"github.com/gocroot/mod/presensi"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/mongo"
)

func Caller(Modulename string, Pesan itmodel.IteungMessage, db *mongo.Database) (reply string) {
	switch Modulename {
	case "idgrup":
		reply = idgrup.IDGroup(Pesan)
	case "presensi-masuk":
		reply = presensi.PresensiMasuk(Pesan, db)
	case "selfie-masuk":
		reply = presensi.CekSelfieMasuk(Pesan, db)
	}
	return
}
