package telebot

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
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

func SendTextMessage(chatID int64, text, botToken string) error {
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

// SendImageMessage sends an image to a specified chatID
func SendImageMessage(chatID int64, photo, caption, botToken string) error {
	url := fmt.Sprintf("%s%s/sendPhoto", telegramAPI, botToken)

	// Create a buffer to write our photo data
	var photoBuf bytes.Buffer
	writer := multipart.NewWriter(&photoBuf)

	// Add the chat_id field
	if err := writer.WriteField("chat_id", fmt.Sprintf("%d", chatID)); err != nil {
		return err
	}

	// Add the caption field if provided
	if caption != "" {
		if err := writer.WriteField("caption", caption); err != nil {
			return err
		}
	}

	// Create a form file field for the photo
	fileWriter, err := writer.CreateFormFile("photo", "image.jpg")
	if err != nil {
		return err
	}

	// Convert base64 string to bytes and write to the form file
	photoBytes, err := base64.StdEncoding.DecodeString(photo)
	if err != nil {
		return err
	}

	if _, err := fileWriter.Write(photoBytes); err != nil {
		return err
	}

	// Close the writer to finalize the multipart form data
	if err := writer.Close(); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, &photoBuf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send image: %s", resp.Status)
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
