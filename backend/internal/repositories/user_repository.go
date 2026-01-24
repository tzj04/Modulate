package repositories

import "modulate/backend/internal/models"

type UserRepository interface {
    Create(user *models.User) error
    GetByUsername(username string) (*models.User, error)
    GetByID(id int64) (*models.User, error)
}