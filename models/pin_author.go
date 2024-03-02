package models

type PinAuthor struct {
	UserId   int64  `db:"user_id" json:"user_id"`
	Nickname string `db:"nickname" json:"nickname"`
}
