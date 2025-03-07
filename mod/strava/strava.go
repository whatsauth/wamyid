package strava

import (
	"regexp"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gocroot/helper/atdb"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var activityId string

func StravaHandler(Pesan itmodel.IteungMessage, db *mongo.Database) string {
	reply := "Strava activity has been scraped"

	c := colly.NewCollector(
		colly.AllowedDomains("strava.app.link"),
	)

	c.OnError(func(_ *colly.Response, err error) {
		reply += "\nSomething went wrong: " + err.Error()
	})

	c.OnHTML("a", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		if strings.Contains(link, "/activities/") {
			parts := strings.Split(link, "/activities/")

			if len(parts) > 1 {
				activityId = strings.Split(parts[1], "/")[0]
				fullActivityURL := "https://www.strava.com/activities/" + activityId

				reply += scrapeStravaActivity(db, fullActivityURL)
			}
		}
	})

	rawUrl := extractStravaLink(Pesan.Message)
	if rawUrl == "" {
		return reply + " activity link not found"
	}

	err := c.Visit(rawUrl)
	if err != nil {
		return "\nError visiting URL1" + err.Error()
	}

	return reply + "link strava activity kamu: " + rawUrl + "\n\n#mental_health"
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
		if int(stravaActivity.Distance[0]) < 5 {
			reply += "\n\nWahhh, kamu malas sekali ya, jangan malas lari terus dong kak! 😏" +
				"\nSatu hari minimal 5 km, masa kamu cuma " + stravaActivity.Distance + " aja 😂 \nxixixixiixi" +
				"\n\nJangan lupa jaga kesehatan dan tetap semangat!! 💪🏻💪🏻💪🏻"
		}

		col := "strava"
		data, err := atdb.GetOneDoc[StravaActivity](db, col, bson.M{"activity_id": stravaActivity.ActivityId})
		if err != nil {
			reply += "\n\nError fetching data from MongoDB: " + err.Error()
			return
		}
		if data.ActivityId == stravaActivity.ActivityId {
			reply += "\n\nHayoolooooo ngapain, Jangan Curang donggg! 😏 Kamu sudah pernah share aktivitas ini sebelumnya." +
				"\n*AOKWOKOKOWKOWKOWKKWOKOK* 🤣🤣" +
				"\nSana Lari lagi jangan malas! \n\n#tetap_semangat💪🏻"
			return
		}

		_, err = atdb.InsertOneDoc(db, col, stravaActivity)
		if err != nil {
			reply += "\n\nError saving data to MongoDB: " + err.Error()
		} else {
			reply += "Haiiiii kak, " + stravaActivity.Name + "! Berikut Progres Aktivitas kamu hari ini yaaa!! 😀" +
				"\n\n- Activity_id: " + stravaActivity.ActivityId +
				"\n- Name: " + stravaActivity.Name +
				"\n- Title: " + stravaActivity.Title +
				"\n- Date Time: " + stravaActivity.DateTime +
				"\n- Type Sport: " + stravaActivity.TypeSport +
				"\n- Distance: " + stravaActivity.Distance +
				"\n- Time Period: " + stravaActivity.TimePeriod +
				"\n- Elevation: " + stravaActivity.Elevation +
				"\n\nSemangat terus, jangan lupa jaga kesehatan dan tetap semangat!! 💪🏻💪🏻💪🏻" +
				"\n\n"
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
