package entity

type SearchRequest struct {
	SearchQuery string `json:"search_query"`
}

type SearchResult struct {
	Users  []SearchUser  `json:"users"`
	Pins   []SearchPin   `json:"pins"`
	Boards []SearchBoard `json:"boards"`
}

type SearchUser struct {
	UserId    UserID `db:"user_id" json:"user_id"`
	Nickname  string `db:"nickname" json:"nickname"`
	AvatarURL string `db:"avatar_url" json:"avatar_url"`
	SubsCount int    `db:"subs" json:"subscribers_count"`
}

type SearchPin struct {
	PinId      PinID  `db:"pin_id" json:"pin_id"`
	Title      string `db:"title" json:"title"`
	ContentURL string `db:"content_url" json:"content_url"`
}

type SearchBoard struct {
	BoardId  BoardID `db:"board_id" json:"board_id"`
	Title    string  `db:"title" json:"title"`
	CoverURL string  `db:"cover_url" json:"cover_url"`
}
