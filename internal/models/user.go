package models

import "time"

type User struct {
	ID        int
	Email     string
	Name      string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}
