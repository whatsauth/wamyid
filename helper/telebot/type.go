package telebot

type Update struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID           int64  `json:"id"`
			IsBot        bool   `json:"is_bot"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			Username     string `json:"username"`
			LanguageCode string `json:"language_code"`
		} `json:"from"`
		Chat struct {
			ID        int64  `json:"id"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
			Type      string `json:"type"`
		} `json:"chat"`
		Date     int    `json:"date"`
		Text     string `json:"text,omitempty"`
		Entities []struct {
			Offset int    `json:"offset"`
			Length int    `json:"length"`
			Type   string `json:"type"`
		} `json:"entities,omitempty"`
		Contact *struct {
			PhoneNumber string `json:"phone_number,omitempty"`
			FirstName   string `json:"first_name,omitempty"`
			LastName    string `json:"last_name,omitempty"`
			UserID      int    `json:"user_id,omitempty"`
		} `json:"contact,omitempty"`
		Location *struct {
			Longitude float64 `json:"longitude"`
			Latitude  float64 `json:"latitude"`
		} `json:"location,omitempty"`
		LiveLocation *struct {
			Longitude float64 `json:"longitude"`
			Latitude  float64 `json:"latitude"`
		} `json:"live_location,omitempty"`
	} `json:"message"`
}

type SendMessagePayload struct {
	ChatID      int64               `json:"chat_id"`
	Text        string              `json:"text"`
	ReplyMarkup ReplyKeyboardMarkup `json:"reply_markup,omitempty"`
}

type ReplyKeyboardMarkup struct {
	Keyboard        [][]KeyboardButton `json:"keyboard"`
	ResizeKeyboard  bool               `json:"resize_keyboard,omitempty"`
	OneTimeKeyboard bool               `json:"one_time_keyboard,omitempty"`
}

type KeyboardButton struct {
	Text           string `json:"text"`
	RequestContact bool   `json:"request_contact,omitempty"`
}
