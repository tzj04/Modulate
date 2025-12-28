package postgres

import (
	"database/sql"
	"errors"

	"modulate/backend/internal/models"
)

type UserRepo struct {
	DB *sql.DB
}

func (r *UserRepo) GetByID(id int64) (*models.User, error) {
	var u models.User
	err := r.DB.QueryRow(`
		SELECT id, username, password, label, is_deleted, created_at
		FROM users
		WHERE id = $1
	`, id).Scan(
		&u.ID,
		&u.Username,
		&u.Password,
		&u.Label,
		&u.IsDeleted,
		&u.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}
