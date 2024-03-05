package db

import "time"

type Pin struct {
	PinId      int64     `db:"pin_id"`
	AuthorId   int64     `db:"author_id"`
	CreatedAt  time.Time `db:"created_at"`
	Caption    string    `db:"caption"`
	ClickUrl   string    `db:"click_url"`
	ContentUrl string    `db:"content_url"`
}

type User struct {
	UserID     int64     `db:"user_id" json:"user_id" swaggerignore:"true"`
	Email      string    `db:"email" json:"email"`
	Nickname   string    `db:"nickname" json:"nickname"`
	Password   string    `db:"password" json:"password"`
	RegisterAt time.Time `db:"register_at" json:"register_at" swaggerignore:"true"`
}
