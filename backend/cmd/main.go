package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"modulate/backend/internal/db"
	"modulate/backend/internal/handlers"
	"modulate/backend/internal/middleware"
	"modulate/backend/internal/repositories/postgres"
	"modulate/backend/internal/routes"

	_ "github.com/lib/pq"
)

func main() {
	// Load the .env file
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using system default envs")
    }

    // Get the DSN from the environment variable
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        log.Fatal("DATABASE_URL is not set in .env")
    }

    // Initialize DB
    err := db.Init(dsn)
    if err != nil {
        log.Fatalf("Error initializing database: %v", err)
    }
    defer db.Close()

    // Initialize repositories
    moduleRepo := &postgres.ModuleRepo{DB: db.DB}
    postRepo := &postgres.PostRepo{DB: db.DB}
    commentRepo := &postgres.CommentRepo{DB: db.DB}
    userRepo := &postgres.UserRepo{DB: db.DB}

    // Initialize handlers
    moduleHandler := handlers.NewModuleHandler(moduleRepo)
    postHandler := handlers.NewPostHandler(postRepo)
    commentHandler := handlers.NewCommentHandler(commentRepo)
    userHandler := handlers.NewUserHandler(userRepo)

    // Initialize router
    router := routes.NewRouter(
        moduleHandler, 
        postHandler, 
        commentHandler,
        userHandler,
    )
    handler := middleware.CORSMiddleware(router)

    // Start server
    log.Println("Server running on :8080")
    if err := http.ListenAndServe(":8080", handler); err != nil {
        log.Fatal(err)
    }
}