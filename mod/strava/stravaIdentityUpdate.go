package strava

import (
	"time"

	"github.com/gocolly/colly"
	"github.com/gocroot/helper/atdb"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func StravaIdentityUpdateHandler(Pesan itmodel.IteungMessage, db *mongo.Database) string {
	reply := "Informasi Profile Stava kamu1: "

	col := "strava_identity"
	data, err := atdb.GetOneDoc[StravaIdentity](db, col, bson.M{"phone_number": Pesan.Phone_number})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "Kak, akun Strava kamu belum terdaftar. Silakan daftar dulu!"
		}
		return "\n\nError fetching data dari MongoDB: " + err.Error()
	}

	if data.LinkIndentity == "" {
		return "Kak, link Strava kamu belum tersimpan di database. Silakan tambahkan dulu!"
	}

	c := colly.NewCollector(
		colly.AllowedDomains("www.strava.com"),
	)

	stravaIdentity := StravaIdentity{}
	stravaIdentity.AthleteId = data.AthleteId
	stravaIdentity.PhoneNumber = data.PhoneNumber

	c.OnHTML("main", func(e *colly.HTMLElement) {
		name := e.ChildText("h2.Details_name__Wz5bH")

		e.ForEach("img", func(_ int, imgEl *colly.HTMLElement) {
			imgTitle := imgEl.Attr("title")
			if imgTitle == name {
				stravaIdentity.Picture = imgEl.Attr("src")
			}
		})
	})

	err = c.Visit(data.LinkIndentity)
	if err != nil {
		return "Link Profile Strava yang anda kirimkan tidak valid. Silakan kirim ulang dengan link yang valid.(3)"
	}

	if data.AthleteId == stravaIdentity.AthleteId {
		if data.PhoneNumber != Pesan.Phone_number {
			return "\n\nBukan akun kamu ini mah kak."
		}

		if data.Picture == stravaIdentity.Picture {
			return "\n\nData Strava kak " + Pesan.Alias_name + " sudah up to date." + stravaIdentity.AthleteId
		}

		stravaIdentity.UpdatedAt = time.Now()

		updateData := bson.M{
			"picture":    stravaIdentity.Picture,
			"updated_at": stravaIdentity.UpdatedAt,
		}

		_, err := atdb.UpdateDoc(db, col, bson.M{"athlete_id": stravaIdentity.AthleteId}, bson.M{"$set": updateData})
		if err != nil {
			return "\n\nError updating data to MongoDB: " + err.Error()

		}

		reply += "\n\nData kak " + Pesan.Alias_name + " sudah berhasil di update."

	}

	return reply
}
