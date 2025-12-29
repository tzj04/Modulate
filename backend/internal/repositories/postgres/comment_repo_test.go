package postgres

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"modulate/backend/internal/models"
)

// TestCreatePost tests the Create method
func TestCreateComment(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := &CommentRepo{DB: db}

	comment := &models.Comment{
		PostID:          1,
		UserID:          2,
		ParentCommentID: nil,
		Content:         "Hello world",
	}

	mock.ExpectQuery(`INSERT INTO comments`).
		WithArgs(
			comment.PostID,
			comment.UserID,
			comment.ParentCommentID,
			comment.Content,
		).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "created_at"}).
				AddRow(1, time.Now()),
		)

	if err := repo.Create(comment); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}


// TestListThreadByTopic tests ListThreadByTopic
func TestListThreadByPost(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := &CommentRepo{DB: db}

	postID := int64(1)

	rows := sqlmock.NewRows([]string{
		"id",
		"post_id",
		"user_id",
		"parent_comment_id",
		"content",
		"is_deleted",
		"created_at",
		"updated_at",
	}).
		AddRow(1, postID, 2, nil, "Parent comment", false, time.Now(), nil).
		AddRow(2, postID, 3, 1, "Child comment", false, time.Now(), nil)

	mock.ExpectQuery(`WITH RECURSIVE thread`).
		WithArgs(postID).
		WillReturnRows(rows)

	comments, err := repo.ListThreadByPost(postID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(comments) != 2 {
		t.Errorf("expected 2 comments, got %d", len(comments))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}


// TestSoftDelete tests SoftDelete
func TestSoftDeleteComment(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := &CommentRepo{DB: db}

	commentID := int64(1)

	mock.ExpectExec(`UPDATE comments`).
		WithArgs(commentID, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := repo.SoftDelete(commentID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
