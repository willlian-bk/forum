package service

import (
	"time"

	"github.com/Akezhan1/forum/internal/app/models"
	"github.com/Akezhan1/forum/internal/app/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo}
}

func (us *UserService) Create(user *models.User) (int64, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return -1, err
	}

	user.Password = string(hashPassword)
	user.Role = "user"
	user.CreatedDate = time.Now()

	return us.repo.Create(user)
}
