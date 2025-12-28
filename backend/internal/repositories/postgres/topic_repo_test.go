package postgres

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"modulate/backend/internal/models"
)

// TestCreateTopic tests the Create method
func TestCreateTopic(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := &TopicRepo{DB: db}

	topic := &models.Topic{
		ModuleID: 1,
		UserID:   2,
		Title:    "New Topic",
		Content:  "Topic content",
	}

	mock.ExpectQuery(`INSERT INTO topics`).
		WithArgs(topic.ModuleID, topic.UserID, topic.Title, topic.Content).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(1, time.Now()))

	err = repo.Create(topic)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

// TestGetByID tests GetByID
func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := &TopicRepo{DB: db}

	topicID := int64(1)
	mock.ExpectQuery(`SELECT id, module_id, user_id, title, content, is_deleted, created_at, updated_at FROM topics`).
		WithArgs(topicID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "module_id", "user_id", "title", "content", "is_deleted", "created_at", "updated_at",
		}).AddRow(topicID, 1, 2, "New Topic", "Topic content", false, time.Now(), nil))

	topic, err := repo.GetByID(topicID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if topic.ID != topicID {
		t.Errorf("expected ID %d, got %d", topicID, topic.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

// TestListByModule tests ListByModule
func TestListByModule(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := &TopicRepo{DB: db}

	moduleID := int64(1)
	rows := sqlmock.NewRows([]string{
		"id", "module_id", "user_id", "title", "content", "is_deleted", "created_at", "updated_at",
	}).
		AddRow(1, moduleID, 2, "Topic 1", "Content 1", false, time.Now(), nil).
		AddRow(2, moduleID, 3, "Topic 2", "Content 2", false, time.Now(), nil)

	mock.ExpectQuery(`SELECT id, module_id, user_id, title, content, is_deleted, created_at, updated_at FROM topics`).
		WithArgs(moduleID).
		WillReturnRows(rows)

	topics, err := repo.ListByModule(moduleID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(topics) != 2 {
		t.Errorf("expected 2 topics, got %d", len(topics))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

// TestSoftDelete tests SoftDelete
func TestSoftDeleteTopic(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := &TopicRepo{DB: db}

	topicID := int64(1)

	mock.ExpectExec(`UPDATE topics SET is_deleted = TRUE`).
		WithArgs(topicID, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.SoftDelete(topicID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
