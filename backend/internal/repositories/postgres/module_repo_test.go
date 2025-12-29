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

	module := &models.Module{
		Code:        "CS1101S",
		Title:       "Programming Methodology",
		Description: "Intro to programming",
	}

	mock.ExpectQuery(`INSERT INTO modules`).
		WithArgs(module.Code, module.Title, module.Description).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "created_at"}).
				AddRow(1, time.Now()),
		)

	if err := repo.Create(module); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if module.ID != 1 {
		t.Errorf("expected module ID 1, got %d", module.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
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

	rows := sqlmock.NewRows([]string{
		"id",
		"code",
		"title",
		"description",
		"created_at",
	}).
		AddRow(1, "CS1101S", "Programming Methodology", "Intro to programming", time.Now()).
		AddRow(2, "CS1231S", "Discrete Structures", "Discrete Math for Computing", time.Now())

	mock.ExpectQuery(`SELECT id, code, title, description, created_at FROM modules`).
		WillReturnRows(rows)

	modules, err := repo.ListAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(modules) != 2 {
		t.Errorf("expected 2 modules, got %d", len(modules))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
