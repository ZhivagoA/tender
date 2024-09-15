package repositories

import (
	"tender/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FeedbackRepository struct {
	DB *gorm.DB
}

func (r *FeedbackRepository) Create(feedback *models.Feedback) error {
	return r.DB.Create(feedback).Error
}

func (r *FeedbackRepository) GetByID(id uuid.UUID) (*models.Feedback, error) {
	var feedback models.Feedback
	err := r.DB.First(&feedback, "id = ?", id).Error
	return &feedback, err
}

func (r *FeedbackRepository) Update(feedback *models.Feedback) error {
	return r.DB.Save(feedback).Error
}

func (r *FeedbackRepository) Delete(id uuid.UUID) error {
	return r.DB.Delete(&models.Feedback{}, id).Error
}
