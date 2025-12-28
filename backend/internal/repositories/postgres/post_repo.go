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
		INSERT INTO posts (topic_id, user_id, parent_post_id, content)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`,
		p.TopicID, p.UserID, p.ParentPostID, p.Content,
	).Scan(&p.ID, &p.CreatedAt)
}

func (r *PostRepo) ListThreadByTopic(topicID int64) ([]models.Post, error) {
	rows, err := r.DB.Query(`
		WITH RECURSIVE thread AS (
			SELECT id, topic_id, user_id, parent_post_id, content, is_deleted, created_at, updated_at
			FROM posts
			WHERE topic_id = $1 AND parent_post_id IS NULL
			UNION ALL
			SELECT p.id, p.topic_id, p.user_id, p.parent_post_id, p.content, p.is_deleted, p.created_at, p.updated_at
			FROM posts p
			JOIN thread t ON p.parent_post_id = t.id
		)
		SELECT * FROM thread ORDER BY created_at`, topicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(
			&p.ID, &p.TopicID, &p.UserID,
			&p.ParentPostID, &p.Content,
			&p.IsDeleted, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *PostRepo) SoftDelete(id int64) error {
	_, err := r.DB.Exec(`
		UPDATE posts SET is_deleted = TRUE, updated_at = $2 WHERE id = $1`,
		id, time.Now(),
	)
	return err
}
