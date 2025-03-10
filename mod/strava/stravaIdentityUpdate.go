package strava

import (
	"net/http/cookiejar"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocroot/helper/atdb"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func StravaIdentityUpdateHandler(Pesan itmodel.IteungMessage, db *mongo.Database) string {
	reply := "Informasi Profile Stava kakak: "

	col := "strava_identity"
	// cek apakah akun strava sudah terdaftar di database
	data, err := atdb.GetOneDoc[StravaIdentity](db, col, bson.M{"phone_number": Pesan.Phone_number})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "Kak, akun Strava kamu belum terdaftar. Silakan daftar dulu!"
		}
		return "\n\nError fetching data dari MongoDB: " + err.Error()
	}
	if data.LinkIndentity == "" {
		return "link Strava kamu belum tersimpan di database!"
	}

	jar, _ := cookiejar.New(nil)

	c := colly.NewCollector(
		colly.AllowedDomains(domWeb),
	)

	c.SetCookieJar(jar)

	stravaIdentity := StravaIdentity{}
	stravaIdentity.AthleteId = data.AthleteId

	c.OnHTML("main", func(e *colly.HTMLElement) {
		name := e.ChildText("h2.Details_name__Wz5bH")

		e.ForEach("img", func(_ int, imgEl *colly.HTMLElement) {
			imgTitle := imgEl.Attr("title")
			if imgTitle == name {
				stravaIdentity.Picture = imgEl.Attr("src")
			}
		})
	})

	c.OnScraped(func(r *colly.Response) {
		if data.AthleteId == "" {
			reply += "\n\nAkun Strava kak " + Pesan.Alias_name + " belum terdaftar."
			return
		}

		if data.AthleteId == stravaIdentity.AthleteId {
			// cek apakah data sudah up to date
			if data.Picture == stravaIdentity.Picture {
				reply += "\n\nData Strava kak " + Pesan.Alias_name + " sudah up to date."
				return
			}

			stravaIdentity.UpdatedAt = time.Now()

			updateData := bson.M{
				"picture":    stravaIdentity.Picture,
				"updated_at": stravaIdentity.UpdatedAt,
			}

			// update data ke database jika ada perubahan
			_, err := atdb.UpdateDoc(db, col, bson.M{"athlete_id": stravaIdentity.AthleteId}, bson.M{"$set": updateData})
			if err != nil {
				reply += "\n\nError updating data to MongoDB: " + err.Error()
				return
			}

			reply += "\n\nData kak " + Pesan.Alias_name + " sudah berhasil di update."
			reply += "\n\nUpdate juga Strava Profile Picture kakak di profile akun do.my.id yaa \n" + stravaIdentity.Picture

			token := getCookieFromColly(c, "login")
			if token == "" {
				reply += "\n\nError membaca token dari cookie. " + token

			} else {
				resp := pushToBackend(stravaIdentity.PhoneNumber, stravaIdentity.Picture, token)
				if resp != "" {
					reply += "\n\nError sending data to Backend: " + resp
				} else {
					reply += "\n\nStrava Profile Picture Kak " + Pesan.Alias_name + " sudah berhasil di update."
					reply += "\n\nCek Ulang di do.my.id yaa kak."
				}
			}

		} else {
			reply += "\n\nData Strava kak " + Pesan.Alias_name + " tidak ditemukan."
			return
		}
	})

	err = c.Visit(data.LinkIndentity)
	if err != nil {
		return "Link Profile Strava yang anda kirimkan tidak valid. Silakan kirim ulang dengan link yang valid.(3)"
	}

	return reply
}
