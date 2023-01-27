package usecase

import "forum/internal/repository"

type UseCase struct {
	AuthorizationUsecase `json:"authorization_usecase,omitempty"`
	PostsUsecase         `json:"posts_usecase,omitempty"`
	PostsVoterUsecase    `json:"posts_voter_usecase,omitempty"`
	CommentsUsecase      `json:"comments_usecase,omitempty"`
	UsersUsecase         `json:"users_usecase,omitempty"`
}

func NewUseCase(r *repository.Repository) *UseCase {
	return &UseCase{
		AuthorizationUsecase: NewAuthUseCase(r.Authorization),
		PostsUsecase:         NewPostUseCase(r.Posts),
		PostsVoterUsecase:    NewPostVotesUsecase(r.PostVoter),
		CommentsUsecase:      NewCommentUsecase(r.Commenter),
		UsersUsecase:         NewUserUsecase(r.User),
	}
}
