package repositories

import "modulate/backend/internal/models"

type ModuleRepository interface {
	Create(module *models.Module) error
	ListAll() ([]models.Module, error)
}
