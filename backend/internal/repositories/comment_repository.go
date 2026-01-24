package repositories

import "modulate/backend/internal/models"

type CommentRepository interface {
	Create(comment *models.Comment) error
	ListThreadByPost(topicID int64) ([]models.Comment, error)
	Update(id int64, userID int64, newContent string) error
	SoftDelete(id int64, userID int64) error
}
