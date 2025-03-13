package strava

var domApp = "strava.app.link"
var domWeb = "www.strava.com"
var isMaintenance = true

func maintenance(phone string) string {
	if phone != "6282268895372" {
		if isMaintenance {
			return "\n\nMaaf kak, sistem sedang maintenance. Coba lagi nanti ya."

		}
	}

	return ""
}
