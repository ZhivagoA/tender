package handlers

import (
	"encoding/json"
	"net/http"
	"tender/internal/models"
	"tender/internal/services"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type TenderHandler struct {
	TenderService services.TenderServiceInterface
}

// Создание тендера
func (h *TenderHandler) CreateTender(w http.ResponseWriter, r *http.Request) {
	var tender models.Tender
	if err := json.NewDecoder(r.Body).Decode(&tender); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	responsibleUserID := tender.ResponsibleUserID

	if err := h.TenderService.CreateTender(&tender, responsibleUserID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tender)
}

// Получение списка всех тендеров
func (h *TenderHandler) ListTenders(w http.ResponseWriter, r *http.Request) {
	tenders, err := h.TenderService.ListTenders()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tenders)
}

// Получение списка тендеров текущего пользователя
func (h *TenderHandler) MyTenders(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		http.Error(w, "Missing User ID", http.StatusBadRequest)
		return
	}

	// Парсинг userID
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	tenders, err := h.TenderService.ListTendersByUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tenders)
}

// Получение статуса тендера
func (h *TenderHandler) GetTenderStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenderID, err := uuid.Parse(vars["tenderID"])
	if err != nil {
		http.Error(w, "Invalid Tender ID", http.StatusBadRequest)
		return
	}

	status, err := h.TenderService.GetTenderStatus(tenderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": status})
}

// Редактирование тендера
func (h *TenderHandler) EditTender(w http.ResponseWriter, r *http.Request) {
	var tender models.Tender
	if err := json.NewDecoder(r.Body).Decode(&tender); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	responsibleUserID := tender.ResponsibleUserID

	if err := h.TenderService.EditTender(&tender, responsibleUserID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tender)
}

// Откат версии тендера
func (h *TenderHandler) RollbackTenderVersion(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		TenderID uuid.UUID `json:"tender_id"`
		Version  int       `json:"version"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.TenderService.RevertTender(payload.TenderID, payload.Version); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Tender version reverted"})
}

func (h *TenderHandler) ListTendersByUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		http.Error(w, "Missing User ID", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	tenders, err := h.TenderService.ListTendersByUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tenders)
}
