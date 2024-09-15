package services

import (
	"errors"
	"fmt"
	"tender/internal/models"
	"tender/internal/repositories"
	"time"

	"github.com/google/uuid"
)

type BidServiceInterface interface {
	CreateBid(bid *models.Bid) error
	PublishBid(bidID uuid.UUID, responsibleUserID uuid.UUID) error
	CancelBid(bidID uuid.UUID, responsibleUserID uuid.UUID) error
	EditBid(bid *models.Bid, responsibleUserID uuid.UUID) error
	ApproveBid(bidID uuid.UUID, tenderID uuid.UUID, responsibleUserID uuid.UUID) error
	RejectBid(bidID uuid.UUID, tenderID uuid.UUID, responsibleUserID uuid.UUID) error
	ApproveBidWithQuorum(bidID uuid.UUID, tenderID uuid.UUID, responsibleUserID uuid.UUID, decision string) error
	ListBids() ([]*models.Bid, error)
	ListBidsByUser(userID uuid.UUID) ([]*models.Bid, error)
	GetBidReviews(bidID uuid.UUID) ([]*models.Feedback, error)
	RevertBid(bidID uuid.UUID, versionNumber int, responsibleUserID uuid.UUID) error
	GetBidStatus(bidID uuid.UUID) (string, error)
	LeaveFeedbackOnBid(feedback *models.Feedback) error
}

type BidService struct {
	Repo                        *repositories.BidRepository
	TenderRepo                  *repositories.TenderRepository
	VersionRepo                 *repositories.BidVersionRepository
	OrganizationResponsibleRepo *repositories.OrganizationResponsibleRepository
}

// Создание предложения
func (s *BidService) CreateBid(bid *models.Bid) error {
	bid.Status = "CREATED"
	bid.Version = 1
	bid.ApprovalCount = 0
	return s.Repo.Create(bid)
}

// Публикация предложения
func (s *BidService) PublishBid(bidID uuid.UUID, responsibleUserID uuid.UUID) error {
	bid, err := s.Repo.GetByID(bidID)
	if err != nil {
		return err
	}

	_, err = s.OrganizationResponsibleRepo.GetByID(responsibleUserID)
	if err != nil {
		return errors.New("user is not responsible for the organization")
	}

	bid.Status = "PUBLISHED"
	return s.Repo.Update(bid)
}

// Отмена предложения
func (s *BidService) CancelBid(bidID uuid.UUID, responsibleUserID uuid.UUID) error {
	bid, err := s.Repo.GetByID(bidID)
	if err != nil {
		return err
	}

	_, err = s.OrganizationResponsibleRepo.GetByID(responsibleUserID)
	if err != nil {
		return errors.New("user is not responsible for the organization")
	}

	bid.Status = "CANCELED"
	return s.Repo.Update(bid)
}

// Редактирование предложения
func (s *BidService) EditBid(bid *models.Bid, responsibleUserID uuid.UUID) error {
	currentBid, err := s.Repo.GetByID(bid.ID)
	if err != nil {
		return errors.New("trying to update unexisted bid")
	}
	if currentBid.UserID != responsibleUserID {
		return errors.New("user is not responsible for the organization")
	}
	fmt.Println(bid.TenderID)
	fmt.Println(currentBid.TenderID)
	if err := s.saveBidVersion(currentBid); err != nil {
		return err
	}

	bid.Version += currentBid.Version + 1
	bid.TenderID = currentBid.TenderID
	bid.ApprovalCount = currentBid.ApprovalCount
	bid.UpdatedAt = time.Now()
	bid.CreatedAt = currentBid.CreatedAt
	return s.Repo.Update(bid)
}

// Сохранение версии предложения
func (s *BidService) saveBidVersion(bid *models.Bid) error {
	version := models.BidVersion{
		BidID:         bid.ID,
		Amount:        bid.Amount,
		TenderID:      bid.TenderID,
		UserID:        bid.UserID,
		Status:        bid.Status,
		Version:       bid.Version,
		ApprovalCount: bid.ApprovalCount,
	}
	return s.VersionRepo.Create(&version)
}

