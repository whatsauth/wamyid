package controller

import (
	"net/http"

	"github.com/gocroot/config"
	"github.com/gocroot/helper"
	"github.com/gocroot/mod/pomokit"
)

func GetPomokitData(respw http.ResponseWriter, req *http.Request) {
	doc, err := helper.GetAllDoc[[]pomokit.PomodoroReport](config.Mongoconn, "pomokit")
	if err != nil {
		helper.WriteResponse(respw, http.StatusInternalServerError, err.Error())
		return
	}
	helper.WriteResponse(respw, http.StatusOK, doc)
}