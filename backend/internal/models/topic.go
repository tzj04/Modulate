package models

import "time"

type Topic struct {
	ID        int64
	Title     string
	ModuleID  int64
	UserID    int64
	CreatedAt time.Time
}
