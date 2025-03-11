package strava

import (
	"net/http"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocroot/helper/atapi"
	"github.com/gocroot/helper/atdb"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func StravaIdentityUpdateHandler(Profile itmodel.Profile, Pesan itmodel.IteungMessage, db *mongo.Database) string {
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

	c := colly.NewCollector(
		colly.AllowedDomains(domWeb),
	)

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

		if stravaIdentity.Picture == "" {
			reply += "\n\nMaaf kak, sistem tidak dapat mengambil foto profil Strava kamu. Pastikan akun Strava kamu dibuat public(everyone). doc: https://www.do.my.id/mentalhealt-strava"
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

			conf, err := atdb.GetOneDoc[Config](db, "config", bson.M{"phonenumber": Profile.Phonenumber})
			if err != nil {
				reply += "\n\nWah kak " + Pesan.Alias_name + " mohon maaf ada kesalahan dalam pengambilan config di database " + err.Error()
				return
			}

			type DataStrava struct {
				StravaProfilePicture string `json:"stravaprofilepicture"`
				PhoneNumber          string `json:"phonenumber"`
			}

			datastrava := DataStrava{
				StravaProfilePicture: stravaIdentity.Picture,
				PhoneNumber:          Pesan.Phone_number,
			}

			statuscode, httpresp, err := atapi.PostStructWithToken[itmodel.Response]("secret", conf.DomyikadoSecret, datastrava, conf.DomyikadoUserURL)
			if err != nil {
				reply += "\n\nAkses ke endpoint domyikado gagal: " + err.Error()
				return
			}

			if statuscode != http.StatusOK {
				reply += "\n\nSalah posting endpoint domyikado: " + httpresp.Response + "\ninfo\n" + httpresp.Info
				return
			}

			reply += "\n\nUpdate Strava Profile Picture berhasil dilakukan di do.my.id, silahkan cek di profile akun do.my.id kakak."

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
