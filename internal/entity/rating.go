package entity

type Rating struct {
	RatingID    uint64 `db:"rating-id" json:"rating-id"`
	Title       string `db:"title" json:"title"`
	RatingCount int    `db:"rating_count" json:"rating_count"`
	User        string `db:"user" json:"user"`
}

type RatingList struct {
	Ratings []Rating
}
