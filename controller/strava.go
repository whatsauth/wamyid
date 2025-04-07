package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gocroot/config"
	"github.com/gocroot/helper"
	"github.com/gocroot/mod/pomokit"
	"github.com/gocroot/mod/strava"
)

func GetStravaActivities(respw http.ResponseWriter, req *http.Request) {
	doc, err := helper.GetAllDoc[[]strava.StravaActivity](config.Mongoconn, "strava_activity")
	if err != nil {
		helper.WriteResponse(respw, http.StatusInternalServerError, err.Error())
		return
	}
	helper.WriteResponse(respw, http.StatusOK, doc)
}

func GetStravaActivitiesWithGrupIDFromPomokit(respw http.ResponseWriter, req *http.Request) {
	// Ambil data dari koleksi strava_activity
	stravaDocs, err := helper.GetAllDoc[[]strava.StravaActivity](config.Mongoconn, "strava_activity")
	if err != nil {
		helper.WriteResponse(respw, http.StatusInternalServerError, err.Error())
		return
	}

	// Ambil data dari koleksi pomokit
	pomokitDocs, err := helper.GetAllDoc[[]pomokit.PomodoroReport](config.Mongoconn, "pomokit")
	if err != nil {
		helper.WriteResponse(respw, http.StatusInternalServerError, err.Error())
		return
	}

	// Buat mapping phonenumber -> grup_id
	phoneToGroup := make(map[string]string)
	for _, p := range pomokitDocs {
		if phoneToGroup[p.PhoneNumber] == "" && p.WaGroupID != "" {
			phoneToGroup[p.PhoneNumber] = p.WaGroupID // Ambil grup ID pertama yang tidak kosong
		}
	}

	// Gabungkan data tanpa mengubah struct
	var response []map[string]interface{}
	for _, s := range stravaDocs {
		mergedData := StructToMap(s)                          // Convert struct ke map
		mergedData["wagroupid"] = phoneToGroup[s.PhoneNumber] // Tambahkan grup_id
		response = append(response, mergedData)
	}

	// Kirim respons
	helper.WriteResponse(respw, http.StatusOK, response)
}

func StructToMap(data interface{}) map[string]interface{} {
	jsonData, _ := json.Marshal(data)
	var result map[string]interface{}
	json.Unmarshal(jsonData, &result)
	return result
}
