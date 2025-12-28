package postgres

import (
	"database/sql"

	"modulate/backend/internal/models"
)

type ModuleRepo struct {
	DB *sql.DB
}

func (r *ModuleRepo) Create(m *models.Module) error {
	return r.DB.QueryRow(`
		INSERT INTO modules (code, title, description)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`, m.Code, m.Title, m.Description).Scan(&m.ID, &m.CreatedAt)
}

func (r *ModuleRepo) ListAll() ([]models.Module, error) {
	rows, err := r.DB.Query(`
		SELECT id, code, title, description, created_at
		FROM modules
		ORDER BY code
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var modules []models.Module
	for rows.Next() {
		var m models.Module
		if err := rows.Scan(
			&m.ID,
			&m.Code,
			&m.Title,
			&m.Description,
			&m.CreatedAt,
		); err != nil {
			return nil, err
		}
		modules = append(modules, m)
	}

	return modules, nil
}
