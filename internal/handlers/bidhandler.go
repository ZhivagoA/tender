package handlers

import (
	"encoding/json"
	"net/http"
	"tender/internal/models"
	"tender/internal/services"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type BidHandler struct {
	BidService services.BidServiceInterface
}

// Создание предложения
func (h *BidHandler) CreateBid(w http.ResponseWriter, r *http.Request) {
	var bid models.Bid
	if err := json.NewDecoder(r.Body).Decode(&bid); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.BidService.CreateBid(&bid); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(bid)
}

// Публикация предложения
func (h *BidHandler) PublishBid(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bidID, err := uuid.Parse(vars["bidID"])
	if err != nil {
		http.Error(w, "Invalid Bid ID", http.StatusBadRequest)
		return
	}

	userIDStr := r.Header.Get("user_id")
	if userIDStr == "" {
		http.Error(w, "Missing User ID", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	if err := h.BidService.PublishBid(bidID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Bid published"})
}

// Отмена предложения
func (h *BidHandler) CancelBid(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bidID, err := uuid.Parse(vars["bidID"])
	if err != nil {
		http.Error(w, "Invalid Bid ID", http.StatusBadRequest)
		return
	}

	userIDStr := r.Header.Get("user_id")
	if userIDStr == "" {
		http.Error(w, "Missing User ID", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	if err := h.BidService.CancelBid(bidID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Bid canceled"})
}

// Редактирование предложения
func (h *BidHandler) EditBid(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bidID, err := uuid.Parse(vars["bidID"])
	if err != nil {
		http.Error(w, "Invalid Bid ID", http.StatusBadRequest)
		return
	}

	var bid models.Bid
	if err := json.NewDecoder(r.Body).Decode(&bid); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	bid.ID = bidID

	responsibleUserID := bid.UserID

	if err := h.BidService.EditBid(&bid, responsibleUserID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bid)
}

// Принятие решения по предложению с кворумом
func (h *BidHandler) SubmitDecision(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		BidID    uuid.UUID `json:"bid_id"`
		TenderID uuid.UUID `json:"tender_id"`
		Decision string    `json:"decision"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Пример идентификатора пользователя
	responsibleUserID := uuid.New()

	if err := h.BidService.ApproveBidWithQuorum(payload.BidID, payload.TenderID, responsibleUserID, payload.Decision); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Decision submitted"})
}

// Получение списка всех предложений
func (h *BidHandler) ListBids(w http.ResponseWriter, r *http.Request) {
	bids, err := h.BidService.ListBids()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bids)
}

// Получение предложений текущего пользователя
func (h *BidHandler) MyBids(w http.ResponseWriter, r *http.Request) {
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

	bids, err := h.BidService.ListBidsByUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bids)
}

// Получение статуса предложения
func (h *BidHandler) GetBidStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bidID, err := uuid.Parse(vars["bidID"])
	if err != nil {
		http.Error(w, "Invalid Bid ID", http.StatusBadRequest)
		return
	}

	status, err := h.BidService.GetBidStatus(bidID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": status})
}

// Откат версии предложения
func (h *BidHandler) RollbackBidVersion(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		BidID   uuid.UUID `json:"bid_id"`
		Version int       `json:"version"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	userIDStr := r.Header.Get("user_id")
	if userIDStr == "" {
		http.Error(w, "Missing User ID", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	if err := h.BidService.RevertBid(payload.BidID, payload.Version, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Bid version reverted"})
}

// Получение отзывов на предложение
func (h *BidHandler) GetBidReviews(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bidID, err := uuid.Parse(vars["bidID"])
	if err != nil {
		http.Error(w, "Invalid Bid ID", http.StatusBadRequest)
		return
	}

	reviews, err := h.BidService.GetBidReviews(bidID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reviews)
}

// Оставление отзыва на предложение
func (h *BidHandler) LeaveBidFeedback(w http.ResponseWriter, r *http.Request) {
	var feedback models.Feedback
	if err := json.NewDecoder(r.Body).Decode(&feedback); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.BidService.LeaveFeedbackOnBid(&feedback); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(feedback)
}
