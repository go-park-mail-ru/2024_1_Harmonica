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
	UserId     int64     `db:"user_id"`
	Email      string    `db:"email"`
	Nickname   string    `db:"nickname"`
	Password   string    `db:"password"`
	RegisterAt time.Time `db:"register_at"`
}
