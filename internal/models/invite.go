package models

import "time"

type Invite struct {
	ID        int
	Email     string
	Token     string
	InvitedBy int
	UsedAt    *time.Time
	ExpiresAt time.Time
	CreatedAt time.Time
}

func (i *Invite) IsExpired() bool {
	return time.Now().After(i.ExpiresAt)
}

func (i *Invite) IsUsed() bool {
	return i.UsedAt != nil
}
