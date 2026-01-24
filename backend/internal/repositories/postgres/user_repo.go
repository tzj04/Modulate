package postgres

import (
	"database/sql"
	"modulate/backend/internal/models"
)

type UserRepo struct {
	DB *sql.DB
}

func (r *UserRepo) Create(user *models.User) error {
	query := `
		INSERT INTO users (username, password, label, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, created_at`
	
	return r.DB.QueryRow(query, user.Username, user.Password, user.Label).
		Scan(&user.ID, &user.CreatedAt)
}

func (r *UserRepo) GetByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, password, label, is_deleted, created_at 
              FROM users WHERE username = $1 AND is_deleted = false`
	
	err := r.DB.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Password, 
		&user.Label, &user.IsDeleted, &user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepo) GetByID(id int64) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, label, is_deleted, created_at 
              FROM users WHERE id = $1`
	
	err := r.DB.QueryRow(query, id).Scan(
		&user.ID, 
		&user.Username, 
		&user.Label, 
		&user.IsDeleted, 
		&user.CreatedAt,
	)
	return user, err
}