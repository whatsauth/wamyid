package controller

import (
	"net/http"

	"github.com/gocroot/config"
	"github.com/gocroot/helper"
	"github.com/gocroot/mod/strava"
)

func GetStravaActivities(respw http.ResponseWriter, req *http.Request) {
	doc, err := helper.GetAllDoc[strava.StravaActivity](config.Mongoconn, "strava_activity")
	if err != nil {
		helper.WriteResponse(respw, http.StatusInternalServerError, err.Error())
		return
	}
	helper.WriteResponse(respw, http.StatusOK, doc)
}
