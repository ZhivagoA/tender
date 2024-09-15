package repositories

import (
	"errors"
	"fmt"
	"tender/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganizationRepository struct {
	DB *gorm.DB
}

func (r *OrganizationRepository) Create(organization *models.Organization) error {
	return r.DB.Create(organization).Error
}

func (r *OrganizationRepository) GetByID(id uuid.UUID) (*models.Organization, error) {
	var organization models.Organization
	if err := r.DB.First(&organization, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("organization not found %v", id)
		}
		return nil, err
	}
	return &organization, nil
}

func (r *OrganizationRepository) Update(organization *models.Organization) error {
	return r.DB.Save(organization).Error
}

func (r *OrganizationRepository) Delete(id uuid.UUID) error {
	return r.DB.Delete(&models.Organization{}, id).Error
}
