package handlers

import (
	"encoding/json"
	"net/http"
	"tender/internal/models"
	"tender/internal/services"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type OrganizationHandler struct {
	OrganizationService services.OrganizationServiceInterface
}

// Создание организации
func (h *OrganizationHandler) CreateOrganization(w http.ResponseWriter, r *http.Request) {
	var organization models.Organization
	if err := json.NewDecoder(r.Body).Decode(&organization); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := h.OrganizationService.CreateOrganization(&organization); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(organization)
}

// Получение организации по ID
func (h *OrganizationHandler) GetOrganization(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orgID, err := uuid.Parse(vars["orgID"])
	if err != nil {
		http.Error(w, "Invalid Organization ID", http.StatusBadRequest)
		return
	}

	org, err := h.OrganizationService.GetOrganization(orgID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(org)
}

// Обновление данных организации
func (h *OrganizationHandler) UpdateOrganization(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orgID, err := uuid.Parse(vars["orgID"])
	if err != nil {
		http.Error(w, "Invalid Organization ID", http.StatusBadRequest)
		return
	}

	var organization models.Organization
	if err := json.NewDecoder(r.Body).Decode(&organization); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	organization.ID = orgID
	if err := h.OrganizationService.UpdateOrganization(&organization); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(organization)
}

// Удаление организации
func (h *OrganizationHandler) DeleteOrganization(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orgID, err := uuid.Parse(vars["orgID"])
	if err != nil {
		http.Error(w, "Invalid Organization ID", http.StatusBadRequest)
		return
	}

	if err := h.OrganizationService.DeleteOrganization(orgID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Organization deleted"})
}
