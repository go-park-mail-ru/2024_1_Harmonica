package postgres

import "time"

type Pin struct {
	PinId      int64     `db:"pin_id"`
	AuthorId   int64     `db:"author_id"`
	CreatedAt  time.Time `db:"created_at"`
	Caption    string    `db:"caption"`
	ClickUrl   string    `db:"click_url"`
	ContentUrl string    `db:"content_url"`
}
