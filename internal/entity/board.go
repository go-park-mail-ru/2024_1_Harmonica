package entity

import "time"

type VisibilityType string

const (
	VisibilityPublic  VisibilityType = "public"
	VisibilityPrivate VisibilityType = "private"
)

var VisibilityTypes = []VisibilityType{VisibilityPublic, VisibilityPrivate}

type Board struct {
	BoardID BoardID `db:"board_id" json:"board_id"`
	//Title          string    `db:"title" json:"title"`
	Title          int64          `db:"title" json:"title"` // исправить на строку!!!
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

type FullBoard struct {
	Board        `json:"board"`
	BoardAuthors []BoardAuthor     `json:"author"`
	Pins         []FeedPinResponse `json:"pins"`
}
