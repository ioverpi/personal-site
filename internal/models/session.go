package models

import "time"

type Session struct {
	ID        int
	UserID    int
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
