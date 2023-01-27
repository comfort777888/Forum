package usecase

import (
	"errors"
	"strings"

	"forum/internal/entity"
	"forum/internal/repository"
)

type UsersUsecase interface {
	GetUserByName(username string) (entity.UserModel, error)
	GetPostsByName(username string, query map[string][]string) ([]entity.Post, error)
}

type UserUsecase struct {
	ur repository.User
}

func NewUserUsecase(r repository.User) *UserUsecase {
	return &UserUsecase{
		r,
	}
}

var ErrInvalidQuery = errors.New("invalid query")

func (u *UserUsecase) GetPostsByName(username string, query map[string][]string) ([]entity.Post, error) {
	var posts []entity.Post
	var err error

	search, ok := query["posts"]
	if !ok {
		return nil, ErrInvalidQuery
	}

	switch strings.Join(search, "") {
	case "created":
		posts, err = u.ur.GetPostsByName(username)
	case "liked":
		posts, err = u.ur.GetLikedPostsByName(username)
	case "disliked":
		posts, err = u.ur.GetDislikedPostsByName(username)
	case "commented":
		posts, err = u.ur.GetCommentedPostsByName(username)
	default:
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	for i := range posts {
		category, err := u.ur.GetAllCategoriesByPostId(posts[i].PostId)
		if err != nil {
			return nil, err
		}
		posts[i].Category = category
	}
	return posts, nil
}

func (u *UserUsecase) GetUserByName(username string) (entity.UserModel, error) {
	return u.ur.GetUser(username)
}
