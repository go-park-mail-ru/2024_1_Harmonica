package models

import "time"

type PinAuthor struct {
	UserId   int64  `db:"user_id" json:"user_id"`
	Nickname string `db:"nickname" json:"nickname"`
}

type Pin struct {
	PinId      int64     `db:"pin_id" json:"pin_id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	Caption    string    `db:"caption" json:"caption"`
	ClickUrl   string    `db:"click_url" json:"click_url"`
	ContentUrl string    `db:"content_url" json:"content_url"`
	PinAuthor  `json:"author"`
}

type Pins struct {
	Pins []Pin `json:"pins"`
}
