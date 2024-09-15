package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tender/internal/models"
	"tender/internal/services"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type OrganizationResponsibleHandler struct {
	ResponsibleService services.OrganizationResponsibleServiceInterface
}

// Назначение ответственного за организацию
func (h *OrganizationResponsibleHandler) AssignResponsible(w http.ResponseWriter, r *http.Request) {
	var responsible models.OrganizationResponsible
	if err := json.NewDecoder(r.Body).Decode(&responsible); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	fmt.Printf("Received responsible: %+v\n", r.Body)
	if err := h.ResponsibleService.AssignResponsible(&responsible); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responsible)
}

// Получение ответственного по ID
func (h *OrganizationResponsibleHandler) GetResponsible(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	responsibleID, err := uuid.Parse(vars["responsibleID"])
	if err != nil {
		http.Error(w, "Invalid Responsible ID", http.StatusBadRequest)
		return
	}

	responsible, err := h.ResponsibleService.GetResponsibleByID(responsibleID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responsible)
}

// Удаление ответственного за организацию
func (h *OrganizationResponsibleHandler) RemoveResponsible(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	responsibleID, err := uuid.Parse(vars["responsibleID"])
	if err != nil {
		http.Error(w, "Invalid Responsible ID", http.StatusBadRequest)
		return
	}

	if err := h.ResponsibleService.RemoveResponsible(responsibleID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Responsible removed"})
}
