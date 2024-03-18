package handler

import "time"

type Session struct {
	UserId int64
	Expiry time.Time
}

func (s Session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}

func CleanupSessions() {
	ticker := time.NewTicker(sessionsCleanupTime)
	for {
		<-ticker.C
		sessions.Range(func(key, value interface{}) bool {
			if session, ok := value.(Session); ok {
				if time.Now().After(session.Expiry) {
					sessions.Delete(key)
				}
			}
			return true
		})
	}
}
