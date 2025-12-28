package postgres

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestGetByIDSuccess tests the successful retrieval of a user
func TestGetByIDSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := &UserRepo{DB: db}

	userID := int64(1)
	mock.ExpectQuery(`SELECT id, username, password, label, is_deleted, created_at FROM users`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "username", "password", "label", "is_deleted", "created_at",
		}).AddRow(userID, "testuser", "hashedpassword", "user", false, time.Now()))

	user, err := repo.GetByID(userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID != userID {
		t.Errorf("expected ID %d, got %d", userID, user.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

// TestGetByIDNotFound tests the case when the user does not exist
func TestGetByIDNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := &UserRepo{DB: db}

	userID := int64(42)
	mock.ExpectQuery(`SELECT id, username, password, label, is_deleted, created_at FROM users`).
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	user, err := repo.GetByID(userID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, sql.ErrNoRows) && err.Error() != "user not found" {
		t.Errorf("expected 'user not found' error, got %v", err)
	}
	if user != nil {
		t.Errorf("expected nil user, got %+v", user)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
