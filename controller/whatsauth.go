package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gocroot/config"
	"github.com/gocroot/helper"
	"github.com/gocroot/model"
)

func HandleRequest(respw http.ResponseWriter, req *http.Request) {
	var resp model.Response
	var msg model.IteungMessage
	url := req.URL.Path
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
	} else {
		dt := &model.WebHook{
			URL:    config.WebhookURL,
			Secret: config.WebhookSecret,
		}
		res, err := helper.RefreshToken(dt, config.WAPhoneNumber, config.WAAPIGetToken, config.Mongoconn)
		if err != nil {
			resp.Response = err.Error()
		} else {
			resp.Response = helper.Jsonstr(res.ModifiedCount)
			resp.Response = "{\"urlpath\":" + url + ",\"updatedecount\":" + resp.Response + "}"
		}

	}
	helper.WriteResponse(respw, http.StatusOK, resp)
}
