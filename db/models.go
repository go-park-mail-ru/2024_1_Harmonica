package db

import (
	"time"
)

type User struct {
	UserId     int64     `db:"user_id" json:"user_id"`
	Email      string    `db:"email" json:"email,omitempty"`
	Nickname   string    `db:"nickname" json:"nickname,omitempty"`
	Password   string    `db:"password" json:"-"`
	RegisterAt time.Time `db:"register_at" json:"-"`
}

type Pin struct {
	PinId      int64     `db:"pin_id" json:"pin_id"`
	AuthorId   int64     `db:"author_id" json:"-"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	Caption    string    `db:"caption" json:"caption,omitempty"`
	ClickUrl   string    `db:"click_url" json:"click_url,omitempty"`
	ContentUrl string    `db:"content_url" json:"content_url,omitempty"`
	User       `json:"author"`
}
