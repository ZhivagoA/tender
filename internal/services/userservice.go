package services

import (
	"tender/internal/models"
	"tender/internal/repositories"

	"github.com/google/uuid"
)

type UserServiceInterface interface {
	CreateUser(user *models.User) error
	GetUser(id uuid.UUID) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id uuid.UUID) error
}

type UserService struct {
	Repo *repositories.UserRepository
}

func (s *UserService) CreateUser(user *models.User) error {
	return s.Repo.Create(user)
}

func (s *UserService) GetUser(id uuid.UUID) (*models.User, error) {
	return s.Repo.GetByID(id)
}

func (s *UserService) UpdateUser(user *models.User) error {
	return s.Repo.Update(user)
}

func (s *UserService) DeleteUser(id uuid.UUID) error {
	return s.Repo.Delete(id)
}
