package strava

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gocolly/colly"
)

func pushToBackend(phone, picture, token string) string {
	data := map[string]string{
		"phonenumber":          phone,
		"stravaprofilepicture": picture,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "Error marshalling data: " + err.Error()
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "Error creating request: " + err.Error()
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("login", token)
	// req.Header.Set("Cookie", "login="+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "Error sending request: " + err.Error()
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "Error response status: " + resp.Status
	}

	return "Data berhasil disimpan"
}

// func extractTokenFromCookie(cookieHeader, cookieName string) string {
// 	cookies := strings.Split(cookieHeader, "; ")
// 	for _, cookie := range cookies {
// 		parts := strings.SplitN(cookie, "=", 2)
// 		if len(parts) == 2 && parts[0] == cookieName {
// 			return parts[1]
// 		}
// 	}
// 	return ""
// }

// func getCookie(r *http.Request, cname string) string {
// 	cookie, err := r.Cookie(cname)
// 	if err != nil {
// 		if err == http.ErrNoCookie {
// 			return "Cookie not found"
// 		}
// 		return "Error getting cookie: " + err.Error()
// 	}
// 	return cookie.Value
// }

func getCookieFromColly(c *colly.Collector, cookieName string) string {
	// Ambil semua cookies dari domain target
	cookies := c.Cookies("https://www.do.my.id") // Ganti dengan domain target

	for _, cookie := range cookies {
		if cookie.Name == cookieName {
			return cookie.Value
		}
	}
	return "Token not found"
}
