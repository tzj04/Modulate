package models

import "time"

type Comment struct {
	ID           int64
	PostID      int64
	UserID       int64
	ParentCommentID *int64

	Content   string
	IsDeleted bool

	CreatedAt time.Time
	UpdatedAt *time.Time

	Children []Comment `json:"children,omitempty"`
}
