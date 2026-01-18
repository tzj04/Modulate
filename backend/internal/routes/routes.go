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
) *mux.Router {

    router := mux.NewRouter()
    router.StrictSlash(true)

    // Global middleware
    router.Use(middleware.LoggingMiddleware)
    router.Use(middleware.RecoveryMiddleware)
    router.Use(middleware.CORSMiddleware)

    // --------------------
    // Public (read-only)
    // --------------------
    // Modules
    router.HandleFunc("/modules", moduleHandler.List).
        Methods("GET")

    // Posts
    router.HandleFunc("/modules/{moduleID}/posts", postHandler.ListByModule).
        Methods("GET")

    // Comments
    router.HandleFunc("/posts/{postID}/comments/thread", commentHandler.Thread).
        Methods("GET")

    // --------------------
    // Protected (write)
    // --------------------
    api := router.PathPrefix("/api").Subrouter()
    api.Use(middleware.AuthMiddleware)

    // Posts
    api.HandleFunc("/modules/{moduleID}/posts", postHandler.Create).
        Methods("POST")

    // Comments
    api.HandleFunc("/posts/{postID}/comments", commentHandler.Create).
        Methods("POST")
    api.HandleFunc("/comments/{commentID}", commentHandler.Delete).
        Methods("DELETE")

    // --------------------
    // Health
    // --------------------
    router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    }).Methods("GET")

    return router
}
