package model

import "time"

type Chats struct {
	IdChats   string  `json:"id_chats" bson:"idChats"`
	Message   string  `json:"message" bson:"message"`
	Responses string  `json:"responses" bson:"responses"`
	Score     float64 `json:"score" bson:"score"`
}

type Requests struct {
	Messages string `json:"messages" bson:"messages"`
}

type LogInbox struct {
	ID        string    `bson:"_id,omitempty"`
	From      string    `bson:"from,omitempty"`
	Message   string    `bson:"messsage,omitempty"`
	CreatedAt time.Time `bson:"createdAt"`
}
