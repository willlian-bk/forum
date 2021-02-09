package service

import (
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/Akezhan1/forum/internal/app/models"
	"github.com/Akezhan1/forum/internal/app/repository"
)

type PostService struct {
	repo repository.Post
}

func NewPostService(repo repository.Post) *PostService {
	return &PostService{repo}
}

func (ps *PostService) Create(post *models.Post) (int, int, error) {
	if err := ps.validateParams(post); err != nil {
		return 400, -1, err
	}

	post.CreatedDate = time.Now()
	post.UpdatedDate = post.CreatedDate
	post.Likes = 0
	post.Dislikes = 0

	id, err := ps.repo.Create(post)
	if err != nil {
		return 500, -1, err
	}

	return 200, int(id), nil
}

func (ps *PostService) Get(id int) (*models.Post, error) {
	if id < 0 {
		return nil, errors.New("Invalid Id")
	}
	post, err := ps.repo.GetPostByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Invalid Id")
		}
		return nil, err
	}

	postCategories, err := ps.repo.GetPostsCategories(id)
	if err != nil {
		return nil, err
	}

	post.Categories = postCategories

	return post, nil
}

func (ps *PostService) GetValidCategories() ([]string, error) {
	return ps.repo.GetValidCategories()
}

func (ps *PostService) EstimatePost(postID, userID, types string) error {
	if types != "like" && types != "dislike" {
		return errors.New("Invalid Type")
	}

	postIDint, err := strconv.Atoi(postID)
	if err != nil {
		return err
	}

	userIDint, err := strconv.Atoi(userID)
	if err != nil {
		return err
	}

	post := &models.Post{
		ID:     postIDint,
		UserID: userIDint,
	}

	return ps.repo.EstimatePost(post, types)
}

func (ps *PostService) validateParams(post *models.Post) error {
	if post.Title == "" {
		return errors.New("Invalid Title")
	}

	if post.Content == "" {
		return errors.New("Invalid Content")
	}

	if categories, err := ps.repo.GetValidCategories(); err != nil {
		return err
	} else {
		for _, postCategory := range post.Categories {
			ok := 1
			for _, category := range categories {
				if postCategory == category {
					ok--
				}
			}
			if ok != 0 {
				return errors.New("Invalid Category")
			}
		}
	}

	return nil
}
