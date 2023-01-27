package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"forum/config"
	"forum/internal/entity"
)

type User interface {
	GetPostsByName(username string) ([]entity.Post, error)
	GetLikedPostsByName(username string) ([]entity.Post, error)
	GetDislikedPostsByName(username string) ([]entity.Post, error)
	GetCommentedPostsByName(username string) ([]entity.Post, error)
	GetAllCategoriesByPostId(postId int) ([]string, error)
	GetUser(username string) (entity.UserModel, error)
}

type UserRepository struct {
	db  *sql.DB
	cnf *config.Config
}

func NewUserRepository(db *sql.DB, cnf *config.Config) *UserRepository {
	return &UserRepository{
		db,
		cnf,
	}
}

func (u *UserRepository) GetPostsByName(username string) ([]entity.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(u.cnf.CtxTimeout)*time.Second)
	defer cancel()

	var posts []entity.Post
	query := `SELECT postId , author, title, content, creationDate, likes, dislikes FROM posts WHERE author = $1;`
	rows, err := u.db.QueryContext(ctx, query, username)
	if err != nil {
		return nil, fmt.Errorf("repository:user:getpostbyname: query %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var post entity.Post

		if err := rows.Scan(&post.PostId, &post.PostAuthor, &post.Title, &post.Content, &post.CreationTime, &post.Likes, &post.Dislikes); err != nil {
			return nil, fmt.Errorf("repository:user:getpostbyname: scan %w", err)
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository:user:getpostbyname: rows error %w", err)
	}

	return posts, nil
}

func (u *UserRepository) GetLikedPostsByName(username string) ([]entity.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(u.cnf.CtxTimeout)*time.Second)
	defer cancel()

	var posts []entity.Post
	query := `SELECT postId , author, title, content, creationDate, likes, dislikes FROM posts WHERE postId IN (SELECT postID FROM likes WHERE username = $1);`
	rows, err := u.db.QueryContext(ctx, query, username)
	if err != nil {
		return nil, fmt.Errorf("repository:user:getlikedpostbyname: query %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var post entity.Post

		if err := rows.Scan(&post.PostId, &post.PostAuthor, &post.Title, &post.Content, &post.CreationTime, &post.Likes, &post.Dislikes); err != nil {
			return nil, fmt.Errorf("repository:user:getlikedpostbyname: scan %w", err)
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository:user:getlikedpostbyname: rows error %w", err)
	}

	return posts, nil
}

func (u *UserRepository) GetDislikedPostsByName(username string) ([]entity.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(u.cnf.CtxTimeout)*time.Second)
	defer cancel()

	var posts []entity.Post
	query := `SELECT postId , author, title, content, creationDate, likes, dislikes FROM posts WHERE postId IN (SELECT postID FROM dislikes WHERE username = $1);`
	rows, err := u.db.QueryContext(ctx, query, username)
	if err != nil {
		return nil, fmt.Errorf("repository:user:getdislikedpostbyname: query %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var post entity.Post

		if err := rows.Scan(&post.PostId, &post.PostAuthor, &post.Title, &post.Content, &post.CreationTime, &post.Likes, &post.Dislikes); err != nil {
			return nil, fmt.Errorf("repository:user:getdislikedpostbyname: scan %w", err)
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository:user:getdislikedpostbyname: rows error %w", err)
	}

	return posts, nil
}

func (u *UserRepository) GetCommentedPostsByName(username string) ([]entity.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(u.cnf.CtxTimeout)*time.Second)
	defer cancel()

	var posts []entity.Post
	query := `SELECT postId , author, title, content, creationDate, likes, dislikes FROM posts WHERE postId IN (SELECT postID FROM comments WHERE author  = $1);`
	rows, err := u.db.QueryContext(ctx, query, username)
	if err != nil {
		return nil, fmt.Errorf("repository:user:getdislikedpostbyname: query %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var post entity.Post

		if err := rows.Scan(&post.PostId, &post.PostAuthor, &post.Title, &post.Content, &post.CreationTime, &post.Likes, &post.Dislikes); err != nil {
			return nil, fmt.Errorf("repository:user:getdislikedpostbyname: scan %w", err)
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository:user:getdislikedpostbyname: rows error %w", err)
	}

	return posts, nil
}

func (u *UserRepository) GetAllCategoriesByPostId(postId int) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(u.cnf.CtxTimeout)*time.Second)
	defer cancel()

	query := `SELECT category FROM posts_category WHERE postCategoryId = $1;`
	rows, err := u.db.QueryContext(ctx, query, postId)
	if err != nil {
		return nil, fmt.Errorf("repository:user:getAllCategories: query %w", err)
	}

	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string

		if err := rows.Scan(&category); err != nil {
			return nil, fmt.Errorf("repository:user:getAllCategories: scan %w", err)
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository:user:getAllCategories: rows error %w", err)
	}

	return categories, nil
}

func (u *UserRepository) GetUser(username string) (entity.UserModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(u.cnf.CtxTimeout)*time.Second)
	defer cancel()

	var user entity.UserModel
	query := `SELECT userId, username, email, posts FROM user WHERE username =$1;`

	if err := u.db.QueryRowContext(ctx, query, username).Scan(&user.UserId, &user.Username, &user.Email, &user.Posts); err != nil {
		return entity.UserModel{}, fmt.Errorf("repository:user:getUser: scan %w", err)
	}

	return user, nil
}
