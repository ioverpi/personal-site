package models

import "time"

type Quote struct {
	ID        int
	Content   string
	Author    string
	IsOwn     bool
	CreatedAt time.Time
}
