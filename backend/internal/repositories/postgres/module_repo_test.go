package postgres

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"modulate/backend/internal/models"
)

// TestCreateModule tests the Create method of ModuleRepo
func TestCreateModule(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := &ModuleRepo{DB: db}

	module := &models.Module{Code: "CS1010", Title: "Programming", Description: "Intro"}

	// Setup expected query and returned row
	mock.ExpectQuery(`INSERT INTO modules`).
		WithArgs(module.Code, module.Title, module.Description).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(1, time.Now()))

	// Call method
	err = repo.Create(module)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

// TestListAllModules tests the ListAll method of ModuleRepo
func TestListAllModules(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := &ModuleRepo{DB: db}

	// Setup expected rows to be returned
	rows := sqlmock.NewRows([]string{"id", "code", "title", "description", "created_at"}).
		AddRow(1, "CS1010", "Programming", "Intro", time.Now()).
		AddRow(2, "CS1020", "Data Structures", "DS Intro", time.Now())

	mock.ExpectQuery(`SELECT id, code, title, description, created_at FROM modules`).
		WillReturnRows(rows)

	// Call method
	modules, err := repo.ListAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(modules) != 2 {
		t.Errorf("expected 2 modules, got %d", len(modules))
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}
