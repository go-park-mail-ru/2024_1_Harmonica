package entity

import (
	"time"
)

type User struct {
	UserID     UserID    `db:"user_id" json:"user_id" swaggerignore:"true"`
	Email      string    `db:"email" json:"email"`
	Nickname   string    `db:"nickname" json:"nickname"`
	Password   string    `db:"password" json:"password"`
	AvatarURL  string    `db:"avatar_url" json:"avatar_url" swaggerignore:"true"`
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

// User list response model
// @Description User information
// @Description with user id, email, nickname and avatar_url
type UserList struct {
	Users []UserResponse `json:"users"`
}

type UserProfileResponse struct {
	User           UserResponse `json:"user"`
	FollowersCount uint64       `json:"followers_count"`
	IsOwner        bool         `json:"is_owner"`
}
