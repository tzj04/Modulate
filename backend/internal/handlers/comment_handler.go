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

type CommentHandler struct {
    Repo repositories.CommentRepository
}

func NewCommentHandler(repo repositories.CommentRepository) *CommentHandler {
    return &CommentHandler{Repo: repo}
}

type createCommentRequest struct {
    Content         string `json:"content"`
    ParentCommentID *int64 `json:"parent_comment_id"`
}

func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    postID, err := strconv.ParseInt(vars["postID"], 10, 64)
    if err != nil {
        http.Error(w, "invalid post id", http.StatusBadRequest)
        return
    }

    var req createCommentRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    if req.Content == "" {
        http.Error(w, "content cannot be empty", http.StatusBadRequest)
        return
    }

    // Extract UserID from Context (set by AuthMiddleware)
    userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
    if !ok {
        http.Error(w, "unauthorized", http.StatusUnauthorized)
        return
    }

    comment := &models.Comment{
        PostID:          postID,
        UserID:          userID,
        ParentCommentID: req.ParentCommentID,
        Content:         req.Content,
    }

    if err := h.Repo.Create(comment); err != nil {
        http.Error(w, "failed to create comment", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    _ = json.NewEncoder(w).Encode(comment)
}

func (h *CommentHandler) Thread(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    postID, err := strconv.ParseInt(vars["postID"], 10, 64)
    if err != nil {
        http.Error(w, "invalid post id", http.StatusBadRequest)
        return
    }

    // Fetch the list from the database
    comments, err := h.Repo.ListThreadByPost(postID)
    if err != nil {
        http.Error(w, "failed to fetch comments", http.StatusInternalServerError)
        return
    }

    // React buildCommentTree hook handles the nesting.
    w.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(w).Encode(comments)
}

func (h *CommentHandler) Delete(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    commentID, err := strconv.ParseInt(vars["commentID"], 10, 64)
    if err != nil {
        http.Error(w, "invalid comment id", http.StatusBadRequest)
        return
    }

    if err := h.Repo.SoftDelete(commentID); err != nil {
        http.Error(w, "failed to delete comment", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}