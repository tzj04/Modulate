package handlers

import (
	"encoding/json"
	"net/http"

	"modulate/backend/internal/repositories"
)

type ModuleHandler struct {
	Repo repositories.ModuleRepository
}

func NewModuleHandler(repo repositories.ModuleRepository) *ModuleHandler {
	return &ModuleHandler{Repo: repo}
}

func (h *ModuleHandler) List(w http.ResponseWriter, r *http.Request) {
	modules, err := h.Repo.ListAll()
	if err != nil {
		http.Error(w, "failed to fetch modules", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(modules)
}
