package service

import (
	"errors"
	"net/http"
	"regexp"
	"time"

	"github.com/Akezhan1/forum/internal/app/models"
	"github.com/Akezhan1/forum/internal/app/repository"
	"golang.org/x/crypto/bcrypt"

	sqlite "github.com/mattn/go-sqlite3"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo}
}

func (us *UserService) Create(user *models.User) (int, int, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return http.StatusInternalServerError, -1, err
	}

	if err := us.validateParams(user); err != nil {
		return http.StatusBadRequest, -1, err
	}

	user.Password = string(hashPassword)
	user.CreatedDate = time.Now()

	id, err := us.repo.Create(user)
	if err != nil {
		if sqliteErr, ok := err.(sqlite.Error); ok {
			if sqliteErr.ExtendedCode == sqlite.ErrConstraintUnique {
				return http.StatusBadRequest, -1, errors.New("User already created")
			}
		}
		return http.StatusInternalServerError, -1, err
	}

	return http.StatusOK, int(id), nil
}

func (us *UserService) validateParams(user *models.User) error {
	patternForEmail := `^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`

	ok, _ := regexp.MatchString(patternForEmail, user.Email)
	if !ok {
		return errors.New("Invalid Email")
	}

	if len(user.Password) < 6 {
		return errors.New("Invalid Password")
	}

	if user.Role != "user" {
		return errors.New("Invalid Role")
	}

	if user.Username == "" || len(user.Username) < 2 {
		return errors.New("Invalid Username")
	}

	return nil
}
