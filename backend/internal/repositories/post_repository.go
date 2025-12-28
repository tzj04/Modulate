package repositories

import "modulate/backend/internal/models"

type PostRepository interface {
	Create(post *models.Post) error
	ListThreadByTopic(topicID int64) ([]models.Post, error)
	SoftDelete(id int64) error
}
