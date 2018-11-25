package models

import strfmt "github.com/go-openapi/strfmt"

var (
	SendMessage          = 1
	CreateChat           = 2
	ConnectToChat        = 3
	ConnectToChatSuccess = 4
	SendMessageSuccess   = 5
)

type Chat struct {
	ID      int           `json:"id"`
	Name    string        `json:"name"`
	Members []string      `json:"members"`
	History []ChatMessage `json:"history"`
}

type ChatMessage struct {
	Chat        int              `json:"chat"`
	User        *User            `json:"user"`
	Event       int              `json:"event"`
	Text        string           `json:"text"`
	NewUsername string           `json:"newusername,omitempty"`
	Data        *[]Chat          `json:"data,omitempty"`
	Created     *strfmt.DateTime `json:"created,omitempty"`
}
