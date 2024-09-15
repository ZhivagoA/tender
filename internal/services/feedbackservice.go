package services

import (
	"errors"
	"tender/internal/models"
	"tender/internal/repositories"

	"github.com/google/uuid"
)

type FeedbackServiceInterface interface {
	AddFeedback(feedback *models.Feedback) error
	GetFeedback(id uuid.UUID) (*models.Feedback, error)
	DeleteFeedback(id uuid.UUID) error
}

type FeedbackService struct {
	Repo                        *repositories.FeedbackRepository
	TenderRepo                  *repositories.TenderRepository
	OrganizationResponsibleRepo *repositories.OrganizationResponsibleRepository
}

func (s *FeedbackService) LeaveFeedbackOnBid(feedback *models.Feedback, responsibleUserID uuid.UUID) error {
	_, err := s.OrganizationResponsibleRepo.GetByID(responsibleUserID)
	if err != nil {
		return errors.New("user is not responsible for the organization")
	}

	return s.Repo.Create(feedback)
}

func (s *FeedbackService) GetReviewsForUserInTender(userID uuid.UUID, tenderID uuid.UUID, responsibleUserID uuid.UUID) ([]*models.Feedback, error) {
	tender, err := s.TenderRepo.GetByID(tenderID)
	if err != nil {
		return nil, err
	}

	_, err = s.OrganizationResponsibleRepo.GetResponsiblesByOrganizationID(tender.OrganizationID)
	if err != nil {
		return nil, errors.New("user is not responsible for the organization")
	}

	var feedbacks []*models.Feedback
	err = s.Repo.DB.Where("user_id = ? AND tender_id = ?", userID, tenderID).Find(&feedbacks).Error
	return feedbacks, err
}

func (s *FeedbackService) GetFeedback(id uuid.UUID) (*models.Feedback, error) {
	return s.Repo.GetByID(id)
}

func (s *FeedbackService) DeleteFeedback(id uuid.UUID) error {
	return s.Repo.Delete(id)
}
