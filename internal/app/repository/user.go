package repository

import (
	"database/sql"

	"github.com/Akezhan1/forum/internal/app/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

func (ur *UserRepository) Create(user *models.User) (int64, error) {
	result, err := ur.db.Exec(`
	INSERT INTO user (email,username,password,role,created_date) 
	VALUES (?,?,?,?,?)`, user.Email, user.Username, user.Password, user.Role, user.CreatedDate)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}
