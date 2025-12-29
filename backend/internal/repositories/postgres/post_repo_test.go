package postgres

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"modulate/backend/internal/models"
)

// TestCreateTopic tests the Create method
func TestCreatePost(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := &PostRepo{DB: db}

	post := &models.Post{
		ModuleID: 1,
		UserID:   2,
		Title:    "New Post",
		Content:  "Post content",
	}

	mock.ExpectQuery(`INSERT INTO posts`).
		WithArgs(
			post.ModuleID,
			post.UserID,
			post.Title,
			post.Content,
		).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "created_at"}).
				AddRow(1, time.Now()),
		)

	if err := repo.Create(post); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if post.ID != 1 {
		t.Errorf("expected post ID 1, got %d", post.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

// TestGetByID tests GetByID
func TestGetPostByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := &PostRepo{DB: db}

	postID := int64(1)

	mock.ExpectQuery(`SELECT\s+id,\s+module_id,\s+user_id,\s+title,\s+content,\s+is_deleted,\s+created_at,\s+updated_at\s+FROM posts\s+WHERE id = \$1`).
		WithArgs(postID).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id",
				"module_id",
				"user_id",
				"title",
				"content",
				"is_deleted",
				"created_at",
				"updated_at",
			}).AddRow(postID, 1, 2, "New Post", "Post content", false, time.Now(), nil),
		)

	post, err := repo.GetByID(postID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if post.ID != postID {
		t.Errorf("expected ID %d, got %d", postID, post.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}


// TestListByModule tests ListByModule
func TestListPostsByModule(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := &PostRepo{DB: db}
	moduleID := int64(1)

	rows := sqlmock.NewRows([]string{
		"id",
		"module_id",
		"user_id",
		"title",
		"content",
		"is_deleted",
		"created_at",
		"updated_at",
	}).
		AddRow(1, moduleID, 2, "Post 1", "Content 1", false, time.Now(), nil).
		AddRow(2, moduleID, 3, "Post 2", "Content 2", false, time.Now(), nil)

	mock.ExpectQuery(`SELECT\s+id,\s+module_id,\s+user_id,\s+title,\s+content,\s+is_deleted,\s+created_at,\s+updated_at\s+FROM posts\s+WHERE module_id = \$1`).
		WithArgs(moduleID).
		WillReturnRows(rows)

	posts, err := repo.ListByModule(moduleID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(posts) != 2 {
		t.Errorf("expected 2 posts, got %d", len(posts))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

// TestSoftDelete tests SoftDelete
func TestSoftDeletePost(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := &PostRepo{DB: db}

	postID := int64(1)

	mock.ExpectExec(`UPDATE posts`).
		WithArgs(postID, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := repo.SoftDelete(postID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}