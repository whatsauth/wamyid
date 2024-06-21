package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gocroot/config"
	"github.com/gocroot/helper"
	"github.com/gocroot/helper/telebot"
	"github.com/whatsauth/itmodel"
)

func TelebotWebhook(w http.ResponseWriter, r *http.Request) {
	var resp itmodel.Response
	waphonenumber := helper.GetParam(r)

	var update telebot.Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		resp.Response = err.Error()
		helper.WriteResponse(w, http.StatusBadRequest, resp)
		return
	}

	chatID := update.Message.Chat.ID
	prof, err := helper.GetAppProfile(waphonenumber, config.Mongoconn)
	if err != nil {
		resp.Response = err.Error()
		helper.WriteResponse(w, http.StatusServiceUnavailable, resp)
		return
	}

	if update.Message.Contact != nil && update.Message.Contact.PhoneNumber != "" {
		text := "Hello, " + update.Message.From.FirstName + " nomor handphone " + update.Message.Contact.PhoneNumber
		if err := telebot.SendMessage(chatID, text, prof.TelegramToken); err != nil {
			resp.Response = err.Error()
			helper.WriteResponse(w, http.StatusInternalServerError, resp)
			return
		}
	} else {
		text := "Hello, " + update.Message.From.FirstName + " nomor handphone tidak tersedia"
		if err := telebot.SendMessage(chatID, text, prof.TelegramToken); err != nil {
			resp.Response = err.Error()
			helper.WriteResponse(w, http.StatusInternalServerError, resp)
			return
		}
		err := telebot.RequestPhoneNumber(chatID, prof.TelegramToken)
		if err != nil {
			resp.Response = err.Error()
			helper.WriteResponse(w, http.StatusExpectationFailed, resp)
			return
		}
	}

	helper.WriteResponse(w, http.StatusOK, resp)
}
