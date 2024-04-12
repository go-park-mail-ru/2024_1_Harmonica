package entity

import (
	"html"
	"time"
)

type VisibilityType string

const (
	VisibilityPublic  VisibilityType = "public"
	VisibilityPrivate VisibilityType = "private"
)

var VisibilityTypes = []VisibilityType{VisibilityPublic, VisibilityPrivate}

type Board struct {
	BoardID     BoardID        `db:"board_id" json:"board_id"`
	Title       string         `db:"title" json:"title"`
	CreatedAt   time.Time      `db:"created_at" json:"created_at"`
	Description string         `db:"description" json:"description"`
	CoverURL    string         `db:"cover_url" json:"cover_url" swaggerignore:"true"`
	Visibility  VisibilityType `db:"visibility" json:"visibility"`
	IsOwner     bool           `json:"is_owner"`
}

func (b *Board) Sanitize() {
	b.Title = html.EscapeString(b.Title)
	b.Description = html.EscapeString(b.Description)
}

type UserBoards struct {
	Boards []Board `json:"boards"`
}

func (b *UserBoards) Sanitize() {
	for _, board := range b.Boards {
		board.Sanitize()
	}
}

type BoardAuthor struct {
	UserId    UserID `db:"user_id" json:"user_id"`
	Nickname  string `db:"nickname" json:"nickname"`
	AvatarURL string `db:"avatar_url" json:"avatar_url"`
}

func (b *BoardAuthor) Sanitize() {
	b.Nickname = html.EscapeString(b.Nickname)
}

type BoardPinResponse struct {
	PinId      PinID  `db:"pin_id" json:"pin_id"`
	ContentUrl string `db:"content_url" json:"content_url"`
	PinAuthor  `json:"author"`
}

func (b *BoardPinResponse) Sanitize() {
	b.PinAuthor.Sanitize()
}

type FullBoard struct {
	Board        `json:"board"`
	BoardAuthors []BoardAuthor      `json:"authors"`
	Pins         []BoardPinResponse `json:"pins"`
}

func (b *FullBoard) Sanitize() {
	b.Board.Sanitize()
	for _, author := range b.BoardAuthors {
		author.Sanitize()
	}
	for _, pin := range b.Pins {
		pin.Sanitize()
	}
}
