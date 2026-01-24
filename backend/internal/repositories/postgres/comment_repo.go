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
    // Initialize as empty slice so JSON is [] instead of null
    comments := make([]models.Comment, 0)

    // Query with JOIN
    rows, err := r.DB.Query(`
        SELECT 
            c.id, c.post_id, c.user_id, c.parent_comment_id, 
            c.content, c.is_deleted, c.created_at, c.updated_at,
            u.username
        FROM comments c
        JOIN users u ON c.user_id = u.id
        WHERE c.post_id = $1
        ORDER BY c.created_at ASC
    `, postID)
    
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var c models.Comment
        if err := rows.Scan(
            &c.ID, &c.PostID, &c.UserID, &c.ParentCommentID, 
            &c.Content, &c.IsDeleted, &c.CreatedAt, &c.UpdatedAt,
            &c.Username,
        ); err != nil {
            return nil, err
        }
        comments = append(comments, c)
    }

    return comments, nil
}


func (r *CommentRepo) Update(id int64, userID int64, content string) error {
    _, err := r.DB.Exec(`
        UPDATE comments 
        SET content = $1, updated_at = $2 
        WHERE id = $3 AND user_id = $4 AND is_deleted = FALSE`,
        content, time.Now(), id, userID,
    )
    return err
}

func (r *CommentRepo) SoftDelete(id int64, userID int64) error {
    _, err := r.DB.Exec(`
        UPDATE comments
        SET is_deleted = TRUE,
            content = '[Comment deleted by user]',
            updated_at = $1
        WHERE id = $2 AND user_id = $3`,
        time.Now(), id, userID,
    )
    return err
}