package postgres

import "time"

type User struct {
	UserId     int64     `db:"user_id"`
	Email      string    `db:"email"`
	Nickname   string    `db:"nickname"`
	Password   string    `db:"password"`
	RegisterAt time.Time `db:"register_at"`
}
