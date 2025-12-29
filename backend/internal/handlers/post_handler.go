package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"modulate/backend/internal/models"
	"modulate/backend/internal/repositories"
)

type PostHandler struct {
	Repo repositories.PostRepository
}

func NewPostHandler(repo repositories.PostRepository) *PostHandler {
	return &PostHandler{Repo: repo}
}

type createPostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	UserID  int64  `json:"user_id"`
}

func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	moduleID, err := strconv.ParseInt(vars["moduleID"], 10, 64)
	if err != nil {
		http.Error(w, "invalid module id", http.StatusBadRequest)
		return
	}

	var req createPostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	Post := &models.Post{
		ModuleID: moduleID,
		UserID:   req.UserID,
		Title:    req.Title,
		Content:  req.Content,
	}

	if err := h.Repo.Create(Post); err != nil {
		http.Error(w, "failed to create Post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Post)
}

func (h *PostHandler) ListByModule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	moduleID, err := strconv.ParseInt(vars["moduleID"], 10, 64)
	if err != nil {
		http.Error(w, "invalid module id", http.StatusBadRequest)
		return
	}

	Posts, err := h.Repo.ListByModule(moduleID)
	if err != nil {
		http.Error(w, "failed to fetch Posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Posts)
}
