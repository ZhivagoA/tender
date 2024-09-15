package repositories

import (
	"tender/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BidVersionRepository struct {
	DB *gorm.DB
}

// Создание новой версии предложения
func (r *BidVersionRepository) Create(version *models.BidVersion) error {
	return r.DB.Create(version).Error
}

// Получение версии предложения по ID и номеру версии
func (r *BidVersionRepository) GetByBidIDAndVersion(bidID uuid.UUID, versionNumber int) (*models.BidVersion, error) {
	var version models.BidVersion
	err := r.DB.Where("bid_id = ? AND version = ?", bidID, versionNumber).First(&version).Error
	return &version, err
}
