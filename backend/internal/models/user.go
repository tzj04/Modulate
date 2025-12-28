package models

import "time"

type User struct {
	ID       int64
	Username string
	Password string
	Label    *string // Label is an optional, self-declared user role (e.g. "Student", "Professor").
	IsDeleted bool
	CreatedAt time.Time
}
