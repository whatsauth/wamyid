package telebot

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/whatsauth/itmodel"
)

// parseUpdateToIteungMessage parses the incoming update to IteungMessage
func ParseUpdateToIteungMessage(update Update, botToken string) itmodel.IteungMessage {
	var iteungMessage itmodel.IteungMessage

	// Map the necessary fields
	if update.Message.Contact != nil {
		iteungMessage.Reply_phone_number = update.Message.Contact.PhoneNumber
	}
	iteungMessage.Phone_number = strconv.FormatInt(update.Message.From.ID, 10)
	iteungMessage.Chat_number = strconv.FormatInt(update.Message.Chat.ID, 10)
	iteungMessage.Alias_name = update.Message.From.FirstName + " " + update.Message.From.LastName
	iteungMessage.Message = update.Message.Text
	iteungMessage.Is_group = update.Message.Chat.Type == "group"
	iteungMessage.Group_id = strconv.FormatInt(update.Message.Chat.ID, 10)
	iteungMessage.Group_name = update.Message.Chat.Username

	if update.Message.Location != nil {
		iteungMessage.Latitude = update.Message.Location.Latitude
		iteungMessage.Longitude = update.Message.Location.Longitude
	}
	if update.Message.LiveLocation != nil {
		iteungMessage.LiveLoc = true
		iteungMessage.Latitude = update.Message.LiveLocation.Latitude
		iteungMessage.Longitude = update.Message.LiveLocation.Longitude
	}
	if update.Message.Photo != nil && len(update.Message.Photo) > 0 {
		// Handle the received photo, for example, store the file_id of the largest size
		largestPhoto := update.Message.Photo[len(update.Message.Photo)-1]
		// Download the photo
		fileURL, err := GetFileURL(largestPhoto.FileID, botToken)
		if err == nil {
			fileBytes, err := DownloadFile(fileURL)
			if err == nil {
				// Convert the photo to base64
				iteungMessage.Filedata = base64.StdEncoding.EncodeToString(fileBytes)
				iteungMessage.Filename = largestPhoto.FileID
				if update.Message.Caption != "" {
					iteungMessage.Message = update.Message.Caption
				}
			}
		}
	}

	// Check if the update is a reply to a message
	if update.Message.ReplyToMessage != nil {
		// Check if the user who replied is the same as the user who sent the original message
		if update.Message.From.ID == update.Message.ReplyToMessage.From.ID {
			// Check if the replied message has location
			if update.Message.ReplyToMessage.Location != nil {
				iteungMessage.LiveLoc = true
				iteungMessage.Latitude = update.Message.ReplyToMessage.Location.Latitude
				iteungMessage.Longitude = update.Message.ReplyToMessage.Location.Longitude
				iteungMessage.Message = update.Message.Text
			}
		}
	}

	return iteungMessage
}

type FileResponse struct {
	OK     bool `json:"ok"`
	Result struct {
		FileID   string `json:"file_id"`
		FileSize int    `json:"file_size"`
		FilePath string `json:"file_path"`
	} `json:"result"`
}

// GetFileURL returns the full URL of a file
func GetFileURL(fileID, botToken string) (string, error) {
	url := fmt.Sprintf("%s%s/getFile?file_id=%s", telegramAPI, botToken, fileID)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get file URL: %s", resp.Status)
	}

	var fileResp FileResponse
	if err := json.NewDecoder(resp.Body).Decode(&fileResp); err != nil {
		return "", err
	}

	if !fileResp.OK {
		return "", fmt.Errorf("failed to get file URL: response not OK")
	}

	fileURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", botToken, fileResp.Result.FilePath)
	return fileURL, nil
}

// DownloadFile downloads the file from the provided URL
func DownloadFile(fileURL string) ([]byte, error) {
	resp, err := http.Get(fileURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download file: %s", resp.Status)
	}

	fileBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return fileBytes, nil
}
