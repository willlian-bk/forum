package service

import (
	"github.com/Akezhan1/forum/internal/app/models"
	"github.com/Akezhan1/forum/internal/app/repository"
)

type User interface {
	Create(*models.User) (int, int, error)
	Authorization(string, string) (*models.User, error)
}

type Service struct {
	User
}

func NewService(r *repository.Repository) *Service {
	return &Service{
		User: NewUserService(r.User),
	}
}
