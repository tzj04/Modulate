package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"modulate/backend/internal/handlers"
	"modulate/backend/internal/repositories/postgres"
	"modulate/backend/internal/routes"
)

func main() {
	// Connect to DB
	db, err := sql.Open("postgres", "postgres://user:password@localhost:5432/modulate?sslmode=disable")
	if err != nil {
		log.Fatal("Error opening database handle:", err)
	}
	defer db.Close()

	// Ping dbto actually verify the connection
	if err := db.Ping(); err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	// Set reasonable connection settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Initialize repositories
	moduleRepo := &postgres.ModuleRepo{DB: db}
	postRepo := &postgres.PostRepo{DB: db}
	commentRepo := &postgres.CommentRepo{DB: db}

	// Initialize handlers
	moduleHandler := handlers.NewModuleHandler(moduleRepo)
	postHandler := handlers.NewPostHandler(postRepo)
	commentHandler := handlers.NewCommentHandler(commentRepo)

	// Initialize router
	router := routes.NewRouter(moduleHandler, postHandler, commentHandler)

	// Start server
	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
