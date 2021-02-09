package repository

import (
	"database/sql"

	"github.com/Akezhan1/forum/internal/app/models"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db}
}

func (pr *PostRepository) Create(post *models.Post) (int64, error) {
	tx, err := pr.db.Begin()
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	result, err := tx.Exec(`
	INSERT INTO post (user_id,title,content,likes,dislikes,created_date,updated_date) 
	VALUES (?,?,?,?,?,?,?)`, post.UserID, post.Title, post.Content, post.Likes, post.Dislikes, post.CreatedDate, post.UpdatedDate)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	for _, category := range post.Categories {
		_, err = tx.Exec(`INSERT INTO category_posts (category_name, post_id) VALUES (?,?)`, category, id)
		if err != nil {
			tx.Rollback()
			return -1, err
		}
	}

	if err = tx.Commit(); err != nil {
		return -1, err
	}

	return id, nil
}

func (pr *PostRepository) GetPostByID(id int) (*models.Post, error) {
	post := &models.Post{}
	if err := pr.db.QueryRow(`
		SELECT id, user_id,title,content,likes,dislikes,created_date,updated_date FROM post WHERE id = ?
	`, id).Scan(&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.Likes,
		&post.Dislikes,
		&post.CreatedDate,
		&post.UpdatedDate); err != nil {
		return nil, err
	}

	rows, err := pr.db.Query("SELECT category_name FROM category_posts WHERE post_id = ?", post.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		name := ""
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		post.Categories = append(post.Categories, name)
	}

	return post, nil
}

func (pr *PostRepository) GetPostsCategories(id int) ([]string, error) {
	return nil, nil
}

func (pr *PostRepository) EstimatePost(post *models.Post, types string) error {
	return nil
}

func (pr *PostRepository) GetValidCategories() ([]string, error) {
	categories := []string{}

	rows, err := pr.db.Query("SELECT name FROM category")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		name := ""
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		categories = append(categories, name)
	}

	return categories, nil
}
