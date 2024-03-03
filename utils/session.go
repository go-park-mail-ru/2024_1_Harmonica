package utils

import "time"

type Session struct {
	UserId int64
	Expiry time.Time
}

func (s Session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}
