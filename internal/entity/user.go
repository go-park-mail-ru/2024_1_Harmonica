package entity

import (
	"html"
	"time"
)

type User struct {
	UserID     UserID    `db:"user_id" json:"user_id" swaggerignore:"true"`
	Email      string    `db:"email" json:"email"`
	Nickname   string    `db:"nickname" json:"nickname"`
	Password   string    `db:"password" json:"password"`
	AvatarURL  string    `db:"avatar_url" json:"avatar_url" swaggerignore:"true"`
	AvatarDX   int64     `json:"avatar_width"`
	AvatarDY   int64     `json:"avatar_height"`
	RegisterAt time.Time `db:"register_at" json:"register_at" swaggerignore:"true"`
}

func (u *User) Sanitize() {
	u.Email = html.EscapeString(u.Email)
	u.Nickname = html.EscapeString(u.Nickname)
}

// User response model
// @Description User information
// @Description with user id, email, nickname and avatar_url
type UserResponse struct {
	UserId    UserID `json:"user_id"`
	Email     string `json:"email"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
	AvatarDX  int64  `json:"avatar_width"`
	AvatarDY  int64  `json:"avatar_height"`
}

func (u *UserResponse) Sanitize() {
	u.Email = html.EscapeString(u.Email)
	u.Nickname = html.EscapeString(u.Nickname)
}

// User list response model
// @Description User information
// @Description with user id, email, nickname and avatar_url
type UserList struct {
	Users []UserResponse `json:"users"`
}

func (u *UserList) Sanitize() {
	for _, user := range u.Users {
		user.Sanitize()
	}
}

type UserProfileResponse struct {
	User               UserResponse `db:"user" json:"user"`
	SubscriptionsCount uint64       `db:"subscriptions_count" json:"subscriptions_count"`
	SubscribersCount   uint64       `db:"subscribers_count" json:"subscribers_count"`
	IsSubscribed       bool         `db:"is_subscribed" json:"is_subscribed"` // важно для авторизованного пользователя
	IsOwner            bool         `db:"is_owner" json:"is_owner"`
}

func (u *UserProfileResponse) Sanitize() {
	u.User.Sanitize()
}

type UserSubscriptionInfo struct {
	UserId           UserID `json:"user_id"`
	Nickname         string `json:"nickname"`
	AvatarURL        string `json:"avatar_url"`
	SubscribersCount uint64 `db:"subscribers_count" json:"subscribers_count"`
}

type UserSubscribers struct {
	Subscribers []UserSubscriptionInfo `json:"subscribers"`
}

type UserSubscriptions struct {
	Subscriptions []UserSubscriptionInfo `json:"subscriptions"`
}

type SearchUsers struct {
	Users []SearchUser `json:"users"`
}
