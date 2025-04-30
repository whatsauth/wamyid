package controller

import (
	"net/http"

	"github.com/gocroot/config"
	"github.com/gocroot/helper"
	"github.com/gocroot/mod/pomokit"
	"go.mongodb.org/mongo-driver/bson"
)

func GetPomokitData(respw http.ResponseWriter, req *http.Request) {
	doc, err := helper.GetAllDoc[[]pomokit.PomodoroReport](config.Mongoconn, "pomokit")
	if err != nil {
		helper.WriteResponse(respw, http.StatusInternalServerError, err.Error())
		return
	}
	helper.WriteResponse(respw, http.StatusOK, doc)
}

func GetPomokitDataByPhonenumber(respw http.ResponseWriter, req *http.Request) {
	phonenumber := helper.GetParam(req)
	filter := bson.M{"phonenumber": phonenumber}
	doc, err := helper.GetAllDocs[[]pomokit.PomodoroReport](config.Mongoconn, "pomokit", filter)
	if err != nil {
		helper.WriteResponse(respw, http.StatusInternalServerError, err.Error())
		return
	}
	helper.WriteResponse(respw, http.StatusOK, doc)
}
