package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"forum/config"
	"forum/internal/entity"
)

type Posts interface {
	CreatePost(post entity.Post) (int, error)
	GetAllPosts() ([]entity.Post, error)
	GetPostbyId(id int) (entity.Post, error)
	GetPostsByCategory(category string) ([]entity.Post, error)
	CategoriesByPostId(id int) ([]string, error)
	GetCreatedPosts(author string) ([]entity.Post, error)
	UpdatePostById(post entity.Post) error
	GetNewstPosts() ([]entity.Post, error)
	GetOldesPosts() ([]entity.Post, error)
	GetMostLikedPosts() ([]entity.Post, error)
	GetMostDisikedPosts() ([]entity.Post, error)
}

type PostRepository struct {
	db  *sql.DB
	cnf *config.Config
}

func NewPostRepository(db *sql.DB, cnf *config.Config) *PostRepository {
	return &PostRepository{
		db,
		cnf,
	}
}

func (p *PostRepository) CreatePost(post entity.Post) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.cnf.CtxTimeout)*time.Second)
	defer cancel()
	// begin transaction
	tx, err := p.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("repository: create post: transaction %w", err)
	}
	query := `INSERT INTO posts (author, title, content) VALUES ($1, $2, $3) RETURNING postId;`
	var id int
	if err := tx.QueryRowContext(ctx, query, post.PostAuthor, post.Title, post.Content).Scan(&id); err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("repository: create post: insert query %w", err)
	}

	query = `UPDATE user SET posts = posts + 1 WHERE username = $1;`
	_, err = tx.ExecContext(ctx, query, post.PostAuthor)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("repository: create post: update post count %w", err)
	}

	query = `INSERT INTO posts_category (postCategoryId, category) VALUES ($1,$2);`
	for _, category := range post.Category {
		_, err = tx.ExecContext(ctx, query, id, category)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("repository: create post: update category %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("repository: create post: commit transaction %w", err)
	}

	return id, nil
}

func (p *PostRepository) GetAllPosts() ([]entity.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.cnf.CtxTimeout)*time.Second)
	defer cancel()

	query := `SELECT postId, author, title, content, creationDate, likes, dislikes FROM posts;`
	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("repository: get all posts: %w", err)
	}

	defer rows.Close()

	var posts []entity.Post
	for rows.Next() {
		var post entity.Post
		if err := rows.Scan(&post.PostId, &post.PostAuthor, &post.Title, &post.Content, &post.CreationTime, &post.Likes, &post.Dislikes); err != nil {
			return nil, fmt.Errorf("repository: get all posts: scan: %w", err)
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: get all posts:  rows.err %w", err)
	}

	return posts, nil
}

func (p *PostRepository) GetPostbyId(id int) (entity.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.cnf.CtxTimeout)*time.Second)
	defer cancel()

	query := `SELECT postId, author, title, content, creationDate, likes, dislikes FROM posts WHERE postId = $1;`
	var post entity.Post
	if err := p.db.QueryRowContext(ctx, query, id).Scan(&post.PostId, &post.PostAuthor, &post.Title, &post.Content, &post.CreationTime, &post.Likes, &post.Dislikes); err != nil {
		return entity.Post{}, fmt.Errorf("repository: get postbyID: query %w", err)
	}

	return post, nil
}

func (p *PostRepository) GetPostsByCategory(category string) ([]entity.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.cnf.CtxTimeout)*time.Second)
	defer cancel()

	query := `SELECT postId, author, title, content, creationDate, likes, dislikes FROM posts WHERE postId IN (SELECT postCategoryId FROM posts_category WHERE category = $1);`
	rows, err := p.db.QueryContext(ctx, query, category)
	if err != nil {
		return nil, fmt.Errorf("repository: get postbyCategory: query %w", err)
	}

	defer rows.Close()

	var posts []entity.Post
	for rows.Next() {
		var post entity.Post

		if err := rows.Scan(&post.PostId, &post.PostAuthor, &post.Title, &post.Content, &post.CreationTime, &post.Likes, &post.Dislikes); err != nil {
			return nil, fmt.Errorf("repository: get postbyCategory: scan %w", err)
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: get postbyCategory: rows error %w", err)
	}

	return posts, nil
}

