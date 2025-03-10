package strava

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

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

func getTokenFromCookie(req *colly.Request, cookieName string) string {
	cookie := req.Headers.Get("cookies")
	cookies := strings.Split(cookie, "; ")
	for _, c := range cookies {
		parts := strings.SplitN(c, "=", 2)
		if len(parts) == 2 && parts[0] == cookieName {
			return parts[1]
		}
	}
	return ""
}
