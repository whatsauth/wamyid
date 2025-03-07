package strava

import (
	"regexp"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gocroot/helper/atdb"
	"github.com/whatsauth/itmodel"
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

	return "link strava activity: " + rawUrl + "\n#mental_health"
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
		_, err := atdb.InsertOneDoc(db, "strava", stravaActivity)
		if err != nil {
			reply += "\n\nError saving data to MongoDB: " + err.Error()
		} else {
			reply += responseStravaData(stravaActivity)
		}
	})

	err := c.Visit(url)
	if err != nil {
		return "\nError visiting URL2" + err.Error()
	}

	return reply
}

func responseStravaData(stravaActivity StravaActivity) string {
	return "\n\nHaiiiiii kak, " + stravaActivity.Name + "! Berikut Progres Aktivitas kamu hari ini yaaa!! 😀" +
		"\n\nActivity_id: " + stravaActivity.ActivityId +
		"\nName: " + stravaActivity.Name +
		"\nTitle: " + stravaActivity.Title +
		"\nDate Time: " + stravaActivity.DateTime +
		"\nType Sport: " + stravaActivity.TypeSport +
		"\nDistance: " + stravaActivity.Distance +
		"\nTime Period: " + stravaActivity.TimePeriod +
		"\nElevation: " + stravaActivity.Elevation +
		"\n\nSemangat terus, jangan lupa jaga kesehatan dan tetap semangat!! 💪🏻💪🏻💪🏻" +
		"\n\n"
}

func extractStravaLink(text string) string {
	re := regexp.MustCompile(`https?://strava\.app\.link/\S+`)
	match := re.FindString(text)

	return match
}
