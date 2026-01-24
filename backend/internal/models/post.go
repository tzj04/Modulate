package models

import "time"

type Post struct {
    ID        int64      `json:"id"`
    ModuleID  int64      `json:"module_id"`

    UserID    int64      `json:"user_id"`
    Username  string     `json:"username"`

    Title     string     `json:"title"`
    Content   string     `json:"content"`
    
    IsDeleted bool       `json:"is_deleted"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt *time.Time `json:"updated_at"`
}