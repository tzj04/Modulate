package models

import "time"

type User struct {
    ID        int64     `json:"id"`
    Username  string    `json:"username"`
    Password  string    `json:"-"`
    Label     *string   `json:"label"`
    IsDeleted bool      `json:"is_deleted"`
    CreatedAt time.Time `json:"created_at"`
}