package entity

type CommentRequest struct {
	Value string `json:"value"`
}

type CommentAuthor struct {
	UserId    UserID `db:"user_id" json:"user_id"`
	Nickname  string `db:"nickname" json:"nickname"`
	AvatarURL string `db:"avatar_url" json:"avatar_url"`
}

type CommentResponse struct {
	CommentId     CommentID `json:"comment_id"`
	Value         string    `json:"value" db:"text"`
	CommentAuthor `json:"user"`
}

type GetCommentsResponse struct {
	Comments []CommentResponse `json:"comments"`
}
