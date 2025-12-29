package postgres

import (
	"database/sql"
	"time"

	"modulate/backend/internal/models"
)

type CommentRepo struct {
	DB *sql.DB
}

func (r *CommentRepo) Create(c *models.Comment) error {
	return r.DB.QueryRow(`
		INSERT INTO comments (post_id, user_id, parent_comment_id, content)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`,
		c.PostID,
		c.UserID,
		c.ParentCommentID,
		c.Content,
	).Scan(&c.ID, &c.CreatedAt)
}


func (r *CommentRepo) ListThreadByPost(postID int64) ([]models.Comment, error) {
	rows, err := r.DB.Query(`
		WITH RECURSIVE thread AS (
			SELECT
				id,
				post_id,
				user_id,
				parent_comment_id,
				content,
				is_deleted,
				created_at,
				updated_at
			FROM comments
			WHERE post_id = $1
			  AND parent_comment_id IS NULL

			UNION ALL

			SELECT
				c.id,
				c.post_id,
				c.user_id,
				c.parent_comment_id,
				c.content,
				c.is_deleted,
				c.created_at,
				c.updated_at
			FROM comments c
			JOIN thread t ON c.parent_comment_id = t.id
		)
		SELECT *
		FROM thread
		ORDER BY created_at ASC
	`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(
			&c.ID,
			&c.PostID,
			&c.UserID,
			&c.ParentCommentID,
			&c.Content,
			&c.IsDeleted,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}


func (r *CommentRepo) SoftDelete(id int64) error {
	_, err := r.DB.Exec(`
		UPDATE comments
		SET is_deleted = TRUE,
		    updated_at = $2
		WHERE id = $1`,
		id,
		time.Now(),
	)
	return err
}