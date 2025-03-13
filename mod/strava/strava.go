package strava

var domApp = "strava.app.link"
var domWeb = "www.strava.com"
var isMaintenance = true

func maintenance(phone string) string {
	reply := ""
	if phone != "6282268895372" {
		if isMaintenance {
			reply += "\n\nMaaf kak, sistem sedang maintenance. Coba lagi nanti ya."
			return reply
		}
	}

	return reply
}
