package services

import (
	"errors"
	"fmt"
	"tender/internal/models"
	"tender/internal/repositories"

	"github.com/google/uuid"
)

type TenderServiceInterface interface {
	CreateTender(tender *models.Tender, responsibleUserID uuid.UUID) error
	PublishTender(tenderID uuid.UUID, responsibleUserID uuid.UUID) error
	CloseTender(tenderID uuid.UUID, responsibleUserID uuid.UUID) error
	EditTender(tender *models.Tender, responsibleUserID uuid.UUID) error
	RevertTender(tenderID uuid.UUID, versionNumber int) error
	ListTenders() ([]*models.Tender, error)
	ListTendersByUser(userID uuid.UUID) ([]*models.Tender, error)
	GetTenderStatus(tenderID uuid.UUID) (string, error)
}

type TenderService struct {
	Repo                        *repositories.TenderRepository
	VersionRepo                 *repositories.TenderVersionRepository
	OrganizationRepo            *repositories.OrganizationRepository
	OrganizationResponsibleRepo *repositories.OrganizationResponsibleRepository
}

// Создание тендера
func (s *TenderService) CreateTender(tender *models.Tender, responsibleUserID uuid.UUID) error {
	organization, err := s.OrganizationRepo.GetByID(tender.OrganizationID)
	if err != nil {
		return fmt.Errorf("organization not found: %v", err)
	}
	fmt.Println(responsibleUserID, organization.ID)
	isResponsible, err := s.OrganizationResponsibleRepo.IsUserResponsibleForOrganization(responsibleUserID, organization.ID)
	if err != nil {
		return fmt.Errorf("error checking user responsibility: %v", err)
	}
	if !isResponsible {
		return errors.New("user is not responsible for the organization")
	}

	tender.Status = "CREATED"
	tender.Version = 1
	return s.Repo.Create(tender)
}

// Публикация тендера
func (s *TenderService) PublishTender(tenderID uuid.UUID, responsibleUserID uuid.UUID) error {
	tender, err := s.Repo.GetByID(tenderID)
	if err != nil {
		return err
	}

	_, err = s.OrganizationRepo.GetByID(responsibleUserID)
	if err != nil {
		return errors.New("user is not responsible for the organization")
	}

	if tender.Status != "CREATED" {
		return errors.New("only tenders with 'CREATED' status can be published")
	}

	tender.Status = "PUBLISHED"
	return s.Repo.Update(tender)
}

// Закрытие тендера
func (s *TenderService) CloseTender(tenderID uuid.UUID, responsibleUserID uuid.UUID) error {
	tender, err := s.Repo.GetByID(tenderID)
	if err != nil {
		return err
	}

	_, err = s.OrganizationRepo.GetByID(responsibleUserID)
	if err != nil {
		return errors.New("user is not responsible for the organization")
	}

	tender.Status = "CLOSED"
	return s.Repo.Update(tender)
}

// Редактирование тендера
func (s *TenderService) EditTender(tender *models.Tender, responsibleUserID uuid.UUID) error {
	responsibleRepo, err := s.OrganizationResponsibleRepo.GetByID(responsibleUserID)
	if err != nil {
		return errors.New("no user with this id was found")
	}
	if responsibleRepo.OrganizationID != tender.OrganizationID {
		return errors.New("user is not responsible for the organization")
	}

	if err := s.saveTenderVersion(tender); err != nil {
		return err
	}

	tender.Version += 1
	return s.Repo.Update(tender)
}

func (s *TenderService) saveTenderVersion(tender *models.Tender) error {
	version := models.TenderVersion{
		TenderID:          tender.ID,
		Title:             tender.Title,
		Description:       tender.Description,
		OrganizationID:    tender.OrganizationID,
		ResponsibleUserID: tender.ResponsibleUserID,
		Status:            tender.Status,
		Version:           tender.Version,
	}
	return s.VersionRepo.Create(&version)
}

func (s *TenderService) RevertTender(tenderID uuid.UUID, versionNumber int) error {
	version, err := s.VersionRepo.GetByTenderIDAndVersion(tenderID, versionNumber)
	if err != nil {
		return err
	}

	tender, err := s.Repo.GetByID(tenderID)
	if err != nil {
		return err
	}

	tender.Title = version.Title
	tender.Description = version.Description
	tender.Status = version.Status
	tender.Version = version.Version

	return s.Repo.Update(tender)
}

// Список всех тендеров
func (s *TenderService) ListTenders() ([]*models.Tender, error) {
	return s.Repo.ListAll()
}

// Список тендеров пользователя
func (s *TenderService) ListTendersByUser(userID uuid.UUID) ([]*models.Tender, error) {
	return s.Repo.ListByUserID(userID)
}

// Получение статуса тендера
func (s *TenderService) GetTenderStatus(tenderID uuid.UUID) (string, error) {
	tender, err := s.Repo.GetByID(tenderID)
	if err != nil {
		return "", err
	}
	return tender.Status, nil
}
