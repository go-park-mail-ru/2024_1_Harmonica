package handler

type UserResponse struct {
	UserId   int64  `json:"user_id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
}
