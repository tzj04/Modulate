package repositories

import "modulate/backend/internal/models"

type CommentRepository interface {
	Create(comment *models.Comment) error
	ListThreadByPost(topicID int64) ([]models.Comment, error)
	SoftDelete(id int64) error
}
