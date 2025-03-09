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

var activityId string

func StravaActivityHandler(Pesan itmodel.IteungMessage, db *mongo.Database) string {
	reply := "Informasi Stava kamu hari ini: "

	c := colly.NewCollector(
		colly.AllowedDomains("strava.app.link"),
	)

	c.OnError(func(_ *colly.Response, err error) {
		reply += "Something went wrong:\n\n" + err.Error()
	})

	c.OnHTML("a", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		path := "/activities/"
		if strings.Contains(link, path) {
			parts := strings.Split(link, path)

			if len(parts) > 1 {
				activityId = strings.Split(parts[1], "/")[0]
				activityId = strings.Split(activityId, "?")[0]
				fullActivityURL := "https://www.strava.com" + path + activityId

				reply += scrapeStravaActivity(db, fullActivityURL, Pesan.Phone_number, Pesan.Alias_name)
			}
		}
	})

	rawUrl := extractStravaLink(Pesan.Message)
	if rawUrl == "" {
		return reply + "\n\nMaaf, pesan yang kamu kirim tidak mengandung link Strava. Silakan kirim link aktivitas Strava untuk mendapatkan informasinya."
	}

	err := c.Visit(rawUrl)
	if err != nil {
		return "Link Strava Activity yang anda kirimkan tidak valid. Silakan kirim ulang dengan link yang valid.(1)"
	}

	reply += "\n\nlink strava activity kamu: " + rawUrl

	return reply
}

func scrapeStravaActivity(db *mongo.Database, url, phone, alias string) string {
	reply := ""

	c := colly.NewCollector(
		colly.AllowedDomains("www.strava.com"),
	)

	var activities []StravaActivity

	stravaActivity := StravaActivity{}
	stravaActivity.ActivityId = activityId
	stravaActivity.LinkActivity = url

	c.OnHTML("main", func(e *colly.HTMLElement) {
		stravaActivity.Name = e.ChildText("h3.styles_name__sPSF9")
		stravaActivity.Title = e.ChildText("h1.styles_name__irvsZ")
		stravaActivity.DateTime = e.ChildText("time.styles_date__Bx7mx")
		stravaActivity.TypeSport = e.ChildText("span.styles_typeText__6DEXK")

		e.ForEach("img", func(_ int, imgEl *colly.HTMLElement) {
			imgTitle := imgEl.Attr("title")
			if imgTitle == stravaActivity.Name {
				stravaActivity.Picture = imgEl.Attr("src")
			}
		})

		activities = append(activities, stravaActivity)
	})

	c.OnHTML("div.Stat_stat__hhbSV", func(e *colly.HTMLElement) {
		label := e.ChildText("span.Stat_statLabel__9Qe6h")
		value := e.ChildText("div.Stat_statValue__jbFOA")

		switch strings.ToLower(label) {
		case "distance":
			stravaActivity.Distance = value
		case "time":
			stravaActivity.MovingTime = value
		case "elevation":
			stravaActivity.Elevation = value
		}
	})

	// found := false

	// c.OnHTML("div.MapAndElevationChart_mapContainer__VIs6u", func(e *colly.HTMLElement) {
	// 	found = true
	// })

	c.OnScraped(func(r *colly.Response) {
		// if !found {
		// 	reply += "\n\nJangan Curang donggg! Silahkan share record aktivitas yang benar dari Strava ya kak, bukan dibikin manual kaya gitu"
		// 	reply += "\nYang semangat dong... yang semangat dong..."
		// 	return
		// }

		distanceFloat := parseDistance(stravaActivity.Distance)
		if distanceFloat < 5 {
			reply += "\n\nWahhh, kamu malas sekali ya, jangan malas lari terus dong kak! 😏"
			reply += "\nSatu hari minimal 5 km, masa kamu cuma " + stravaActivity.Distance + " aja 😂 \nxixixixiixi"
			reply += "\n\nJangan lupa jaga kesehatan dan tetap semangat!! 💪🏻💪🏻💪🏻"
			return
		}

		Idata, err := atdb.GetOneDoc[StravaIdentity](db, "strava_identity", bson.M{"phone_number": phone})
		if err != nil && err != mongo.ErrNoDocuments {
			reply += "\n\nError fetching data from MongoDB: " + err.Error()
			return
		}
		if Idata.Picture != stravaActivity.Picture {
			reply += "\n\nAda yang salah nih dengan akun strava kamu, selahkan lakukan update dengan perintah dibawah yaaa"
			reply += "\n\n *strava update in*"
			return
		}

		col := "strava_activity"
		data, err := atdb.GetOneDoc[StravaActivity](db, col, bson.M{"activity_id": stravaActivity.ActivityId})
		if err != nil && err != mongo.ErrNoDocuments {
			reply += "\n\nError fetching data from MongoDB: " + err.Error()
			return
		}
		if data.ActivityId == stravaActivity.ActivityId {
			reply += "\n\n*AOKWOKOKOWKOWKOWKKWOKOK* 🤣🤣"
			reply += "\nHayoolooooo ngapain, Jangan Curang donggg! 😏 Kamu sudah pernah share aktivitas ini sebelumnya."
			reply += "\nSana Lari lagi jangan malas!"
			return
		}

		stravaActivity.CreatedAt = time.Now()

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
	})

	err := c.Visit(url)
	if err != nil {
		return "Link Strava Activity yang anda kirimkan tidak valid. Silakan kirim ulang dengan link yang valid.(2)"
	}

	return reply
}
