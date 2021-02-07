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

func (ur *UserRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	if err := ur.db.QueryRow(`
		SELECT id,username,email,role,password,created_date FROM user WHERE email = ?
	`, email).Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.Password, &user.CreatedDate); err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRepository) GetByUsername(username string) (*models.User, error) {
	user := &models.User{}
	if err := ur.db.QueryRow(`
		SELECT id,username,email,role,password,created_date FROM user WHERE username = ?
	`, username).Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.Password, &user.CreatedDate); err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRepository) GetByID(id int) (*models.User, error) {
	user := &models.User{}
	if err := ur.db.QueryRow(`
		SELECT id,username,email,role,password,created_date FROM user WHERE id = ?
	`, id).Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.Password, &user.CreatedDate); err != nil {
		return nil, err
	}
	return user, nil
}
