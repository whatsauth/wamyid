package helpdesk

import "github.com/whatsauth/itmodel"

func PenugasanOperator(Pesan itmodel.IteungMessage) (reply string) {

	return "Selamat datang di HelpDesk Pamong Desa :\n" + Pesan.Group_id + "\nAnda kami hubungkan dengan operator kami Asep di nomor wa.me/62817898989\nTerima kasih"
}
