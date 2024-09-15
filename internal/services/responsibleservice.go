package services

import (
	"fmt"
	"tender/internal/models"
	"tender/internal/repositories"

	"github.com/google/uuid"
)

type OrganizationResponsibleServiceInterface interface {
	AssignResponsible(responsible *models.OrganizationResponsible) error
	GetResponsibleByID(id uuid.UUID) (*models.OrganizationResponsible, error)
	RemoveResponsible(id uuid.UUID) error
}

type OrganizationResponsibleService struct {
	Repo             *repositories.OrganizationResponsibleRepository
	OrganizationRepo *repositories.OrganizationRepository
	UserRepo         *repositories.UserRepository
}

func (s *OrganizationResponsibleService) AssignResponsible(responsible *models.OrganizationResponsible) error {
	_, err := s.OrganizationRepo.GetByID(responsible.OrganizationID)
	if err != nil {
		return fmt.Errorf("organization not found: %v", err)
	}
	_, err = s.UserRepo.GetByID(responsible.UserID)
	if err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	return s.Repo.Create(responsible)
}

func (s *OrganizationResponsibleService) GetResponsibleByID(id uuid.UUID) (*models.OrganizationResponsible, error) {
	return s.Repo.GetByID(id)
}

func (s *OrganizationResponsibleService) RemoveResponsible(id uuid.UUID) error {
	return s.Repo.Delete(id)
}
