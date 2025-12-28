package postgres

import (
	"database/sql"
	"time"

	"modulate/backend/internal/models"
)

type TopicRepo struct {
	DB *sql.DB
}

func (r *TopicRepo) Create(t *models.Topic) error {
	return r.DB.QueryRow(`
		INSERT INTO topics (module_id, user_id, title, content)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`,
		t.ModuleID, t.UserID, t.Title, t.Content,
	).Scan(&t.ID, &t.CreatedAt)
}

func (r *TopicRepo) GetByID(id int64) (*models.Topic, error) {
	var t models.Topic
	err := r.DB.QueryRow(`
		SELECT id, module_id, user_id, title, content, is_deleted, created_at, updated_at
		FROM topics WHERE id = $1`, id,
	).Scan(
		&t.ID, &t.ModuleID, &t.UserID,
		&t.Title, &t.Content, &t.IsDeleted,
		&t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TopicRepo) ListByModule(moduleID int64) ([]models.Topic, error) {
	rows, err := r.DB.Query(`
		SELECT id, module_id, user_id, title, content, is_deleted, created_at, updated_at
		FROM topics
		WHERE module_id = $1
		ORDER BY created_at DESC`, moduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []models.Topic
	for rows.Next() {
		var t models.Topic
		if err := rows.Scan(
			&t.ID, &t.ModuleID, &t.UserID,
			&t.Title, &t.Content, &t.IsDeleted,
			&t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		topics = append(topics, t)
	}
	return topics, nil
}

func (r *TopicRepo) SoftDelete(id int64) error {
	_, err := r.DB.Exec(`
		UPDATE topics SET is_deleted = TRUE, updated_at = $2 WHERE id = $1`,
		id, time.Now(),
	)
	return err
}
