package repositories

import "modulate/backend/internal/models"

type TopicRepository interface {
	Create(topic *models.Topic) error
	GetByID(id int64) (*models.Topic, error)
	ListByModule(moduleID int64) ([]models.Topic, error)
	SoftDelete(id int64) error
}
