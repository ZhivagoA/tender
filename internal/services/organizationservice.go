package services

import (
	"tender/internal/models"
	"tender/internal/repositories"

	"github.com/google/uuid"
)

type OrganizationServiceInterface interface {
	CreateOrganization(org *models.Organization) error
	GetOrganization(id uuid.UUID) (*models.Organization, error)
	UpdateOrganization(org *models.Organization) error
	DeleteOrganization(id uuid.UUID) error
}

type OrganizationService struct {
	Repo *repositories.OrganizationRepository
}

func (s *OrganizationService) CreateOrganization(org *models.Organization) error {
	return s.Repo.Create(org)
}

func (s *OrganizationService) GetOrganization(id uuid.UUID) (*models.Organization, error) {
	return s.Repo.GetByID(id)
}

func (s *OrganizationService) UpdateOrganization(org *models.Organization) error {
	return s.Repo.Update(org)
}

func (s *OrganizationService) DeleteOrganization(id uuid.UUID) error {
	return s.Repo.Delete(id)
}
