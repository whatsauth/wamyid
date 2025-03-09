package strava

import (
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocroot/helper/atdb"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var athleteId string

func StravaIdentityHandler(Pesan itmodel.IteungMessage, db *mongo.Database) string {
	reply := "Informasi Profile Stava kamu: "

	c := colly.NewCollector(
		colly.AllowedDomains("strava.app.link"),
	)

	c.OnError(func(_ *colly.Response, err error) {
		reply += "Something went wrong:\n\n" + err.Error()
	})

	c.OnHTML("a", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		path := "/athletes/"
		if strings.Contains(link, path) {
			parts := strings.Split(link, path)

			if len(parts) > 1 {
				athleteId = strings.Split(parts[1], "/")[0]
				athleteId = strings.Split(athleteId, "?")[0]
				fullAthleteURL := "https://www.strava.com" + path + athleteId

				reply += scrapeStravaIdentity(db, fullAthleteURL, Pesan.Phone_number, Pesan.Alias_name)
			}
		}
	})

	rawUrl := extractStravaLink(Pesan.Message)
	if rawUrl == "" {
		return reply + "\n\nMaaf, pesan yang kamu kirim tidak mengandung link Strava. Silakan kirim link aktivitas Strava untuk mendapatkan informasinya."
	}

	err := c.Visit(rawUrl)
	if err != nil {
		return "Link Profile Strava yang anda kirimkan tidak valid. Silakan kirim ulang dengan link yang valid.(1)"
	}

	reply += "\n\nlink Profile Strava kamu: " + rawUrl

	return reply
}

func scrapeStravaIdentity(db *mongo.Database, url, phone, alias string) string {
	reply := ""

	c := colly.NewCollector(
		colly.AllowedDomains("www.strava.com"),
	)

	var identities []StravaIdentity

	stravaIdentity := StravaIdentity{}
	stravaIdentity.AthleteId = athleteId
	stravaIdentity.LinkIndentity = url
	stravaIdentity.PhoneNumber = phone

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

	c.OnScraped(func(r *colly.Response) {
		col := "strava_identity"
		// cek apakah data sudah ada di database
		data, err := atdb.GetOneDoc[StravaIdentity](db, col, bson.M{"athlete_id": stravaIdentity.AthleteId})
		if err != nil && err != mongo.ErrNoDocuments {
			reply += "\n\nError fetching data from MongoDB: " + err.Error()
			return
		}
		if data.AthleteId == stravaIdentity.AthleteId {
			reply += "\n\nLink Profile Strava kamu sudah pernah di share sebelumnya."
			return
		}

		stravaIdentity.CreatedAt = time.Now()

		// simpan data ke database jika data belum ada
		_, err = atdb.InsertOneDoc(db, col, stravaIdentity)
		if err != nil {
			reply += "\n\nError saving data to MongoDB: " + err.Error()
		} else {
			reply += "\n\nData Strava Kak " + alias + " sudah berhasil di simpan."
			reply += "\n\nTambahin Strava Profile Picture kamu ke profile akun do.my.id kamu yaa \n" + data.Picture
		}
	})

	err := c.Visit(url)
	if err != nil {
		return "Link Profile Strava yang anda kirimkan tidak valid. Silakan kirim ulang dengan link yang valid.(2)"
	}

	return reply
}
