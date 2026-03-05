package models

type ChatsInformations struct {
	ChatID    int64  `bson:"chat_id,omitempty" json:"chat_id,omitempty"`
	Username  string `bson:"username" json:"username" default:"nil"`
	ChatTitle string `bson:"chat_title" json:"chat_title" default:"nil"`
	Language  string `bson:"language" json:"language" default:"en_us"`
}
