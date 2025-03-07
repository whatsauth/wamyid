package mod

import (
	"github.com/gocroot/mod/daftar"
	"github.com/gocroot/mod/helpdesk"
	"github.com/gocroot/mod/idgrup"
	"github.com/gocroot/mod/kyc"
	"github.com/gocroot/mod/lms"
	"github.com/gocroot/mod/lmsdesa"
	"github.com/gocroot/mod/pomodoro"
	"github.com/gocroot/mod/posint"
	"github.com/gocroot/mod/presensi"
	"github.com/gocroot/mod/siakad"
	"github.com/gocroot/mod/tasklist"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/mongo"
)

func Caller(Profile itmodel.Profile, Modulename string, Pesan itmodel.IteungMessage, db *mongo.Database) (reply string) {
	switch Modulename {
	case "idgrup":
		reply = idgrup.IDGroup(Pesan)
	case "feedbackhelpdesk":
		reply = helpdesk.FeedbackHelpdesk(Profile, Pesan, db)
	case "endhelpdesk":
		reply = helpdesk.EndHelpdesk(Profile, Pesan, db)
	case "helpdesk":
		reply = helpdesk.StartHelpdesk(Pesan, db)
	case "presensi-masuk":
		reply = presensi.PresensiMasuk(Pesan, db)
	case "presensi-pulang":
		reply = presensi.PresensiPulang(Pesan, db)
	case "upload-lmsdesa-file":
		reply = lmsdesa.ArsipFile(Pesan, db)
	case "upload-lmsdesa-gambar":
		reply = lmsdesa.ArsipGambar(Pesan, db)
	case "lms":
		reply = lms.ReplyRekapUsers(Profile, Pesan, db)
	case "cek-ktp":
		reply = kyc.CekKTP(Profile, Pesan, db)
	case "selfie-masuk":
		reply = presensi.CekSelfieMasuk(Profile, Pesan, db)
	case "selfie-pulang":
		reply = presensi.CekSelfiePulang(Profile, Pesan, db)
	case "tasklist-append":
		reply = tasklist.TaskListAppend(Pesan, db)
	case "tasklist-reset":
		reply = tasklist.TaskListReset(Pesan, db)
	case "tasklist-save":
		reply = tasklist.TaskListSave(Pesan, db)
	case "domyikado-user":
		reply = daftar.DaftarDomyikado(Profile, Pesan, db)
	case "panduan-siakad":
		reply = siakad.PanduanDosen(Pesan, db)
	case "login-siakad":
		reply = siakad.LoginSiakad(Pesan, db)
	case "approve-bap":
		reply = siakad.ApproveBAP(Pesan, db)
	case "cek-approval":
		reply = siakad.CekApprovalBAP(Pesan, db)
	case "cetak-bap":
		reply = siakad.CetakBAP(Pesan, db)
	case "approve-bimbingan":
		reply = siakad.ApproveBimbingan(Pesan, db)
	case "approve-bimbinganbypoin":
		reply = siakad.ApproveBimbinganbyPoin(Pesan, db)
	case "prohibited-items":
		reply = posint.GetProhibitedItems(Pesan, db)
	case "pomodoro-cycle":
		reply = pomodoro.HandlePomodoroStart(Pesan, db)
	}

	return
}
