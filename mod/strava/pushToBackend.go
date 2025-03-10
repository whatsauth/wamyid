package strava

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func PushToBackend(phone, picture string) string {
	apiURL := "https://asia-southeast2-awangga.cloudfunctions.net/domyid/data/user"

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
