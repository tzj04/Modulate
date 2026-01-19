package models

import "time"

type Module struct {
    ID          int64     `json:"id"`
    Code        string    `json:"code"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"created_at"`
}