package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"modulate/backend/internal/middleware"
	"modulate/backend/internal/models"
	"modulate/backend/internal/repositories"
)

type PostHandler struct {
	Repo repositories.PostRepository
}

func NewPostHandler(repo repositories.PostRepository) *PostHandler {
	return &PostHandler{Repo: repo}
}

func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Get UserID from Context, key defined in middleware
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		http.Error(w, "Unauthorized: could not get user from context", http.StatusUnauthorized)
		return
	}

	// Decode the body
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Get ModuleID from URL
	vars := mux.Vars(r)
	moduleID, err := strconv.Atoi(vars["moduleID"])
	if err != nil {
		http.Error(w, "invalid module id", http.StatusBadRequest)
		return
	}

	// Save to Repo
	post := &models.Post{
		ModuleID: int64(moduleID),
		UserID:   userID,
		Title:    req.Title,
		Content:  req.Content,
	}

	if err := h.Repo.Create(post); err != nil {
		http.Error(w, "failed to create post: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func (h *PostHandler) ListByModule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	moduleID, err := strconv.ParseInt(vars["moduleID"], 10, 64)
	if err != nil {
		http.Error(w, "invalid module id", http.StatusBadRequest)
		return
	}

	posts, err := h.Repo.ListByModule(moduleID)
	if err != nil {
		http.Error(w, "failed to fetch posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (h *PostHandler) GetByID(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    // use "postID" to match the route: /posts/{postID}
    id, err := strconv.ParseInt(vars["postID"], 10, 64)
    if err != nil {
        http.Error(w, "invalid post id", http.StatusBadRequest)
        return
    }

    post, err := h.Repo.GetByID(id)
    if err != nil {
        // If the repo returns sql.ErrNoRows
        http.Error(w, "post not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(post)
}