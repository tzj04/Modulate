package repositories

import "modulate/backend/internal/models"

type PostRepository interface {
	Create(post *models.Post) error
	Update(id int64, userID int64, title string, content string) error
	GetByID(id int64) (*models.Post, error)
	ListByModule(moduleID int64) ([]models.Post, error)
	SoftDelete(id int64) error
}
