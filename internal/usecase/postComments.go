package usecase

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"forum/internal/entity"
	"forum/internal/repository"
)

type CommentsUsecase interface {
	LikeComment(commentId int, username string) error
	DislikeComment(commentId int, username string) error
	GetCommentsByPostId(commentId int) ([]entity.Comments, error)
	CreateComment(comment entity.Comments) error
	GetCommentById(commentId int) (entity.Comments, error)
	CommentsLikes(postId int) (map[int][]string, error)
	CommentsDislikes(postId int) (map[int][]string, error)
}

type CommentUsecase struct {
	CommentRepository repository.Commenter
}

func NewCommentUsecase(c repository.Commenter) *CommentUsecase {
	return &CommentUsecase{
		c,
	}
}

// LikeComment likes
func (c *CommentUsecase) LikeComment(commentId int, username string) error {
	if err := c.CommentRepository.CommentLiked(commentId, username); err == nil {
		if err := c.CommentRepository.RemoveLikeFromComment(commentId, username); err != nil {
			return err
		}

		return nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("usecase:likecomment: %w", err)
	}

	if err := c.CommentRepository.CommentDisliked(commentId, username); err == nil {
		if err := c.CommentRepository.RemoveDislikeFromComment(commentId, username); err != nil {
			return err
		}
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("usecase:likecomment: %w", err)
	}

	if err := c.CommentRepository.LikeComment(commentId, username); err != nil {
		return err
	}

	return nil
}

func (c *CommentUsecase) DislikeComment(commentId int, username string) error {
	if err := c.CommentRepository.CommentDisliked(commentId, username); err == nil {
		if err := c.CommentRepository.RemoveDislikeFromComment(commentId, username); err != nil {
			return err
		}
		return nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("usecase:dislikecomment: %w", err)
	}

	if err := c.CommentRepository.CommentLiked(commentId, username); err == nil {
		if err := c.CommentRepository.RemoveLikeFromComment(commentId, username); err != nil {
			return err
		}
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("usecase:likecomment: %w", err)
	}

	if err := c.CommentRepository.DislikeComment(commentId, username); err != nil {
		return err
	}
	return nil
}

func (c *CommentUsecase) CreateComment(comment entity.Comments) error {
	if err := checkComment(comment); err != nil {
		return err
	}

	if err := c.CommentRepository.CreateComment(comment); err != nil {
		return err
	}

	return nil
}

func (c *CommentUsecase) GetCommentById(commentId int) (entity.Comments, error) {
	comment, err := c.CommentRepository.GetCommentById(commentId)
	if err != nil {
		return entity.Comments{}, err
	}

	return comment, nil
}

func (c *CommentUsecase) GetCommentsByPostId(commentId int) ([]entity.Comments, error) {
	comments, err := c.CommentRepository.GetCommentsByPostId(commentId)
	if err != nil {
		return []entity.Comments{}, err
	}

	return comments, nil
}

func (c *CommentUsecase) CommentsLikes(postId int) (map[int][]string, error) {
	users, err := c.CommentRepository.GetCommentLikes(postId)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (c *CommentUsecase) CommentsDislikes(postId int) (map[int][]string, error) {
	users, err := c.CommentRepository.GetCommentDislikes(postId)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func checkComment(comment entity.Comments) error {
	if len(comment.Content) > 500 {
		return fmt.Errorf("usecase:create comment: %w", ErrInvalidContentLength)
	}
	comment.Content = strings.Trim(comment.Content, " \n\r")
	comment.Content = strings.Trim(comment.Content, "\t") // ???????????

	for _, w := range comment.Content {
		if (w < 32 && w > 126) && (w != 13 && w != 10) {
			return fmt.Errorf("usecase: create comment: %w", ErrInvalidCharacter)
		}
	}

	if comment.Content == "" {
		return fmt.Errorf("usecase: create comment: %w", ErrInvalidCharacter)
	}

	return nil
}