func (p *PostRepository) CategoriesByPostId(id int) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.cnf.CtxTimeout)*time.Second)
	defer cancel()

	query := `SELECT category FROM posts_category WHERE postCategoryId = $1;`

	rows, err := p.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("repository: categoriesbyId: query: %w", err)
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string

		if err := rows.Scan(&category); err != nil {
			return nil, fmt.Errorf("repository: categoriesbyId: scan: %w", err)
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: categoriesbyId: rows error %w", err)
	}

	return categories, nil
}

func (p *PostRepository) GetCreatedPosts(author string) ([]entity.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.cnf.CtxTimeout)*time.Second)
	defer cancel()

	query := `SELECT postId, author, title, context, creationDate, likes, dislikes FROM posts WHERE author = $1;`

	rows, err := p.db.QueryContext(ctx, query, author)
	if err != nil {
		return nil, fmt.Errorf("repository: get createdtposts: query %w", err)
	}

	defer rows.Close()

	var posts []entity.Post
	for rows.Next() {
		var post entity.Post

		if err := rows.Scan(&post.PostId, &post.PostAuthor, &post.Title, &post.Content, &post.CreationTime, &post.Likes, &post.Dislikes); err != nil {
			return nil, fmt.Errorf("repository: get createdtposts: scan %w", err)
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: get createdtposts: rows error %w", err)
	}

	return posts, nil
}

func (p *PostRepository) UpdatePostById(post entity.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.cnf.CtxTimeout)*time.Second)
	defer cancel()

	query := `UPDATE posts SET title = $1, content = $2 WHERE postId = $3;`
	_, err := p.db.ExecContext(ctx, query, post.Title, post.Content, post.PostId)
	if err != nil {
		return fmt.Errorf("repository: update by id: query %w", err)
	}

	return nil
}

func (p *PostRepository) GetNewstPosts() ([]entity.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.cnf.CtxTimeout)*time.Second)
	defer cancel()

	query := `SELECT * FROM posts ORDER BY creationDate DESC;`
	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []entity.Post
	for rows.Next() {
		var post entity.Post

		if err := rows.Scan(&post.PostId, &post.PostAuthor, &post.Title, &post.Content, &post.CreationTime, &post.Likes, &post.Dislikes); err != nil {
			return nil, fmt.Errorf("repository: get Newsttposts: scan %w", err)
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: get Newsttposts: rows error %w", err)
	}

	return posts, nil
}

func (p *PostRepository) GetOldesPosts() ([]entity.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.cnf.CtxTimeout)*time.Second)
	defer cancel()

	query := `SELECT postId, author, title, content, creationDate, likes, dislikes FROM posts ORDER BY creationDate;`
	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("repository: get newstpost: query %w", err)
	}

	defer rows.Close()

	var posts []entity.Post
	for rows.Next() {
		var post entity.Post

		if err := rows.Scan(&post.PostId, &post.PostAuthor, &post.Title, &post.Content, &post.CreationTime, &post.Likes, &post.Dislikes); err != nil {
			return nil, fmt.Errorf("repository: get newstpost: scan %w", err)
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: get newstpost: rows error %w", err)
	}

	return posts, nil
}

func (p *PostRepository) GetMostLikedPosts() ([]entity.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.cnf.CtxTimeout)*time.Second)
	defer cancel()

	query := `SELECT postId, author, title, context, creationDate, likes, dislikes FROM posts ORDER BY likes DESC; `
	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []entity.Post
	for rows.Next() {
		var post entity.Post

		if err := rows.Scan(&post.PostId, &post.PostAuthor, &post.Title, &post.Content, &post.CreationTime, &post.Likes, &post.Dislikes); err != nil {
			return nil, fmt.Errorf("repository: get newstpost: scan %w", err)
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: get newstpost: rows error %w", err)
	}

	return posts, nil
}

func (p *PostRepository) GetMostDisikedPosts() ([]entity.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.cnf.CtxTimeout)*time.Second)
	defer cancel()

	query := `SELECT postId, author, title, context, creationDate, likes, dislikes FROM posts ORDER BY dislikes DESC; `
	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []entity.Post
	for rows.Next() {
		var post entity.Post

		if err := rows.Scan(&post.PostId, &post.PostAuthor, &post.Title, &post.Content, &post.CreationTime, &post.Likes, &post.Dislikes); err != nil {
			return nil, fmt.Errorf("repository: get newstpost: scan %w", err)
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: get newstpost: rows error %w", err)
	}

	return posts, nil
}
