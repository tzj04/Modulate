package models

import "time"

type Post struct {
	ID           int64
	TopicID      int64
	UserID       int64
	ParentPostID *int64

	Content   string
	IsDeleted bool

	CreatedAt time.Time
	UpdatedAt *time.Time
}
