package strava

import (
	"net/http"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocroot/helper/atapi"
	"github.com/gocroot/helper/atdb"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func StravaActivityUpdateIfEmptyDataHandler(Profile itmodel.Profile, Pesan itmodel.IteungMessage, db *mongo.Database) string {
	reply := "Informasi Stava kamu hari ini: "

	c := colly.NewCollector(
		colly.AllowedDomains(domWeb, domApp),
	)

	var fullActivityURL string

	// ambil link strava activity dari pesan
	rawUrl := extractStravaLink(Pesan.Message)
	if rawUrl == "" {
		return reply + "\n\nMaaf, pesan yang kamu kirim tidak mengandung link Strava. Silakan kirim link aktivitas Strava untuk mendapatkan informasinya."
	}

	path := "/activities/"
	if strings.Contains(rawUrl, domWeb) {
		activityId, fullActivityURL = extractContains(rawUrl, path, false)
		if activityId != "" {
			reply += scrapeStravaActivityUpdate(db, fullActivityURL, Profile.Phonenumber, Pesan.Phone_number, Pesan.Alias_name)
		}

	} else if strings.Contains(rawUrl, domApp) {
		c.OnHTML("a", func(e *colly.HTMLElement) {
			link := e.Attr("href")

			activityId, fullActivityURL = extractContains(link, path, true)
			if activityId != "" {
				reply += scrapeStravaActivityUpdate(db, fullActivityURL, Profile.Phonenumber, Pesan.Phone_number, Pesan.Alias_name)
			}
		})
	}

	err := c.Visit(rawUrl)
	if err != nil {
		return "Link Profile Strava yang anda kirimkan tidak valid. Silakan kirim ulang dengan link yang valid.(3) "
	}

	return reply
}

func scrapeStravaActivityUpdate(db *mongo.Database, url, profilePhone, phone, alias string) string {
	reply := ""

	if msg := maintenance(phone); msg != "" {
		reply += msg
		return reply
	}

	c := colly.NewCollector(
		colly.AllowedDomains(domWeb),
	)

	// cek apakah akun strava sudah terdaftar di database
	Idata, err := atdb.GetOneDoc[StravaIdentity](db, "strava_identity", bson.M{"phone_number": phone})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "Kak, akun Strava kamu belum terdaftar. Silakan daftar dulu!"
		}
		return "\n\nError fetching data dari MongoDB: " + err.Error()
	}

	col := "strava_activity"
	filter := bson.M{
		"athlete_id":  Idata.AthleteId,
		"activity_id": activityId,
		"distance":    bson.M{"$eq": ""},
		"moving_time": bson.M{"$eq": ""},
	}
	// cek apakah akun strava sudah terdaftar di database
	data, err := atdb.GetOneDoc[StravaActivity](db, col, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "Kak, Strava Activity kamu tidak ditemukan di database. Silakan cek kembali."
		}
		return "\n\nError fetching data dari MongoDB: " + err.Error()
	}

	stravaActivity := StravaActivity{}
	stravaActivity.ActivityId = data.ActivityId

	c.OnHTML("main", func(e *colly.HTMLElement) {
		stravaActivity.Picture = extractStravaProfileImg(e, stravaActivity.Name)
	})

	c.OnHTML("div", func(e *colly.HTMLElement) {
		extractStravaActivitySpan(e, &stravaActivity)
	})

	found := false
	c.OnHTML("div[class^='MapAndElevationChart_mapContainer__']", func(e *colly.HTMLElement) {
		found = true
	})

	c.OnScraped(func(r *colly.Response) {
		if data.ActivityId != stravaActivity.ActivityId {
			reply += "\n\nStrava Activity ID kak " + alias + " tidak sama."
			return
		}

		if stravaActivity.Distance == "" || stravaActivity.MovingTime == "" {
			reply += "\n\nMaaf kak, kami tidak dapat mengambil data aktivitas kamu. Coba hubungi admin ya."
			return
		}

		if stravaActivity.Picture == "" {
			reply += "\n\nMaaf kak, sistem tidak dapat mengambil foto profil Strava kamu. Pastikan akun Strava kamu dibuat public(everyone). doc: https://www.do.my.id/mentalhealt-strava"
			return
		}

		if Idata.AthleteId != data.AthleteId {
			reply += "\n\nAda yang salah nih dengan akun strava kamu, coba lakukan update dengan perintah dibawah yaaa"
			reply += "\n\n *strava update in*"
			reply += "\n\nAtau mungkin link yang kamu share bukan punya kamu 😏"
			return
		}

		if strings.TrimSpace(data.Distance) != "" && strings.TrimSpace(data.MovingTime) != "" {
			reply += "\n\nData Strava kak " + alias + " sudah up to date."
			return
		}

		// cek apakah ada map atau tidak di halaman strava
		if !found {
			reply += "\n\nJangan Curang donggg! Silahkan share record aktivitas yang benar dari Strava ya kak, bukan dibikin manual kaya gitu"
			reply += "\nYang semangat dong... yang semangat dong..."
			return
		}

		if stravaActivity.TypeSport == "Ride" {
			reply += "\n\nMaaf kak, sistem hanya dapat mengambil data aktivitas jalan dan lari. Silakan share link aktivitas jalan dan lari Strava kamu."
			return
		}

		if data.Distance == "" && data.MovingTime == "" && data.ActivityId == stravaActivity.ActivityId {
			if data.Status == "Invalid" {
				distanceFloat := parseDistance(stravaActivity.Distance)
				if distanceFloat < 3 {
					stravaActivity.UpdatedAt = time.Now()
					stravaActivity.Status = "Invalid"

					updateData := bson.M{
						"distance":    stravaActivity.Distance,
						"moving_time": stravaActivity.MovingTime,
						"elevation":   stravaActivity.Elevation,
						"updated_at":  stravaActivity.UpdatedAt,
						"status":      stravaActivity.Status,
					}

					// update data ke database jika ada perubahan
					_, err := atdb.UpdateDoc(db, col, bson.M{"activity_id": stravaActivity.ActivityId}, bson.M{"$set": updateData})
					if err != nil {
						reply += "\n\nError updating data to MongoDB: " + err.Error()
						return
					}

					reply += "\n\nWahhh, kamu malas sekali ya, jangan malas lari terus dong kak! 😏"
					reply += "\nSatu hari minimal 3 km, masa kamu cuma " + stravaActivity.Distance + " aja"
					reply += "\n\nJangan lupa jaga kesehatan dan tetap semangat!! 💪🏻💪🏻💪🏻"
					return
				} else {
					// simpan data ke database jika data belum ada
					stravaActivity.UpdatedAt = time.Now()
					stravaActivity.Status = "Valid"

					updateData := bson.M{
						"distance":    stravaActivity.Distance,
						"moving_time": stravaActivity.MovingTime,
						"elevation":   stravaActivity.Elevation,
						"updated_at":  stravaActivity.UpdatedAt,
						"status":      stravaActivity.Status,
					}

					// update data ke database jika ada perubahan
					_, err := atdb.UpdateDoc(db, col, bson.M{"activity_id": stravaActivity.ActivityId}, bson.M{"$set": updateData})
					if err != nil {
						reply += "\n\nError updating data to MongoDB: " + err.Error()
						return

					} else {
						reply += "\n\nHaiiiii kak, " + "*" + alias + "*" + "! Berikut Progres Aktivitas kamu hari ini yaaa yang di update!! 😀"
						reply += "\n\n- Activity ID: " + stravaActivity.ActivityId
						reply += "\n- Name: " + data.Name
						reply += "\n- Title: " + data.Title
						reply += "\n- Date Time: " + data.DateTime
						reply += "\n- Type Sport: " + data.TypeSport
						reply += "\n- Distance: " + stravaActivity.Distance
						reply += "\n- Moving Time: " + stravaActivity.MovingTime
						reply += "\n- Elevation: " + stravaActivity.Elevation
						reply += "\n- Status: " + stravaActivity.Status
						reply += "\n\nSemangat terus, jangan lupa jaga kesehatan dan tetap semangat!! 💪🏻💪🏻💪🏻"
					}

					conf, err := atdb.GetOneDoc[Config](db, "config", bson.M{"phonenumber": profilePhone})
					if err != nil {
						reply += "\n\nWah kak " + alias + " mohon maaf ada kesalahan dalam pengambilan config di database " + err.Error()
						return
					}

					datastrava := map[string]interface{}{
						"stravaprofilepicture": stravaActivity.Picture,
						"athleteid":            stravaActivity.AthleteId,
						"phonenumber":          Idata.PhoneNumber,
						"name":                 alias,
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

					reply += "\n\nStrava Profile Picture: " + stravaActivity.Picture
					reply += "\n\nCek link di atas apakah sudah sama dengan Strava Profile Picture di profile akun do.my.id yaa"
				}

			} else {
				reply += "\n\nMaaf kak, Tidak bisa mengambil data aktivitas kamu.(1)"
				return
			}
		} else {
			reply += "\n\nMaaf kak, Tidak bisa mengambil data aktivitas kamu.(2)"
			return
		}
	})

	// rawUrl = data.LinkActivity

	err = c.Visit(url)
	if err != nil {
		return "Link Profile Strava yang anda kirimkan tidak valid. Silakan kirim ulang dengan link yang valid.(3) "
	}

	return reply
}
