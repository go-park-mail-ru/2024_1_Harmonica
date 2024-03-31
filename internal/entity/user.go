package entity

import (
	"time"
)

type User struct {
	UserID     UserID    `db:"user_id" json:"user_id" swaggerignore:"true"`
	Email      string    `db:"email" json:"email"`
	Nickname   string    `db:"nickname" json:"nickname"`
	Password   string    `db:"password" json:"password"`
	AvatarURL  string    `db:"avatar_url" json:"avatar_url"`
	RegisterAt time.Time `db:"register_at" json:"register_at" swaggerignore:"true"`
}

// User response model
// @Description User information
// @Description with user id, email, nickname and avatar_url
type UserResponse struct {
	UserId    UserID `json:"user_id"`
	Email     string `json:"email"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
}

// ответ - для хэндлера => в сторону хэндлера перенести

// User list response model
// @Description User information
// @Description with user id, email, nickname and avatar_url
type UserList struct {
	Users []UserResponse `json:"users"`
}

type UserProfileResponse struct {
	User            UserResponse `json:"user"`
	FollowersNumber uint64       `json:"followers_number"`
	IsOwner         bool         `json:"is_owner"`
}
