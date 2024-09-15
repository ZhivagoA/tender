package repositories

import (
	"fmt"
	"tender/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TenderRepository struct {
	DB *gorm.DB
}

func (r *TenderRepository) Create(tender *models.Tender) error {
	return r.DB.Create(tender).Error
}

func (r *TenderRepository) Update(tender *models.Tender) error {
	return r.DB.Save(tender).Error
}

func (r *TenderRepository) Delete(id uuid.UUID) error {
	return r.DB.Delete(&models.Tender{}, id).Error
}
func (r *TenderRepository) ListAll() ([]*models.Tender, error) {
	var tenders []*models.Tender
	if err := r.DB.Find(&tenders).Error; err != nil {
		return nil, err
	}
	return tenders, nil
}

// Получение тендеров по ID пользователя
func (r *TenderRepository) ListByUserID(userID uuid.UUID) ([]*models.Tender, error) {
	var tenders []*models.Tender
	if err := r.DB.Where("responsible_user_id = ?", userID).Find(&tenders).Error; err != nil {
		return nil, err
	}
	fmt.Println(tenders)
	return tenders, nil
}

// Получение тендера по ID
func (r *TenderRepository) GetByID(tenderID uuid.UUID) (*models.Tender, error) {
	var tender models.Tender
	if err := r.DB.First(&tender, "id = ?", tenderID).Error; err != nil {
		return nil, err
	}
	return &tender, nil
}
