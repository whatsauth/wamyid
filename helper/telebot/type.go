package telebot

type Update struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID           int    `json:"id"`
			IsBot        bool   `json:"is_bot"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			Username     string `json:"username"`
			LanguageCode string `json:"language_code"`
		} `json:"from"`
		Chat struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
			Type      string `json:"type"`
		} `json:"chat"`
		Date     int    `json:"date"`
		Text     string `json:"text,omitempty"`
		Caption  string `json:"caption,omitempty"`
		Entities []struct {
			Offset int    `json:"offset"`
			Length int    `json:"length"`
			Type   string `json:"type"`
		} `json:"entities,omitempty"`
		Contact  *Contact `json:"contact,omitempty"`
		Location *struct {
			Longitude            float64 `json:"longitude"`
			Latitude             float64 `json:"latitude"`
			LivePeriod           int     `json:"live_period,omitempty"`
			Heading              int     `json:"heading,omitempty"`
			ProximityAlertRadius int     `json:"proximity_alert_radius,omitempty"`
		} `json:"location,omitempty"`
		Photo          []PhotoSize `json:"photo,omitempty"`
		ReplyToMessage *Message    `json:"reply_to_message,omitempty"` // Perubahan di sini

	} `json:"message"`
}
type Message struct {
	MessageID int `json:"message_id"`
	From      struct {
		ID           int    `json:"id"`
		IsBot        bool   `json:"is_bot"`
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		Username     string `json:"username"`
		LanguageCode string `json:"language_code"`
	} `json:"from"`
	Date     int           `json:"date"`
	Text     string        `json:"text,omitempty"`
	Caption  string        `json:"caption,omitempty"`
	Entities []interface{} `json:"entities,omitempty"`
	Contact  *Contact      `json:"contact,omitempty"`
	Location *struct {
		Longitude            float64 `json:"longitude"`
		Latitude             float64 `json:"latitude"`
		LivePeriod           int     `json:"live_period,omitempty"`
		Heading              int     `json:"heading,omitempty"`
		ProximityAlertRadius int     `json:"proximity_alert_radius,omitempty"`
	} `json:"location,omitempty"`
	Photo []PhotoSize `json:"photo,omitempty"`
}

type Contact struct {
	PhoneNumber string `json:"phone_number,omitempty"`
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	UserID      int    `json:"user_id,omitempty"`
}

// PhotoSize represents one size of a photo or a file/sticker thumbnail.
type PhotoSize struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	FileSize     int    `json:"file_size,omitempty"`
}

type SendMessagePayload struct {
	ChatID      int                 `json:"chat_id"`
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
