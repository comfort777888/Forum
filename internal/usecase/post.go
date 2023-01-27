package usecase

import (
	"errors"
	"strings"

	"forum/internal/entity"
	"forum/internal/repository"
)

type PostsUsecase interface {
	CreatePost(post entity.Post) error
	GetAllPosts() ([]entity.Post, error)
	GetPostsByCategory(category string) ([]entity.Post, error)
	GetPostsById(id int) (entity.Post, error)
	GetCreatedPosts(author string) ([]entity.Post, error)
	GetAllPostsFromFilter(user entity.UserModel, query map[string][]string) ([]entity.Post, error)
}

type PostUseCase struct {
	PostRepository repository.Posts
}

func NewPostUseCase(p repository.Posts) *PostUseCase {
	return &PostUseCase{
		p,
	}
}

func (pu *PostUseCase) CreatePost(post entity.Post) error {
	if err := verificatePost(post); err != nil {
		return err
	}
	_, err := pu.PostRepository.CreatePost(post)
	if err != nil {
		return err
	}
	return nil
}

func (pu *PostUseCase) GetAllPosts() ([]entity.Post, error) {
	posts, err := pu.PostRepository.GetAllPosts()
	if err != nil {
		return nil, err
	}

	for i := range posts {
		category, err := pu.PostRepository.CategoriesByPostId(posts[i].PostId)
		if err != nil {
			return nil, err
		}
		posts[i].Category = category
	}

	return posts, nil
}

func (pu *PostUseCase) GetPostsByCategory(category string) ([]entity.Post, error) {
	posts, err := pu.PostRepository.GetPostsByCategory(category)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (pu *PostUseCase) GetPostsById(id int) (entity.Post, error) {
	post, err := pu.PostRepository.GetPostbyId(id)
	if err != nil {
		return entity.Post{}, err
	}

	post.Category, err = pu.PostRepository.CategoriesByPostId(post.PostId)
	if err != nil {
		return entity.Post{}, err
	}

	return post, nil
}

func (pu *PostUseCase) GetCreatedPosts(author string) ([]entity.Post, error) {
	posts, err := pu.PostRepository.GetCreatedPosts(author)
	if err != nil {
		return nil, err
	}

	for i := range posts {
		category, err := pu.PostRepository.CategoriesByPostId(posts[i].PostId)
		if err != nil {
			return nil, err
		}
		posts[i].Category = category
	}

	return posts, nil
}

func (pu *PostUseCase) GetAllPostsFromFilter(user entity.UserModel, query map[string][]string) ([]entity.Post, error) {
	var posts []entity.Post
	var err error

	for key, value := range query {
		switch key {
		case "category":
			posts, err = pu.PostRepository.GetPostsByCategory(strings.Join(value, ""))
			if err != nil {
				return nil, err
			}
		case "time":
			switch strings.Join(value, "") {
			case "new":
				posts, err = pu.PostRepository.GetNewstPosts()
				if err != nil {
					return nil, err
				}
			case "old":
				posts, err = pu.PostRepository.GetOldesPosts()
				if err != nil {
					return nil, err
				}
			default:
				return nil, err
			}
			if err != nil {
				return nil, err
			}
		case "vote":
			switch strings.Join(value, "") {
			case "like":
				posts, err = pu.PostRepository.GetMostLikedPosts()
				if err != nil {
					return nil, err
				}
			case "dislike":
				posts, err = pu.PostRepository.GetMostDisikedPosts()
				if err != nil {
					return nil, err
				}
			default:
				return nil, err
			}
			if err != nil {
				return nil, err
			}
		case "clean":
			switch strings.Join(value, "") {
			case "true":
				posts, err = pu.PostRepository.GetAllPosts()
				if err != nil {
					return nil, err
				}
			}
		default:
			return nil, errors.New("error posts from filters")
		}
		for i := range posts {
			categories, err := pu.PostRepository.CategoriesByPostId(posts[i].PostId)
			if err != nil {
				return nil, err
			}
			posts[i].Category = categories
		}
	}

	return posts, nil
}

var (
	ErrInvalidTitleLength   = errors.New("invalid length for title")
	ErrInvalidContentLength = errors.New("invalid length for content")
	ErrNoTitle              = errors.New("no title")
	ErrNoContent            = errors.New("empty content section")
	ErrInvalidContent       = errors.New("invalid content")
)

func verificatePost(post entity.Post) error {
	if len(post.Title) > 100 {
		return ErrInvalidTitleLength
	}
	if len(post.Content) > 2000 {
		return ErrInvalidContentLength
	}

	post.Content = strings.Trim(post.Content, " \n\r")
	if post.Content == "" {
		return ErrNoContent
	}

	for _, w := range post.Content {
		if (w < 32 || w > 126) && (w != 13 && w != 10) {
			return ErrInvalidContent
		}
	}

	post.Title = strings.Trim(post.Title, " \n\r")
	if post.Title == "" {
		return ErrNoTitle
	}

	for _, w := range post.Title {
		if (w < 32 || w > 126) && (w != 13 && w != 10) {
			return ErrInvalidContent
		}
	}

	return nil
}
