package models

import "time"

type Post struct {
	ID        int64
	Content   string // Content stores the body of a post (length constrained at DB level)
	TopicID   int64
	UserID    int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
