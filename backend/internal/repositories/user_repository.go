package repositories

import "modulate/backend/internal/models"

type UserRepository interface {
	GetByID(id int64) (*models.User, error)
}
