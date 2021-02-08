package repository

import (
	"database/sql"

	"github.com/Akezhan1/forum/internal/app/models"
)

type User interface {
	CreateUser(*models.User) (int64, error)
	GetUserByEmail(string) (*models.User, error)
	GetUserByUsername(string) (*models.User, error)
	GetUserByID(int) (*models.User, error)

	CreateSession(*models.Session) error
	UpdateSession(*models.Session) error
	DeleteSession(string) error
	GetSession(string) (*models.Session, error)
}

type Post interface {
}

type Comment interface {
}

type Repository struct {
	User
	Post
	Comment
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		User: NewUserRepository(db),
	}
}