// Согласование предложения
func (s *BidService) ApproveBid(bidID uuid.UUID, tenderID uuid.UUID, responsibleUserID uuid.UUID) error {
	tender, err := s.TenderRepo.GetByID(tenderID)
	if err != nil {
		return err
	}

	if tender.Status == "CLOSED" {
		return errors.New("tender is already closed")
	}

	bid, err := s.Repo.GetByID(bidID)
	if err != nil {
		return err
	}

	_, err = s.OrganizationResponsibleRepo.GetByID(responsibleUserID)
	if err != nil {
		return errors.New("user is not responsible for the organization")
	}

	bid.Status = "APPROVED"
	tender.Status = "CLOSED"

	if err := s.Repo.Update(bid); err != nil {
		return err
	}
	return s.TenderRepo.Update(tender)
}

// Отклонение предложения
func (s *BidService) RejectBid(bidID uuid.UUID, tenderID uuid.UUID, responsibleUserID uuid.UUID) error {
	bid, err := s.Repo.GetByID(bidID)
	if err != nil {
		return err
	}

	_, err = s.OrganizationResponsibleRepo.GetByID(responsibleUserID)
	if err != nil {
		return errors.New("user is not responsible for the organization")
	}

	bid.Status = "REJECTED"
	return s.Repo.Update(bid)
}

// Логика одобрения с кворумом
func (s *BidService) ApproveBidWithQuorum(bidID uuid.UUID, tenderID uuid.UUID, responsibleUserID uuid.UUID, decision string) error {
	tender, err := s.TenderRepo.GetByID(tenderID)
	if err != nil {
		return err
	}

	if tender.Status == "CLOSED" {
		return errors.New("tender is already closed")
	}

	bid, err := s.Repo.GetByID(bidID)
	if err != nil {
		return err
	}

	responsibleUsers, err := s.OrganizationResponsibleRepo.GetResponsiblesByOrganizationID(tender.OrganizationID)
	if err != nil {
		return err
	}

	quorum := min(3, len(responsibleUsers))

	if decision == "reject" {
		bid.Status = "REJECTED"
		return s.Repo.Update(bid)
	}

	bid.ApprovalCount += 1

	if bid.ApprovalCount >= quorum {
		bid.Status = "APPROVED"
		tender.Status = "CLOSED"
		if err := s.TenderRepo.Update(tender); err != nil {
			return err
		}
	}

	return s.Repo.Update(bid)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Получение всех предложений
func (s *BidService) ListBids() ([]*models.Bid, error) {
	return s.Repo.ListAll()
}

// Получение предложений пользователя
func (s *BidService) ListBidsByUser(userID uuid.UUID) ([]*models.Bid, error) {
	return s.Repo.ListByUserID(userID)
}

// Получение отзывов на предложение
func (s *BidService) GetBidReviews(bidID uuid.UUID) ([]*models.Feedback, error) {
	return s.Repo.GetReviewsForBid(bidID)
}

// Откат версии предложения
func (s *BidService) RevertBid(bidID uuid.UUID, versionNumber int, responsibleUserID uuid.UUID) error {
	version, err := s.VersionRepo.GetByBidIDAndVersion(bidID, versionNumber)
	if err != nil {
		return err
	}

	bid, err := s.Repo.GetByID(bidID)
	if err != nil {
		return err
	}

	bid.Amount = version.Amount
	bid.Status = version.Status
	bid.Version = version.Version
	bid.ApprovalCount = version.ApprovalCount

	return s.Repo.Update(bid)
}

// Получение статуса предложения
func (s *BidService) GetBidStatus(bidID uuid.UUID) (string, error) {
	bid, err := s.Repo.GetByID(bidID)
	if err != nil {
		return "", err
	}
	return bid.Status, nil
}

// Оставление отзыва на предложение
func (s *BidService) LeaveFeedbackOnBid(feedback *models.Feedback) error {
	return s.Repo.CreateFeedback(feedback)
}
