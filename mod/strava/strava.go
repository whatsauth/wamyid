package strava

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gocroot/helper/atdb"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var activityId string

func StravaHandler(Pesan itmodel.IteungMessage, db *mongo.Database) string {
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
				fullActivityURL := "https://www.strava.com" + path + activityId

				reply += scrapeStravaActivity(db, fullActivityURL)
			}
		}
	})

	rawUrl := extractStravaLink(Pesan.Message)
	if rawUrl == "" {
		return reply + "\n\nMaaf, pesan yang kamu kirim tidak mengandung link Strava. " +
			"Silakan kirim link aktivitas Strava untuk mendapatkan informasinya. 😊"
	}

	err := c.Visit(rawUrl)
	if err != nil {
		return "\nError visiting URL1" + err.Error()
	}

	return reply + "\n\nlink strava activity kamu: " + rawUrl
}

func scrapeStravaActivity(db *mongo.Database, url string) string {
	reply := ""

	c := colly.NewCollector(
		colly.AllowedDomains("www.strava.com"),
	)

	var activities []StravaActivity

	stravaActivity := StravaActivity{
		ActivityId: activityId,
	}

	c.OnHTML("main", func(e *colly.HTMLElement) {
		stravaActivity.Name = e.ChildText("h3.styles_name__sPSF9")
		stravaActivity.Title = e.ChildText("h1.styles_name__irvsZ")
		stravaActivity.DateTime = e.ChildText("time.styles_date__Bx7mx")
		stravaActivity.TypeSport = e.ChildText("span.styles_typeText__6DEXK")

		activities = append(activities, stravaActivity)
	})

	c.OnHTML("div.Stat_stat__hhbSV", func(e *colly.HTMLElement) {
		label := e.ChildText("span.Stat_statLabel__9Qe6h")
		value := e.ChildText("div.Stat_statValue__jbFOA")

		switch strings.ToLower(label) {
		case "distance":
			stravaActivity.Distance = value
		case "time":
			stravaActivity.TimePeriod = value
		case "elevation":
			stravaActivity.Elevation = value
		}
	})

	c.OnScraped(func(r *colly.Response) {
		distanceFloat := parseDistance(stravaActivity.Distance)
		if distanceFloat < 5 {
			reply += "\n\nWahhh, kamu malas sekali ya, jangan malas lari terus dong kak! 😏" +
				"\nSatu hari minimal 5 km, masa kamu cuma " + stravaActivity.Distance + " aja 😂 \nxixixixiixi" +
				"\n\nJangan lupa jaga kesehatan dan tetap semangat!! 💪🏻💪🏻💪🏻"
			return
		}

		col := "strava"
		data, err := atdb.GetOneDoc[StravaActivity](db, col, bson.M{"activity_id": stravaActivity.ActivityId})
		if err != nil && err != mongo.ErrNoDocuments {
			reply += "\n\nError fetching data from MongoDB: " + err.Error()
			return
		}
		if data.ActivityId == stravaActivity.ActivityId {
			reply += "\n\nHayoolooooo ngapain, Jangan Curang donggg! 😏 Kamu sudah pernah share aktivitas ini sebelumnya." +
				"\n*AOKWOKOKOWKOWKOWKKWOKOK* 🤣🤣" +
				"\nSana Lari lagi jangan malas!"
			return
		}

		_, err = atdb.InsertOneDoc(db, col, stravaActivity)
		if err != nil {
			reply += "\n\nError saving data to MongoDB: " + err.Error()
		} else {
			reply += "\n\nHaiiiii kak, " + stravaActivity.Name + "! Berikut Progres Aktivitas kamu hari ini yaaa!! 😀" +
				"\n\n- Activity_id: " + stravaActivity.ActivityId +
				"\n- Name: " + stravaActivity.Name +
				"\n- Title: " + stravaActivity.Title +
				"\n- Date Time: " + stravaActivity.DateTime +
				"\n- Type Sport: " + stravaActivity.TypeSport +
				"\n- Distance: " + stravaActivity.Distance +
				"\n- Time Period: " + stravaActivity.TimePeriod +
				"\n- Elevation: " + stravaActivity.Elevation +
				"\n\nSemangat terus, jangan lupa jaga kesehatan dan tetap semangat!! 💪🏻💪🏻💪🏻"
		}
	})

	err := c.Visit(url)
	if err != nil {
		return "\nError visiting URL2" + err.Error()
	}

	return reply
}

func extractStravaLink(text string) string {
	re := regexp.MustCompile(`https?://strava\.app\.link/\S+`)
	match := re.FindString(text)

	return match
}

func parseDistance(distance string) float64 {
	reply := ""
	distance = strings.TrimSpace(distance)
	if len(distance) == 0 {
		return 0
	}

	re := regexp.MustCompile(`[0-9]+(\.[0-9]+)?`)
	number := re.FindString(distance)

	if number == "" {
		return 0
	}

	distanceFloat, err := strconv.ParseFloat(number, 64)
	if err != nil {
		reply += "\nError parsing distance: " + err.Error()
		return 0
	}

	return distanceFloat
}
