package telebot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const telegramAPI = "https://api.telegram.org/bot"

func SetWebhook(webhookURL, botToken string) error {
	url := fmt.Sprintf("%s%s/setWebhook", telegramAPI, botToken)
	payload := map[string]string{
		"url": webhookURL,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, io.NopCloser(bytes.NewBuffer(body)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to set webhook: %s", resp.Status)
	}

	return nil
}

func SendMessage(chatID int64, text, botToken string) error {
	url := fmt.Sprintf("%s%s/sendMessage", telegramAPI, botToken)
	payload := map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, io.NopCloser(bytes.NewBuffer(body)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message: %s", resp.Status)
	}

	return nil
}

func RequestPhoneNumber(chatID int64, botToken string) error {
	url := fmt.Sprintf("%s%s/sendMessage", telegramAPI, botToken)
	payload := SendMessagePayload{
		ChatID: chatID,
		Text:   "Please share your phone number.",
		ReplyMarkup: ReplyKeyboardMarkup{
			Keyboard: [][]KeyboardButton{
				{
					{Text: "Share Contact", RequestContact: true},
				},
			},
			OneTimeKeyboard: true,
			ResizeKeyboard:  true,
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	return nil
}
