package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"modulate/backend/internal/handlers"
	"modulate/backend/internal/middleware"
)

func NewRouter(
    moduleHandler *handlers.ModuleHandler,
    postHandler *handlers.PostHandler,
    commentHandler *handlers.CommentHandler,
    userHandler *handlers.UserHandler,
) *mux.Router {

    router := mux.NewRouter()
    router.StrictSlash(true)

    // Global middleware
    router.Use(middleware.LoggingMiddleware)
    router.Use(middleware.RecoveryMiddleware)

    // --------------------
    // Public (read-only)
    // --------------------
    // Modules
    router.HandleFunc("/modules", moduleHandler.List).
        Methods("GET")
    router.HandleFunc("/modules/{id}", moduleHandler.GetByID).
        Methods("GET")

    // Posts
    router.HandleFunc("/modules/{moduleID}/posts", postHandler.ListByModule).
        Methods("GET")
    router.HandleFunc("/posts/{postID}", postHandler.GetByID).
        Methods("GET")

    // Comments
    router.HandleFunc("/posts/{postID}/comments/thread", commentHandler.Thread).
        Methods("GET")

    // --------------------
    // Authentication
    // --------------------
    router.HandleFunc("/auth/register", userHandler.Register).
        Methods("POST")
    router.HandleFunc("/auth/login", userHandler.Login).
        Methods("POST")
    router.HandleFunc("/auth/refresh", userHandler.Refresh).
        Methods("POST")
    router.HandleFunc("/auth/logout", userHandler.Logout).
        Methods("POST")

    // --------------------
    // Protected (write)
    // --------------------
    api := router.PathPrefix("/api").Subrouter()
    api.Use(middleware.AuthMiddleware)

    // Posts
    api.HandleFunc("/modules/{moduleID}/posts", postHandler.Create).
        Methods("POST", "OPTIONS")
    api.HandleFunc("/posts/{postID}", postHandler.Update).
        Methods("PUT", "OPTIONS")
    api.HandleFunc("/posts/{postID}", postHandler.Delete).
        Methods("DELETE", "OPTIONS")

    // Comments
    api.HandleFunc("/posts/{postID}/comments", commentHandler.Create).
        Methods("POST", "OPTIONS")
    api.HandleFunc("/comments/{commentID}", commentHandler.Update).
        Methods("PUT", "OPTIONS")
    api.HandleFunc("/comments/{commentID}", commentHandler.Delete).
        Methods("DELETE", "OPTIONS")

    // --------------------
    // Health
    // --------------------
    router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    }).Methods("GET")

    return router
}
