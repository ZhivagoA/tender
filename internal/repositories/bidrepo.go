package repositories

import (
	"tender/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BidRepository struct {
	DB *gorm.DB
}

func (r *BidRepository) Create(bid *models.Bid) error {
	return r.DB.Create(bid).Error
}

func (r *BidRepository) Update(bid *models.Bid) error {
	return r.DB.Save(bid).Error
}

func (r *BidRepository) Delete(id uuid.UUID) error {
	return r.DB.Delete(&models.Bid{}, id).Error
}

// Получение всех предложений
func (r *BidRepository) ListAll() ([]*models.Bid, error) {
	var bids []*models.Bid
	if err := r.DB.Find(&bids).Error; err != nil {
		return nil, err
	}
	return bids, nil
}

// Получение предложений по ID пользователя
func (r *BidRepository) ListByUserID(userID uuid.UUID) ([]*models.Bid, error) {
	var bids []*models.Bid
	if err := r.DB.Where("user_id = ?", userID).Find(&bids).Error; err != nil {
		return nil, err
	}
	return bids, nil
}

// Получение отзывов на предложение
func (r *BidRepository) GetReviewsForBid(bidID uuid.UUID) ([]*models.Feedback, error) {
	var reviews []*models.Feedback
	if err := r.DB.Where("bid_id = ?", bidID).Find(&reviews).Error; err != nil {
		return nil, err
	}
	return reviews, nil
}

// Получение предложения по ID
func (r *BidRepository) GetByID(bidID uuid.UUID) (*models.Bid, error) {
	var bid models.Bid
	if err := r.DB.First(&bid, "id = ?", bidID).Error; err != nil {
		return nil, err
	}
	return &bid, nil
}

// Добавление отзыва на предложение
func (r *BidRepository) CreateFeedback(feedback *models.Feedback) error {
	if err := r.DB.Create(feedback).Error; err != nil {
		return err
	}
	return nil
}
