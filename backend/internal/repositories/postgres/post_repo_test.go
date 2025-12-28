package postgres

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"modulate/backend/internal/models"
)

// TestCreatePost tests the Create method
func TestCreatePost(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := &PostRepo{DB: db}

	post := &models.Post{
		TopicID:      1,
		UserID:       2,
		ParentPostID: nil,
		Content:      "Hello world",
	}

	mock.ExpectQuery(`INSERT INTO posts`).
		WithArgs(post.TopicID, post.UserID, post.ParentPostID, post.Content).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(1, time.Now()))

	err = repo.Create(post)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

// TestListThreadByTopic tests ListThreadByTopic
func TestListThreadByTopic(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := &PostRepo{DB: db}

	topicID := int64(1)

	rows := sqlmock.NewRows([]string{"id", "topic_id", "user_id", "parent_post_id", "content", "is_deleted", "created_at", "updated_at"}).
		AddRow(1, topicID, 2, nil, "Parent post", false, time.Now(), nil).
		AddRow(2, topicID, 3, 1, "Child post", false, time.Now(), nil)

	mock.ExpectQuery(`WITH RECURSIVE thread`).WithArgs(topicID).WillReturnRows(rows)

	posts, err := repo.ListThreadByTopic(topicID)
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

	mock.ExpectExec(`UPDATE posts SET is_deleted = TRUE`).
		WithArgs(postID, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.SoftDelete(postID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
