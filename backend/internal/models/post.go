package models

import "time"

type Post struct {
	ID        int64
	ModuleID  int64
	UserID    int64

	Title     string
	Content   string

	IsDeleted bool

	CreatedAt time.Time
	UpdatedAt *time.Time
}
