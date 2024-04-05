package entity

import (
	"time"
)

// Pin model from DB
// @Description Full pin information
type Pin struct {
	PinId         PinID     `db:"pin_id" json:"pin_id" swaggerignore:"true"`
	AuthorId      UserID    `db:"author_id" json:"author_id" swaggerignore:"true"`
	CreatedAt     time.Time `db:"created_at" json:"created_at" swaggerignore:"true"`
	Title         string    `db:"title" json:"title"`
	Description   string    `db:"description" json:"description"`
	AllowComments bool      `db:"allow_comments" json:"allow_comments"`
	ClickUrl      string    `db:"click_url" json:"click_url"`
	ContentUrl    string    `db:"content_url" json:"content_url" swaggerignore:"true"`
}

// Pin response author model
// @Description User-author information
// @Description with user id, nickname and avatar
type PinAuthor struct {
	UserId    UserID `db:"user_id" json:"user_id"`
	Nickname  string `db:"nickname" json:"nickname"`
	AvatarURL string `db:"avatar_url" json:"avatar_url"`
}

// Pin page response model
// @Description Full pin information
type PinPageResponse struct {
	PinId         PinID     `db:"pin_id" json:"pin_id"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	Title         string    `db:"title" json:"title"`
	Description   string    `db:"description" json:"description"`
	AllowComments bool      `db:"allow_comments" json:"allow_comments"`
	ClickUrl      string    `db:"click_url" json:"click_url"`
	ContentUrl    string    `db:"content_url" json:"content_url"`
	LikesCount    uint64    `db:"likes_count" json:"likes_count"`
	IsOwner       bool      `json:"is_owner"`
	PinAuthor     `json:"author"`
}

// Feed pin response model
// @Description PinResponse information
// @Description with author, pin id and content URL.
type FeedPinResponse struct {
	PinId      PinID  `db:"pin_id" json:"pin_id"`
	ContentUrl string `db:"content_url" json:"content_url"`
	PinAuthor  `json:"author"`
}

// Pins model
// @Description Pins array of FeedPinResponse
type FeedPins struct {
	Pins []FeedPinResponse `json:"pins"`
}

type UserPinResponse struct {
	PinId      PinID  `db:"pin_id" json:"pin_id"`
	ContentUrl string `db:"content_url" json:"content_url"`
}

type UserPins struct {
	Pins []UserPinResponse `json:"pins"`
}
