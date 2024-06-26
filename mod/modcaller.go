package mod

import (
	"github.com/gocroot/mod/idgrup"
	"github.com/gocroot/mod/kyc"
	"github.com/gocroot/mod/lmsdesa"
	"github.com/gocroot/mod/presensi"
	"github.com/gocroot/mod/tasklist"
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
	case "upload-lmsdesa-file":
		reply = lmsdesa.ArsipFile(Pesan, db)
	case "upload-lmsdesa-gambar":
		reply = lmsdesa.ArsipGambar(Pesan, db)
	case "cek-ktp":
		reply = kyc.CekKTP(Profile, Pesan, db)
	case "selfie-masuk":
		reply = presensi.CekSelfieMasuk(Profile, Pesan, db)
	case "selfie-pulang":
		reply = presensi.CekSelfiePulang(Pesan, db)
	case "tasklist-append":
		reply = tasklist.TaskListAppend(Pesan, db)
	case "tasklist-reset":
		reply = tasklist.TaskListReset(Pesan, db)
	case "tasklist-save":
		reply = tasklist.TaskListSave(Pesan, db)
	}
	return
}
