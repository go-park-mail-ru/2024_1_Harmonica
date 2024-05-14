package entity

import (
	"html"
	"time"
)

type MessageStatus string

const (
	MessageStatusUnread MessageStatus = "unread"
	MessageStatusRead   MessageStatus = "read"
)

var MessageStatuses = []MessageStatus{MessageStatusUnread, MessageStatusRead}

type Message struct {
	MessageId  int64         `db:"message_id" json:"message_id"`
	SenderId   UserID        `db:"sender_id" json:"sender_id"`
	ReceiverId UserID        `db:"receiver_id" json:"receiver_id"`
	Text       string        `db:"text" json:"text"`
	Status     MessageStatus `db:"status" json:"status"`
	SentAt     time.Time     `db:"sent_at" json:"sent_at"`
}

//type MessageRequest struct {
//	ReceiverId UserID `json:"receiver_id"`
//	Text       string `json:"text"`
//}

func (m *Message) Sanitize() {
	m.Text = html.EscapeString(m.Text)
}

type MessageResponse struct {
	SenderId UserID `db:"sender_id" json:"sender_id"`
	Text     string `db:"text" json:"text"`
	//Status   MessageStatus `db:"status" json:"status"`
	IsRead bool      `db:"message_read" json:"message_read"`
	SentAt time.Time `db:"sent_at" json:"sent_at"`
}

type Action string

const (
	ActionMessage Action = "CHAT_MESSAGE"
)

var Actions = []Action{ActionMessage}

type ChatMessage struct {
	Action  Action `json:"action"`
	Payload struct {
		Text       string `json:"text"`
		SenderId   UserID `json:"sender_id"`
		ReceiverId UserID `json:"receiver_id"`
	} `json:"payload"`
}

type UserFromChat struct {
	UserID    UserID `db:"user_id" json:"user_id" swaggerignore:"true"`
	Nickname  string `db:"nickname" json:"nickname"`
	AvatarURL string `db:"avatar_url" json:"avatar_url" swaggerignore:"true"`
}

type Messages struct {
	User     UserFromChat      `db:"user" json:"user"`
	Messages []MessageResponse `db:"messages" json:"messages"`
}

type UserChat struct {
	User        UserFromChat    `db:"user" json:"user"`
	LastMessage MessageResponse `db:"chat_last_message" json:"chat_last_message"`
}

type UserChats struct {
	Chats []UserChat `json:"chats"`
}
