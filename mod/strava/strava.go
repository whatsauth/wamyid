package strava

var domApp = "strava.app.link"
var domWeb = "www.strava.com"
var isMaintenance = false

const AdminPicture = "https://lh3.googleusercontent.com/a/ACg8ocK27sU9YXcfmLm9Zw_MtUW0kT--NA5XmjMGaJUvEdKl65cx6QQ=s96-c"

func maintenance(phone string) string {
	if phone != "6282268895372" {
		if isMaintenance {
			return "\n\nMaaf kak, sistem sedang maintenance. Coba lagi nanti ya."

		}
	}

	return ""
}
