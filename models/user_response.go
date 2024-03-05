package models

// User response model
// @Description User information
// @Description with user id, email and nickname
type UserResponse struct {
	UserId   int64  `json:"user_id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
}
