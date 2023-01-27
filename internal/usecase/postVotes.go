package usecase

import (
	"database/sql"
	"errors"
	"fmt"

	"forum/internal/repository"
)

type PostsVoterUsecase interface {
	GetPostDislikes(postId int) ([]string, error)
	GetPostLikes(postId int) ([]string, error)
	DisikePost(postId int, username string) error
	LikePost(postId int, username string) error
}

type PostVoterUsecase struct {
	PostVotesRepository repository.PostVoter
}

func NewPostVotesUsecase(p repository.PostVoter) *PostVoterUsecase {
	return &PostVoterUsecase{
		p,
	}
}

func (p *PostVoterUsecase) LikePost(postId int, username string) error {
	if err := p.PostVotesRepository.LikePostByUser(postId, username); err == nil {
		if err := p.PostVotesRepository.RemoveLikePost(postId, username); err != nil {
			return err
		}
		return nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("usecase:like post: %w", err)
	}

	if err := p.PostVotesRepository.DislikePostByUser(postId, username); err == nil {
		if err := p.PostVotesRepository.RemoveDislikePost(postId, username); err != nil {
			return err
		}
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("usecase:dislike post: %w", err)
	}

	if err := p.PostVotesRepository.LikePost(postId, username); err != nil {
		return err
	}

	return nil
}

func (p *PostVoterUsecase) DisikePost(postId int, username string) error {
	if err := p.PostVotesRepository.DislikePostByUser(postId, username); err == nil {

		if err := p.PostVotesRepository.RemoveDislikePost(postId, username); err != nil {
			return err
		}
		return nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("usecase: like post: %w", err)
	}

	if err := p.PostVotesRepository.LikePostByUser(postId, username); err == nil {
		if err := p.PostVotesRepository.RemoveLikePost(postId, username); err != nil {
			return err
		}
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("usecase: dislike post: %w", err)
	}

	if err := p.PostVotesRepository.DislikePost(postId, username); err != nil {
		return err
	}

	return nil
}

func (p *PostVoterUsecase) GetPostLikes(postId int) ([]string, error) {
	users, err := p.PostVotesRepository.GetPostLikes(postId)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (p *PostVoterUsecase) GetPostDislikes(postId int) ([]string, error) {
	users, err := p.PostVotesRepository.GetPostDislikes(postId)
	if err != nil {
		return nil, err
	}

	return users, nil
}
