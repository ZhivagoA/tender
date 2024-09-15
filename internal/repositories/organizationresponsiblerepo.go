package repositories

import (
	"fmt"
	"tender/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganizationResponsibleRepository struct {
	DB *gorm.DB
}

func (r *OrganizationResponsibleRepository) Create(responsible *models.OrganizationResponsible) error {
	return r.DB.Create(responsible).Error
}

func (r *OrganizationResponsibleRepository) GetResponsiblesByOrganizationID(orgID uuid.UUID) ([]*models.OrganizationResponsible, error) {
	var responsibles []*models.OrganizationResponsible
	err := r.DB.Where("organization_id = ?", orgID).Find(&responsibles).Error
	return responsibles, err
}

func (r *OrganizationResponsibleRepository) GetByID(id uuid.UUID) (*models.OrganizationResponsible, error) {
	var responsible models.OrganizationResponsible
	err := r.DB.First(&responsible, "id = ?", id).Error
	return &responsible, err
}

func (r *OrganizationResponsibleRepository) Update(responsible *models.OrganizationResponsible) error {
	return r.DB.Save(responsible).Error
}

func (r *OrganizationResponsibleRepository) Delete(id uuid.UUID) error {
	return r.DB.Delete(&models.OrganizationResponsible{}, id).Error
}

func (r *OrganizationResponsibleRepository) IsUserResponsibleForOrganization(userID, organizationID uuid.UUID) (bool, error) {
	var count int64
	fmt.Println(userID.String(), "repositories")
	err := r.DB.Model(&models.OrganizationResponsible{}).
		Where("id = ? AND organization_id = ?", userID, organizationID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
