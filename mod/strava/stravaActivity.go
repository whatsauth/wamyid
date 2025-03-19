package strava

import (
	"net/http"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocroot/helper/atapi"
	"github.com/gocroot/helper/atdb"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var activityId string

func StravaActivityHandler(Profile itmodel.Profile, Pesan itmodel.IteungMessage, db *mongo.Database) string {
	reply := "Informasi Stava kamu hari ini: "

	var fullActivityURL string

	c := colly.NewCollector(
		colly.AllowedDomains(domApp, domWeb),
	)

	c.OnError(func(_ *colly.Response, err error) {
		reply += "Something went wrong:\n\n" + err.Error()
	})

	// ambil link strava activity dari pesan
	rawUrl := extractStravaLink(Pesan.Message)
	if rawUrl == "" {
		return reply + "\n\nMaaf, pesan yang kamu kirim tidak mengandung link Strava. Silakan kirim link aktivitas Strava untuk mendapatkan informasinya."
	}

	if strings.Contains(rawUrl, domWeb) {
		reply += scrapeStravaActivity(db, rawUrl, Profile.Phonenumber, Pesan.Phone_number, Pesan.Alias_name)

	} else if strings.Contains(rawUrl, domApp) {
		c.OnHTML("a", func(e *colly.HTMLElement) {
			link := e.Attr("href")

			path := "/activities/"
			if strings.Contains(link, path) {
				parts := strings.Split(link, path)

				if len(parts) > 1 {
					activityId = strings.Split(parts[1], "/")[0]
					activityId = strings.Split(activityId, "?")[0]
					fullActivityURL = "https://www.strava.com" + path + activityId

					reply += scrapeStravaActivity(db, fullActivityURL, Profile.Phonenumber, Pesan.Phone_number, Pesan.Alias_name)
				}
			}
		})
	}

	err := c.Visit(rawUrl)
	if err != nil {
		return "Link Strava Activity yang anda kirimkan tidak valid. Silakan kirim ulang dengan link yang valid.(1)"
	}

	if fullActivityURL != "" {
		reply += "\n\nLink Activity Strava kamu: " + fullActivityURL
	} else {
		reply += "\n\nLink Activity Strava kamu: " + rawUrl
	}

	return reply
}

func scrapeStravaActivity(db *mongo.Database, url, profilePhone, phone, alias string) string {
	reply := ""

	if msg := maintenance(phone); msg != "" {
		reply += msg
		return reply
	}

	c := colly.NewCollector(
		colly.AllowedDomains(domWeb),
	)

	var activities []StravaActivity

	stravaActivity := StravaActivity{}
	stravaActivity.ActivityId = activityId
	stravaActivity.LinkActivity = url
	stravaActivity.PhoneNumber = phone

	c.OnHTML("main", func(e *colly.HTMLElement) {
		stravaActivity.Name = e.ChildText("h3[class^='styles_name__']")
		stravaActivity.Title = e.ChildText("h1[class^='styles_name__']")
		stravaActivity.TypeSport = e.ChildText("span[class^='styles_typeText__']")
		// stravaActivity.DateTime = e.ChildText("time.styles_date__Bx7mx")

		e.ForEach("time[class^='styles_date__']", func(_ int, timeEl *colly.HTMLElement) {
			dt := timeEl.Attr("datetime")
			if dt != "" {
				stravaActivity.DateTime = formatDateTimeToIndo(dt)
			} else {
				stravaActivity.DateTime = dt
			}
		})

		e.ForEach("img", func(_ int, imgEl *colly.HTMLElement) {
			imgTitle := imgEl.Attr("title")
			if imgTitle == stravaActivity.Name {
				stravaActivity.Picture = imgEl.Attr("src")
			}
		})

		activities = append(activities, stravaActivity)
	})

	c.OnHTML("div", func(e *colly.HTMLElement) {
		e.ForEach("span", func(_ int, el *colly.HTMLElement) {
			label := strings.ToLower(strings.TrimSpace(el.Text))
			value := strings.TrimSpace(el.DOM.Next().Text()) // Ambil elemen di sebelahnya

			switch label {
			case "distance":
				stravaActivity.Distance = value
			case "time":
				stravaActivity.MovingTime = value
			case "elevation":
				stravaActivity.Elevation = value
			}
		})
	})

	// cek apakah yang share link strava activity adalah pemilik akun strava
	Idata, err := atdb.GetOneDoc[StravaIdentity](db, "strava_identity", bson.M{"phone_number": phone})
	if err != nil && err != mongo.ErrNoDocuments {
		reply += "\n\nError fetching data from MongoDB: " + err.Error()
		return reply
	}

	if Idata.AthleteId != "" {
		stravaActivity.AthleteId = Idata.AthleteId
	} else {
		reply += "\n\nAthlete ID tidak ditemukan. Pastikan akun Strava kamu sudah terhubung dengan sistem."
		return reply
	}

	found := false
	c.OnHTML("div[class^='MapAndElevationChart_mapContainer__']", func(e *colly.HTMLElement) {
		found = true
	})

	c.OnScraped(func(r *colly.Response) {
		if stravaActivity.Distance == "" || stravaActivity.MovingTime == "" {
			reply += "\n\nMaaf kak, kami tidak dapat mengambil data aktivitas kamu. Coba hubungi admin ya."
			return
		}

		if stravaActivity.Picture == "" {
			reply += "\n\nMaaf kak, sistem tidak dapat mengambil foto profil Strava kamu. Pastikan profil dan activity Strava kamu dibuat public(everyone). doc: https://www.do.my.id/mentalhealt-strava"
			return
		}

		// cek apakah data sudah up to date
		if Idata.Picture != stravaActivity.Picture || Idata.Name != stravaActivity.Name {
			reply += "\n\nAda yang salah nih dengan akun strava kamu, coba lakukan update dengan perintah dibawah yaaa"
			reply += "\n\n *strava update in*"
			reply += "\n\nAtau mungkin link yang kamu share bukan punya kamu 😏"
			return
		}

		col := "strava_activity"
		// cek apakah data sudah ada di database
		data, err := atdb.GetOneDoc[StravaActivity](db, col, bson.M{"activity_id": stravaActivity.ActivityId})
		if err != nil && err != mongo.ErrNoDocuments {
			reply += "\n\nError fetching data from MongoDB: " + err.Error()
			return
		}
		if data.ActivityId == stravaActivity.ActivityId {
			createdAtFormated := formatDateTimeToIndo(data.CreatedAt.Format("2006-01-02T15:04:05"))
			reply += "\n\nMaaf kak, " + "*" + alias + "*" + "! Kamu sudah pernah share aktivitas ini sebelumnya pada tanggal " + createdAtFormated + "! Berikut data aktivitas kamu yang sudah tersimpan."
			reply += "\n\n- Name: " + stravaActivity.Name
			reply += "\n- Title: " + stravaActivity.Title
			reply += "\n- Date Time: " + stravaActivity.DateTime
			reply += "\n- Type Sport: " + stravaActivity.TypeSport
			reply += "\n- Distance: " + stravaActivity.Distance
			reply += "\n- Moving Time: " + stravaActivity.MovingTime
			reply += "\n- Elevation: " + stravaActivity.Elevation
			reply += "\n\nSemangat terus, jangan lupa jaga kesehatan dan tetap semangat!! 💪🏻💪🏻💪🏻"

			reply += "\nJangan Curang yaaaa!"
			reply += "\nSana Lari lagi jangan malas!"

			// simpan data ke database jika data belum ada
			stravaActivity.CreatedAt = time.Now()
			stravaActivity.Status = "Duplicate"

			_, err = atdb.InsertOneDoc(db, col, stravaActivity)
			if err != nil {
				reply += "\n\nError saving data to MongoDB: " + err.Error()
			}
			return
		}

		// cek apakah ada map atau tidak di halaman strava
		if !found {
			// simpan data ke database jika data belum ada
			stravaActivity.CreatedAt = time.Now()
			stravaActivity.Status = "Fraudulent"

			_, err = atdb.InsertOneDoc(db, col, stravaActivity)
			if err != nil {
				reply += "\n\nError saving data to MongoDB: " + err.Error()
			}

			reply += "\n\nJangan Curang donggg! Silahkan share record aktivitas yang benar dari Strava ya kak, bukan dibikin manual kaya gitu"
			reply += "\nYang semangat dong... yang semangat dong..."
			return
		}

		if stravaActivity.TypeSport == "Ride" {
			// simpan data ke database jika data belum ada
			stravaActivity.CreatedAt = time.Now()
			stravaActivity.Status = "Invalid"

			_, err = atdb.InsertOneDoc(db, col, stravaActivity)
			if err != nil {
				reply += "\n\nError saving data to MongoDB: " + err.Error()
			}

			reply += "\n\nMaaf kak, sistem hanya dapat mengambil data aktivitas jalan dan lari. Silakan share link aktivitas jalan dan lari Strava kamu."
			return
		}

		// cek apakah jarak lari kurang dari 5 km
		distanceFloat := parseDistance(stravaActivity.Distance)
		if distanceFloat < 3 {
			// simpan data ke database jika data belum ada
			stravaActivity.CreatedAt = time.Now()
			stravaActivity.Status = "Invalid"

			_, err = atdb.InsertOneDoc(db, col, stravaActivity)
			if err != nil {
				reply += "\n\nError saving data to MongoDB: " + err.Error()
			}

			reply += "\n\nWahhh, kamu malas sekali ya, jangan malas lari terus dong kak! 😏"
			reply += "\nSatu hari minimal 3 km, masa kamu cuma " + stravaActivity.Distance + " aja 😂 \nxixixixiixi"
			reply += "\n\nJangan lupa jaga kesehatan dan tetap semangat!! 💪🏻💪🏻💪🏻"
			return

		} else {
			// simpan data ke database jika data belum ada
			stravaActivity.CreatedAt = time.Now()
			stravaActivity.Status = "Valid"

			_, err = atdb.InsertOneDoc(db, col, stravaActivity)
			if err != nil {
				reply += "\n\nError saving data to MongoDB: " + err.Error()
			} else {
				reply += "\n\nHaiiiii kak, " + "*" + alias + "*" + "! Berikut Progres Aktivitas kamu hari ini yaaa!! 😀"
				reply += "\n\n- Name: " + stravaActivity.Name
				reply += "\n- Title: " + stravaActivity.Title
				reply += "\n- Date Time: " + stravaActivity.DateTime
				reply += "\n- Type Sport: " + stravaActivity.TypeSport
				reply += "\n- Distance: " + stravaActivity.Distance
				reply += "\n- Moving Time: " + stravaActivity.MovingTime
				reply += "\n- Elevation: " + stravaActivity.Elevation
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
	})

	err = c.Visit(url)
	if err != nil {
		return "Link Strava Activity yang anda kirimkan tidak valid. Silakan kirim ulang dengan link yang valid.(2)"
	}

	return reply
}
