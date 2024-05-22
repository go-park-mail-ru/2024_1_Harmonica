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
	BoardID        BoardID        `db:"board_id" json:"board_id"`
	Title          string         `db:"title" json:"title"`
	CreatedAt      time.Time      `db:"created_at" json:"created_at"`
	Description    string         `db:"description" json:"description"`
	CoverURL       string         `db:"cover_url" json:"cover_url" swaggerignore:"true"`
	CoverDX        int64          `json:"cover_width"`
	CoverDY        int64          `json:"cover_height"`
	VisibilityType VisibilityType `db:"visibility_type" json:"visibility_type"`
	IsOwner        bool           `json:"is_owner"`
}

func (b *Board) Sanitize() {
	b.Title = html.EscapeString(b.Title)
	b.Description = html.EscapeString(b.Description)
}

type UserBoard struct {
	BoardID              BoardID        `db:"board_id" json:"board_id"`
	Title                string         `db:"title" json:"title"`
	CreatedAt            time.Time      `db:"created_at" json:"created_at"`
	Description          string         `db:"description" json:"description"`
	CoverURL             string         `db:"cover_url" json:"cover_url" swaggerignore:"true"`
	VisibilityType       VisibilityType `db:"visibility_type" json:"visibility_type"`
	RecentPinContentUrls []string       `db:"recent_pins" json:"recent_pins"`
}

func (b *UserBoard) Sanitize() {
	b.Title = html.EscapeString(b.Title)
	b.Description = html.EscapeString(b.Description)
}

type UserBoards struct {
	Boards []UserBoard `json:"boards"`
}

func (b *UserBoards) Sanitize() {
	for _, board := range b.Boards {
		board.Sanitize()
	}
}

type UserBoardWithoutPin struct {
	BoardID        BoardID        `db:"board_id" json:"board_id"`
	Title          string         `db:"title" json:"title"`
	Description    string         `db:"description" json:"description"`
	CoverURL       string         `db:"cover_url" json:"cover_url" swaggerignore:"true"`
	VisibilityType VisibilityType `db:"visibility_type" json:"visibility_type"`
}

func (b *UserBoardWithoutPin) Sanitize() {
	b.Title = html.EscapeString(b.Title)
	b.Description = html.EscapeString(b.Description)
}

type UserBoardsWithoutPin struct {
	Boards []UserBoardWithoutPin `json:"boards"`
}

func (b *UserBoardsWithoutPin) Sanitize() {
	for _, board := range b.Boards {
		board.Sanitize()
	}
}

type BoardAuthor struct {
	UserId    UserID `db:"user_id" json:"user_id"`
	Nickname  string `db:"nickname" json:"nickname"`
	AvatarURL string `db:"avatar_url" json:"avatar_url"`
	AvatarDX  int64  `json:"avatar_width"`
	AvatarDY  int64  `json:"avatar_height"`
}

func (b *BoardAuthor) Sanitize() {
	b.Nickname = html.EscapeString(b.Nickname)
}

type BoardPinResponse struct {
	PinId      PinID  `db:"pin_id" json:"pin_id"`
	ContentUrl string `db:"content_url" json:"content_url"`
	ContentDX  int64  `json:"content_width"`
	ContentDY  int64  `json:"content_height"`
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
