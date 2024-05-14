package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gocroot/config"
	"github.com/gocroot/helper"
	"github.com/gocroot/model"
)

func GetHome(respw http.ResponseWriter, req *http.Request) {
	var resp model.Response
	resp.Response = helper.GetIPaddress()
	helper.WriteResponse(respw, http.StatusOK, resp)
}

func PostInbox(respw http.ResponseWriter, req *http.Request) {
	var resp model.Response
	var msg model.IteungMessage
	httpstatus := http.StatusUnauthorized
	resp.Response = "Wrong Secret"
	if helper.GetSecretFromHeader(req) == config.WebhookSecret {
		err := json.NewDecoder(req.Body).Decode(&msg)
		if err != nil {
			resp.Response = err.Error()
		} else {
			resp, err = helper.WebHook(config.WAKeyword, config.WAPhoneNumber, config.WAAPIQRLogin, config.WAAPIMessage, msg, config.Mongoconn)
			if err != nil {
				resp.Response = err.Error()
			}
		}
	}
	helper.WriteResponse(respw, httpstatus, resp)
}

func GetNewToken(respw http.ResponseWriter, req *http.Request) {
	var resp model.Response
	httpstatus := http.StatusServiceUnavailable
	dt := &model.WebHook{
		URL:    config.WebhookURL,
		Secret: config.WebhookSecret,
	}
	res, err := helper.RefreshToken(dt, config.WAPhoneNumber, config.WAAPIGetToken, config.Mongoconn)
	if err != nil {
		resp.Response = err.Error()
	} else {
		resp.Response = helper.Jsonstr(res.ModifiedCount)
		httpstatus = http.StatusOK
	}
	helper.WriteResponse(respw, httpstatus, resp)
}

func NotFound(respw http.ResponseWriter, req *http.Request) {
	var resp model.Response
	resp.Response = "Not Found"
	helper.WriteResponse(respw, http.StatusNotFound, resp)
}
