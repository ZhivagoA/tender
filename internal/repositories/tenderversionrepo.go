package repositories

import (
	"tender/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TenderVersionRepository struct {
	DB *gorm.DB
}

// Создание новой версии тендера
func (r *TenderVersionRepository) Create(version *models.TenderVersion) error {
	return r.DB.Create(version).Error
}

// Получение версии тендера по ID и номеру версии
func (r *TenderVersionRepository) GetByTenderIDAndVersion(tenderID uuid.UUID, versionNumber int) (*models.TenderVersion, error) {
	var version models.TenderVersion
	err := r.DB.Where("tender_id = ? AND version = ?", tenderID, versionNumber).First(&version).Error
	return &version, err
}
