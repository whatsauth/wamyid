package strava

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

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

func formatDateTimeToIndo(dateTime string) string {
	layout := "2006-01-02T15:04:05"
	t, err := time.ParseInLocation(layout, dateTime, time.Local)
	if err != nil {
		return "Error parsing date time: " + err.Error()
	}

	return t.Format("02 Jan 2006 15:04 WIB")
}
