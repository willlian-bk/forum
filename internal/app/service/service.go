package service

import (
	"github.com/Akezhan1/forum/internal/app/models"
	"github.com/Akezhan1/forum/internal/app/repository"
)

type User interface {
	Create(*models.User) (int, int, error)
	Authorization(string, string) (*models.Session, error)
	Logout(string) error
	IsValidToken(string) bool
	GetUserIDByToken(string) (int, error)
}

type Post interface {
	Create(*models.Post) (int, int, error)
	Get(int) (*models.Post, error)
	GetValidCategories() ([]string, error)
}

type Service struct {
	User
	Post
}

func NewService(r *repository.Repository) *Service {
	return &Service{
		User: NewUserService(r.User),
		Post: NewPostService(r.Post),
	}
}
