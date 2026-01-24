package models

import "time"

type Comment struct {
    ID              int64      `json:"id"`
    PostID          int64      `json:"post_id"`
    UserID          int64      `json:"user_id"`
    ParentCommentID *int64     `json:"parent_comment_id"`

    Username        string     `json:"username"`
    Content         string     `json:"content"`

    IsDeleted       bool       `json:"is_deleted"`
    CreatedAt       time.Time  `json:"created_at"`
    UpdatedAt       *time.Time `json:"updated_at"`
	
    Children        []Comment  `json:"children,omitempty"` // omitempty hides field if no replies
}