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
	SenderId   UserID        `db:"sender_id" json:"sender_id"`
	ReceiverId UserID        `db:"receiver_id" json:"receiver_id"`
	Text       string        `db:"text" json:"text"`
	Status     MessageStatus `db:"status" json:"status"`
	SentAt     time.Time     `db:"sent_at" json:"sent_at"`
}

type Messages struct {
	Messages []MessageResponse `json:"messages"`
}
