package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

func (h *ModuleHandler) GetByID(w http.ResponseWriter, r *http.Request) {
    // Extract the "id" from the URL path
    vars := mux.Vars(r)
    idStr := vars["id"]

    // Convert string ID to integer
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "invalid module ID", http.StatusBadRequest)
        return
    }

    // Fetch from repository
    module, err := h.Repo.GetByID(id)
    if err != nil {
        // If the module doesn't exist, return 404
        http.Error(w, "module not found", http.StatusNotFound)
        return
    }

    // Return as JSON
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(module)
}