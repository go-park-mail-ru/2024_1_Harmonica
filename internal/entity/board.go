package entity

import "time"

type VisibilityType string

const (
	VisibilityPublic  VisibilityType = "public"
	VisibilityPrivate VisibilityType = "private"
)

var VisibilityTypes = []VisibilityType{VisibilityPublic, VisibilityPrivate}

type Board struct {
	BoardID        BoardID        `db:"board_id" json:"board_id"`
	Title          string         `db:"title" json:"title"`
	CreatedAt      time.Time      `db:"created_at" json:"created_at"`
	Description    string         `db:"description" json:"description"`
	CoverURL       string         `db:"cover_url" json:"cover_url"`
	VisibilityType VisibilityType `db:"visibility_type" json:"visibility_type"`
}

type UserBoards struct {
	Boards []Board `json:"boards"`
}

type BoardAuthor struct {
	UserId    UserID `db:"user_id" json:"user_id"`
	Nickname  string `db:"nickname" json:"nickname"`
	AvatarURL string `db:"avatar_url" json:"avatar_url"`
}

type BoardPinResponse struct {
	PinId      PinID  `db:"pin_id" json:"pin_id"`
	ContentUrl string `db:"content_url" json:"content_url"`
	PinAuthor  `json:"author"`
}

type FullBoard struct {
	Board        `json:"board"`
	BoardAuthors []BoardAuthor      `json:"authors"`
	Pins         []BoardPinResponse `json:"pins"`
}
