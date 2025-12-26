package models

import "time"

type Module struct {
	ID          int64
	Code        string
	Title       string
	Description string
	CreatedAt time.Time
}
