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
	loc, _ := time.LoadLocation("Asia/Jakarta")
	t, err := time.ParseInLocation(time.RFC3339, dateTime, loc)
	if err != nil {
		return ""
	}

	return t.Format("02 Jan 2006 15:04")
}
