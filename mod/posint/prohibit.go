package posint

import (
	"strconv"
	"strings"

	"github.com/gocroot/helper/atdb"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetProhibitedItems(Pesan itmodel.IteungMessage, db *mongo.Database) (reply string) {
	country, err := GetCountryFromMessage(Pesan.Message, db)
	if err != nil {
		return "Terdapat kesalahan pada  GetCountryFromMessage " + err.Error()
	}
	if country == "" {
		return "Nama negara tidak ada kak di database kita"
	}
	listprob, err := atdb.GetAllDoc[[]Item](db, "prohibited_items", bson.M{"Destination": country})
	if err != nil {
		return "Terdapat kesalahan pada  GetAllDoc " + err.Error()
	}
	msg := "ini dia list prohibited item dari negara yang kakak minta:\n"
	for i, probitem := range listprob {
		msg += strconv.Itoa(i+1) + ". " + probitem.ProhibitedItems + "\n"
	}
	return msg

}

func GetCountryFromMessage(message string, db *mongo.Database) (country string, err error) {
	// Ubah pesan menjadi huruf kecil
	lowerMessage := strings.ToLower(message)
	// Mendapatkan nama negara
	countries, err := atdb.GetAllDistinctDoc(db, bson.M{}, "Destination", "prohibited_items")
	if err != nil {
		return
	}
	// Iterasi melalui daftar negara
	for _, country := range countries {
		lowerCountry := strings.ToLower(country.(string))
		if strings.Contains(lowerMessage, lowerCountry) {
			return country.(string), nil
		}
	}
	return "", nil
}
