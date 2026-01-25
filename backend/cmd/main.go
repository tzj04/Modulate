package main

import (
	"log"
	"net/http"
	"os"
    "fmt"

	"github.com/joho/godotenv"

	"modulate/backend/internal/db"
	"modulate/backend/internal/handlers"
	"modulate/backend/internal/middleware"
	"modulate/backend/internal/repositories/postgres"
	"modulate/backend/internal/routes"

	_ "github.com/lib/pq"
)

func main() {
    fmt.Println("GOPROD: App is starting...")
    
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
    err = db.RunMigrations()
    if err != nil {
        log.Printf("Migration notice: %v", err) 
    }

    defer db.Close()

    fmt.Println("Initialising Repos...")
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

    fmt.Println("Initialising Handlers...")
    // Initialize router
    router := routes.NewRouter(
        moduleHandler, 
        postHandler, 
        commentHandler,
        userHandler,
    )
    handler := middleware.CORSMiddleware(router)

    fmt.Println("Starting Server...")
    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // Fallback for local development
    }

    log.Printf("Server running on port %s", port)
    if err := http.ListenAndServe("0.0.0.0:"+port, handler); err != nil {
        log.Fatal(err)
    }
}