package strava

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gocroot/helper/atdb"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/mongo"
)

var activityId string
var rawUrl = "https://strava.app.link/NlALiL9oxRb"

func StravaHandler(Pesan itmodel.IteungMessage, db *mongo.Database) string {
	c := colly.NewCollector(
		colly.AllowedDomains("strava.app.link"),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong: ", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Page visited: ", r.Request.URL)
	})

	c.OnHTML("a", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		if strings.Contains(link, "/activities/") {
			parts := strings.Split(link, "/activities/")
			if len(parts) > 1 {
				activityId = strings.Split(parts[1], "/")[0]
				fullActivityURL := "https://www.strava.com/activities/" + activityId

				scrapeStravaActivity(db, fullActivityURL)
			}
		}
	})

	err := c.Visit(rawUrl)
	if err != nil {
		fmt.Println("Error visiting URL1", err)
	}

	return "Strava activity has been scraped"
}

func scrapeStravaActivity(db *mongo.Database, url string) {
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

		activities = append(activities, stravaActivity)
	})

	c.OnScraped(func(r *colly.Response) {
		atdb.InsertOneDoc(db, "stravas", activities)
	})

	err := c.Visit(url)
	if err != nil {
		fmt.Println("Error visiting URL", err)
	}
}
