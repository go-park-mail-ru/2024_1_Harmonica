package server

import (
	"harmonica/internal/entity"
	"time"
)

type Session struct {
	UserId entity.UserID
	Expiry time.Time
}

func (s Session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}

func CleanupSessions() {
	ticker := time.NewTicker(sessionsCleanupTime)
	for {
		<-ticker.C
		Sessions.Range(func(key, value interface{}) bool {
			if session, ok := value.(Session); ok {
				if time.Now().After(session.Expiry) {
					Sessions.Delete(key)
				}
			}
			return true
		})
	}
}
