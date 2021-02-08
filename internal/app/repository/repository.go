package repository

import (
	"database/sql"

	"github.com/Akezhan1/forum/internal/app/models"
)

type User interface {
	Create(*models.User) (int64, error)
	GetByEmail(string) (*models.User, error)
	GetByUsername(string) (*models.User, error)
	GetByID(int) (*models.User, error)

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
