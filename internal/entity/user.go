package entity

import "time"

type User struct {
	UserID     int64     `db:"user_id" json:"user_id" swaggerignore:"true"`
	Email      string    `db:"email" json:"email"`
	Nickname   string    `db:"nickname" json:"nickname"`
	Password   string    `db:"password" json:"password"`
	RegisterAt time.Time `db:"register_at" json:"register_at" swaggerignore:"true"`
}

// User response model
// @Description User information
// @Description with user id, email and nickname
type UserResponse struct {
	UserId   int64  `json:"user_id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
}
