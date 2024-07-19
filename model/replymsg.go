package model

type Chats struct {
	IdChats   string  `json:"id_chats" bson:"idChats"`
	Message   string  `json:"message" bson:"message"`
	Responses string  `json:"responses" bson:"responses"`
	Score     float64 `json:"score" bson:"score"`
}

type Requests struct {
	Messages string `json:"messages" bson:"messages"`
}
