package repository

import (
	"database/sql"
	"forum/config"
)

type Repository struct {
	Authorization
	Posts
	PostVoter
	Commenter
	User
}

func NewRepository(db *sql.DB, cnf *config.Config) *Repository {
	return &Repository{
		Authorization: NewAuthRepostiry(db, cnf),
		Posts:         NewPostRepository(db, cnf),
		PostVoter:     NewPostVotingRepostiry(db, cnf),
		Commenter:     NewCommentsRepostiry(db, cnf),
		User:          NewUserRepository(db, cnf),
	}
}
