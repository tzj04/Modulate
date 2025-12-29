package repositories

import "modulate/backend/internal/models"

type PostRepository interface {
	Create(post *models.Post) error
	GetByID(id int64) (*models.Post, error)
	ListByModule(moduleID int64) ([]models.Post, error)
	SoftDelete(id int64) error
}
