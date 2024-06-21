package controller

import (
	"encoding/json"
	"log"
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

	log.Printf("Message from %s: %s", update.Message.From.Username, update.Message.Text)

	chatID := update.Message.Chat.ID
	text := "Hello, " + update.Message.From.FirstName

	prof, err := helper.GetAppProfile(waphonenumber, config.Mongoconn)
	if err != nil {
		resp.Response = err.Error()
		helper.WriteResponse(w, http.StatusServiceUnavailable, resp)
		return
	}
	if err := telebot.SendMessage(chatID, text, prof.TelegramToken); err != nil {
		resp.Response = err.Error()
		helper.WriteResponse(w, http.StatusServiceUnavailable, resp)
		return
	}
	helper.WriteResponse(w, http.StatusOK, resp)
}
