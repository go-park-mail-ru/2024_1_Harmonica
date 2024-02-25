package postgres

type User struct {
	User_id  int64  `db:"user_id"`
	Email    string `db:"email"`
	Nickname string `db:"nickname"`
	Password string `db:"password"`
}
