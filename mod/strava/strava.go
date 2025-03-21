package strava

import (
	"strings"

	"github.com/gocolly/colly"
)

var domApp = "strava.app.link"
var domWeb = "www.strava.com"
var isMaintenance = false

func maintenance(phone string) string {
	if phone != "6282268895372" {
		if isMaintenance {
			return "\n\nMaaf kak, sistem sedang maintenance. Coba lagi nanti ya."

		}
	}

	return ""
}

func extractStravaActivitySpan(e *colly.HTMLElement, activity *StravaActivity) {
	e.ForEach("span", func(_ int, el *colly.HTMLElement) {
		label := strings.ToLower(strings.TrimSpace(el.Text))
		value := strings.TrimSpace(el.DOM.Next().Text()) // Ambil elemen di sebelahnya

		switch label {
		case "distance":
			activity.Distance = value
		case "time":
			activity.MovingTime = value
		case "elevation":
			activity.Elevation = value
		}
	})

	// Jika elevation kosong, beri nilai default
	if activity.Elevation == "" {
		activity.Elevation = "0 m"
	}
}

func extractStravaProfileImg(e *colly.HTMLElement, name string) string {
	var imageURL string
	e.ForEach("img", func(_ int, imgEl *colly.HTMLElement) {
		imgTitle := imgEl.Attr("title")
		if imgTitle == name {
			imageURL = imgEl.Attr("src")
		}
	})

	return imageURL
}
