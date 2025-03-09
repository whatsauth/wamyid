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
	reply := "Informasi Profile Stava kamu: "

	c := colly.NewCollector(
		colly.AllowedDomains("www.strava.com"),
	)

	var identities []StravaIdentity

	stravaIdentity := StravaIdentity{}
	stravaIdentity.AthleteId = athleteId

	c.OnHTML("main", func(e *colly.HTMLElement) {
		name := e.ChildText("h2.Details_name__Wz5bH")

		e.ForEach("img", func(_ int, imgEl *colly.HTMLElement) {
			imgTitle := imgEl.Attr("title")
			if imgTitle == name {
				stravaIdentity.Picture = imgEl.Attr("src")
			}
		})

		identities = append(identities, stravaIdentity)
	})

	col := "strava_identity"
	data, err := atdb.GetOneDoc[StravaIdentity](db, col, bson.M{"athlete_id": stravaIdentity.AthleteId})
	if err != nil && err != mongo.ErrNoDocuments {
		return "\n\nError fetching data from MongoDB: " + err.Error()
	}

	c.OnScraped(func(r *colly.Response) {
		if data.AthleteId == stravaIdentity.AthleteId {
			if data.PhoneNumber != Pesan.Phone_number {
				reply += "\n\nBukan akun kamu ini mah kak."
				return
			}

			if data.Picture != stravaIdentity.Picture {
				stravaIdentity.UpdatedAt = time.Now()

				updateData := bson.M{
					"picture":    stravaIdentity.Picture,
					"updated_at": stravaIdentity.UpdatedAt,
				}

				_, err := atdb.UpdateDoc(db, col, bson.M{"athlete_id": stravaIdentity.AthleteId}, bson.M{"$set": updateData})
				if err != nil {
					reply += "\n\nError updating data to MongoDB: " + err.Error()
					return
				}

				reply += "\n\nData kak " + Pesan.Alias_name + " sudah berhasil di update."
			}
		}
	})

	err = c.Visit(data.LinkIndentity)
	if err != nil {
		return "Link Profile Strava yang anda kirimkan tidak valid. Silakan kirim ulang dengan link yang valid.(3)"
	}

	return reply
}
