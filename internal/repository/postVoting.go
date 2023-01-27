package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"forum/config"
	"forum/internal/entity"
)

type PostVoter interface {
	LikePost(postId int, username string) error
	DislikePost(postId int, username string) error
	RemoveDislikePost(postId int, username string) error
	RemoveLikePost(postId int, username string) error
	LikePostByUser(postId int, username string) error
	DislikePostByUser(postId int, username string) error
	GetPostLikes(postId int) ([]string, error)
	GetPostDislikes(postId int) ([]string, error)
}

type PostVotingRepository struct {
	db  *sql.DB
	cnf *config.Config
}

func NewPostVotingRepostiry(db *sql.DB, config *config.Config) *PostVotingRepository {
	return &PostVotingRepository{
		db:  db,
		cnf: config,
	}
}

func (p *PostVotingRepository) LikePost(postId int, username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.cnf.CtxTimeout)*time.Second)
	defer cancel()

	tx, err := p.db.Begin()
	if err != nil {
		return fmt.Errorf("repository:postvoting:like: begin %w", err)
	}

	query := `INSERT INTO likes (username, postId) VALUES ($1,$2);`
	_, err = tx.ExecContext(ctx, query, username, postId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("repository:postvoting:like: query1 %w", err)
	}

	query = `UPDATE posts SET likes = likes + 1 WHERE postId = $1;`
	_, err = tx.ExecContext(ctx, query, postId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("repository:postvoting:like: update %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("repository:postvoting:like: commit %w", err)
	}

	return nil
}

func (p *PostVotingRepository) DislikePost(postId int, username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.cnf.CtxTimeout)*time.Second)
	defer cancel()

	tx, err := p.db.Begin()
	if err != nil {
		return fmt.Errorf("repository:postvoting:dislike: begin %w", err)
	}
	query := `INSERT INTO dislikes (username, postId) VALUES ($1, $2);`
	_, err = tx.ExecContext(ctx, query, username, postId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("repository:postvoting:dislike:query1 %w", err)
	}

	query = `UPDATE posts SET dislikes = dislikes + 1 WHERE postId = $1;`
	_, err = tx.ExecContext(ctx, query, postId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("repository:postvoting:dislike: update %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("repository:postvoting:dislike: commit %w", err)
	}

	return nil
}

func (p *PostVotingRepository) RemoveDislikePost(postId int, username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.cnf.CtxTimeout)*time.Second)
	defer cancel()

	tx, err := p.db.Begin()
	if err != nil {
		return fmt.Errorf("repository: postvoting:removedislike: begin %w", err)
	}

	query := `DELETE FROM dislikes WHERE username = $1 AND postId = $2;`
	_, err = tx.ExecContext(ctx, query, username, postId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("repository: postvoting:removedislike: delete %w", err)
	}

	query = `UPDATE posts SET dislikes = dislikes - 1 WHERE postId = $1;`
	_, err = tx.ExecContext(ctx, query, postId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("repository: postvoting:removedislike: update %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("repository: postvoting:removedislike: commit %w", err)
	}

	return nil
}

func (p *PostVotingRepository) RemoveLikePost(postId int, username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.cnf.CtxTimeout)*time.Second)
	defer cancel()

	tx, err := p.db.Begin()
	if err != nil {
		return fmt.Errorf("repository: postvoting:removelike: begin %w", err)
	}
	query := `DELETE FROM likes WHERE username = $1 AND postId = $2;`
	_, err = tx.ExecContext(ctx, query, username, postId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("repository: postvoting:removelike: delete %w", err)
	}

	query = `UPDATE posts SET likes = likes - 1 WHERE postId = $1;`
	_, err = tx.ExecContext(ctx, query, postId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("repository: postvoting:removelike: update%w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("repository: postvoting:removelike: commit %w", err)
	}

	return nil
}

// postliked
func (p *PostVotingRepository) LikePostByUser(postId int, username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.cnf.CtxTimeout)*time.Second)
	defer cancel()

	var user entity.UserModel
	query := `SELECT username FROM likes WHERE postId = $1 AND username = $2;`
	if err := p.db.QueryRowContext(ctx, query, postId, username).Scan(&user.Username); err != nil {
		return fmt.Errorf("repository: postvoting:likepostbyusert %w", err)
	}

	return nil
}

func (p *PostVotingRepository) DislikePostByUser(postId int, username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.cnf.CtxTimeout)*time.Second)
	defer cancel()

	var user string
	query := `SELECT username FROM dislikes WHERE postId = $1 AND username = $2;`
	if err := p.db.QueryRowContext(ctx, query, postId, username).Scan(&user); err != nil {
		return fmt.Errorf("repository: postvoting:likepostbyusert %w", err)
	}

	return nil
}

func (p *PostVotingRepository) GetPostLikes(postId int) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.cnf.CtxTimeout)*time.Second)
	defer cancel()

	query := `SELECT username FROM likes WHERE postId = $1;`
	rows, err := p.db.QueryContext(ctx, query, postId)
	if err != nil {
		return nil, fmt.Errorf("repository: postvoting:getpostlikes:query %w", err)
	}

	defer rows.Close()

	var likes []string
	for rows.Next() {
		var like string

		if err := rows.Scan(&like); err != nil {
			return nil, fmt.Errorf("repository: postvoting:getpostlikes:scan %w", err)
		}

		likes = append(likes, like)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: postvoting:getpostlikes: rows error %w", err)
	}

	return likes, nil
}

func (p *PostVotingRepository) GetPostDislikes(postId int) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.cnf.CtxTimeout)*time.Second)
	defer cancel()

	query := `SELECT username FROM dislikes WHERE postId = $1;`
	rows, err := p.db.QueryContext(ctx, query, postId)
	if err != nil {
		return nil, fmt.Errorf("repository: postvoting:getpostlikes:query %w", err)
	}

	defer rows.Close()

	var dislikes []string
	for rows.Next() {
		var dislike string

		if err := rows.Scan(&dislike); err != nil {
			return nil, fmt.Errorf("repository: postvoting:getpostlikes:scan %w", err)
		}

		dislikes = append(dislikes, dislike)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: postvoting:getpostlikes: rows error %w", err)
	}

	return dislikes, nil
}
