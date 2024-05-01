package entity

type SearchRequest struct {
	SearchUsers  bool   `json:"users"`
	SearchPins   bool   `json:"pins"`
	SearchBoards bool   `json:"boards"`
	SearchQuery  string `json:"search_query"`
}

type SearchResult struct {
	ResultType string `json:"type"`
	Object     any    `json:"object"`
}

type SearchUser struct {
	UserId    UserID `db:"user_id" json:"user_id"`
	Nickname  string `db:"nickname" json:"nickname"`
	AvatarURL string `db:"avatar_url" json:"avatar_url"`
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
