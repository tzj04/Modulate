package postgres

import (
	"database/sql"
	"time"

	"modulate/backend/internal/models"
)

type PostRepo struct {
	DB *sql.DB
}

func (r *PostRepo) Create(p *models.Post) error {
	return r.DB.QueryRow(`
		INSERT INTO posts (module_id, user_id, title, content)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`,
		p.ModuleID,
		p.UserID,
		p.Title,
		p.Content,
	).Scan(&p.ID, &p.CreatedAt)
}


func (r *PostRepo) GetByID(id int64) (*models.Post, error) {
	var p models.Post

	err := r.DB.QueryRow(`
		SELECT
			id,
			module_id,
			user_id,
			title,
			content,
			is_deleted,
			created_at,
			updated_at
		FROM posts
		WHERE id = $1
	`, id).Scan(
		&p.ID,
		&p.ModuleID,
		&p.UserID,
		&p.Title,
		&p.Content,
		&p.IsDeleted,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &p, nil
}


func (r *PostRepo) ListByModule(moduleID int64) ([]models.Post, error) {
	rows, err := r.DB.Query(`
		SELECT
			id,
			module_id,
			user_id,
			title,
			content,
			is_deleted,
			created_at,
			updated_at
		FROM posts
		WHERE module_id = $1
		ORDER BY created_at DESC
	`, moduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]models.Post, 0)

	for rows.Next() {
		var p models.Post
		if err := rows.Scan(
			&p.ID,
			&p.ModuleID,
			&p.UserID,
			&p.Title,
			&p.Content,
			&p.IsDeleted,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}


func (r *PostRepo) SoftDelete(id int64) error {
	_, err := r.DB.Exec(`
		UPDATE posts
		SET is_deleted = TRUE,
		    updated_at = $2
		WHERE id = $1
	`,
		id,
		time.Now(),
	)
	return err
}
