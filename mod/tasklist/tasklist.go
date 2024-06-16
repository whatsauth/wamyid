package tasklist

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gocroot/helper/atapi"
	"github.com/gocroot/helper/atdb"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TaskListSave(Pesan itmodel.IteungMessage, db *mongo.Database) (reply string) {
	id, _ := GetIDandTask(Pesan.Message)
	idp, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "gagal mendapatkan id laporan kak"
	}
	taskall, err := atdb.GetAllDoc[[]TaskList](db, "tasklist", bson.M{"laporanid": idp})
	if err != nil {
		return "Data task tidak ditemukan kak"
	}
	conf, err := atdb.GetOneDoc[Config](db, "config", bson.M{"phonenumber": "62895601060000"})
	if err != nil {
		return "Wah kak " + Pesan.Alias_name + " mohon maaf ada kesalahan dalam pengambilan config di database " + err.Error()
	}
	statuscode, httpresp, err := atapi.PostStructWithToken[itmodel.Response]("secret", conf.DomyikadoSecret, taskall, conf.DomyikadoTaskListURL)
	if err != nil {
		return "Akses ke endpoint domyikado gagal: " + err.Error()
	}
	if statuscode != http.StatusOK {
		return "Salah posting endpoint domyikado: " + httpresp.Response + "\ninfo\n" + httpresp.Info
	}
	msg := "Pertemuan " + id + "\nTask Lisk:\n"
	// Loop melalui slice menggunakan range tanpa indeks
	for _, taskone := range taskall {
		msg += taskone.Task + "\n"
	}
	msg += "\n*Sudah disimpan permanen*"
	return msg
}

// input := "https://wa.me/62895800006000?text=-.-T@$kl1$t-.-98suf8usdf0s98dfoi0sid9f|||++task list pertama ini"
func TaskListAppend(Pesan itmodel.IteungMessage, db *mongo.Database) (reply string) {
	id, task := GetIDandTask(Pesan.Message)
	idp, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "gagal mendapatkan id laporan kak"
	}
	dt := TaskList{
		LaporanID:   idp,
		Task:        task,
		PhoneNumber: Pesan.Phone_number,
		Name:        Pesan.Alias_name,
	}
	_, err = atdb.InsertOneDoc(db, "tasklist", dt)
	if err != nil {
		return "gagal insert db kak"
	}
	taskall, err := atdb.GetAllDoc[[]TaskList](db, "tasklist", bson.M{"laporanid": idp})
	if err != nil {
		return "Data task tidak ditemukan kak"
	}
	msg := "Pertemuan " + id + "\nTask Lisk:\n"
	// Loop melalui slice menggunakan range tanpa indeks
	for _, taskone := range taskall {
		msg += taskone.Task + "\n"
	}
	msg += "Untuk menambah task klik:\n" + "https://wa.me/62895800006000?text=-.-T@$kl1$t-.-98suf8usdf0s98dfoi0sid9f|||++" + "\nUntuk Reset Isi Task klik:\n" + "https://wa.me/62895800006000?text=-.-T@$kl1$tR35t-.-98suf8usdf0s98dfoi0sid9f|||++" + "\nUntuk simpan permanen klik:\n" + "https://wa.me/62895800006000?text=-.-T@$kl1$tS@v3-.-98suf8usdf0s98dfoi0sid9f|||++"
	return msg
}

func TaskListReset(Pesan itmodel.IteungMessage, db *mongo.Database) (reply string) {
	id, _ := GetIDandTask(Pesan.Message)
	idp, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "gagal mendapatkan id laporan kak"
	}
	_, err = atdb.DeleteManyDocs(db, "tasklist", bson.M{"laporanid": idp})
	if err != nil {
		return "gagal hapus db kak"
	}
	msg := "Pertemuan " + id + "\nTask Lisk:0\n"
	msg += "Untuk menambah task klik:\n" + "https://wa.me/62895800006000?text=-.-T@$kl1$t-.-98suf8usdf0s98dfoi0sid9f|||++" + "\nUntuk Reset Isi Task klik:\n" + "https://wa.me/62895800006000?text=-.-T@$kl1$tR35t-.-98suf8usdf0s98dfoi0sid9f|||++" + "\nUntuk simpan permanen klik:\n" + "https://wa.me/62895800006000?text=-.-T@$kl1$tS@v3-.-98suf8usdf0s98dfoi0sid9f|||++"
	return msg
}

func GetIDandTask(input string) (cleanedStrBefore, cleanedStrAfter string) {
	// input := "&&T@$kl1$t&&98suf8usdf0s98dfoi0sid9f|||\ntask list pertama ini"

	// Find the position of the delimiter "|||"
	pos := strings.Index(input, "|||")
	if pos == -1 {
		fmt.Println("Delimiter not found")
		return
	}

	// Extract the substring after the delimiter "|||"
	substrAfter := input[pos+len("|||"):]

	// Remove newline characters from the substring after "|||"
	cleanedStrAfter = strings.ReplaceAll(substrAfter, "\n", "")

	// Extract the substring before the delimiter "|||"
	substrBefore := input[:pos]

	// Find the position of the last occurrence of "&&"
	posLastAnd := strings.LastIndex(substrBefore, "-.-")
	if posLastAnd == -1 {
		fmt.Println("Delimiter '&&' not found")
		return
	}

	// Extract the part after the last "&&"
	cleanedStrBefore = substrBefore[posLastAnd+len("-.-"):]

	return
}
