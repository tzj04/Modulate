package postgres

import (
    "database/sql"

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

func (r *PostRepo) Update(id int64, userID int64, title string, content string) error {
    result, err := r.DB.Exec(`
        UPDATE posts 
        SET title = $1, content = $2, updated_at = CURRENT_TIMESTAMP
        WHERE id = $3 AND user_id = $4`, // Ensure only the owner can edit
        title, content, id, userID,
    )
    if err != nil {
        return err
    }

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        return sql.ErrNoRows // Means either post doesn't exist or user doesn't own it
    }
    return nil
}

func (r *PostRepo) GetByID(id int64) (*models.Post, error) {
    var p models.Post
    err := r.DB.QueryRow(`
        SELECT p.id, p.module_id, p.user_id, p.title, p.content, 
               p.is_deleted, p.created_at, p.updated_at, u.username
        FROM posts p
        JOIN users u ON p.user_id = u.id
        WHERE p.id = $1
    `, id).Scan(
        &p.ID, &p.ModuleID, &p.UserID, &p.Title, &p.Content,
        &p.IsDeleted, &p.CreatedAt, &p.UpdatedAt, &p.Username,
    )
    if err != nil {
        return nil, err
    }
    return &p, nil
}

func (r *PostRepo) ListByModule(moduleID int64) ([]models.Post, error) {
    rows, err := r.DB.Query(`
        SELECT p.id, p.module_id, p.user_id, p.title, p.content, 
               p.is_deleted, p.created_at, p.updated_at, u.username
        FROM posts p
        JOIN users u ON p.user_id = u.id
        WHERE p.module_id = $1 AND p.is_deleted = FALSE
        ORDER BY p.created_at DESC
    `, moduleID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    posts := make([]models.Post, 0)
    for rows.Next() {
        var p models.Post
        if err := rows.Scan(
            &p.ID, &p.ModuleID, &p.UserID, &p.Title, &p.Content,
            &p.IsDeleted, &p.CreatedAt, &p.UpdatedAt, &p.Username,
        ); err != nil {
            return nil, err
        }
        posts = append(posts, p)
    }
    return posts, nil
}

func (r *PostRepo) SoftDelete(id int64) error {
    query := `
        UPDATE posts 
        SET title = '[Deleted]', 
            content = '[This post has been removed by the author]', 
            is_deleted = true 
        WHERE id = $1`
    result, err := r.DB.Exec(query, id)
    if err != nil {
        return err
    }

    rows, _ := result.RowsAffected()
    if rows == 0 {
        return sql.ErrNoRows
    }
    return nil
}