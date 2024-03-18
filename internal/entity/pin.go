package entity

import "time"

type Pin struct {
	PinId      int64     `db:"pin_id"`
	AuthorId   int64     `db:"author_id"`
	CreatedAt  time.Time `db:"created_at"`
	Caption    string    `db:"caption"`
	ClickUrl   string    `db:"click_url"`
	ContentUrl string    `db:"content_url"`
}

// PinResponse author model
// @Description User-author information
// @Description with user id and nickname
type PinAuthor struct {
	UserId   int64  `db:"user_id" json:"user_id"`
	Nickname string `db:"nickname" json:"nickname"`
}

// PinResponse model
// @Description PinResponse information
// @Description with author, pin id, created date, caption, click and content URLs.
type PinResponse struct {
	PinId      int64     `db:"pin_id" json:"pin_id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	Caption    string    `db:"caption" json:"caption"`
	ClickUrl   string    `db:"click_url" json:"click_url"`
	ContentUrl string    `db:"content_url" json:"content_url"`
	PinAuthor  `json:"author"`
}

// Pins model
// @Description Pins array of PinResponse
type Pins struct {
	Pins []PinResponse `json:"pins"`
}
