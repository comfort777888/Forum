package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"forum/config"
	"forum/internal/entity"
)

type Commenter interface {
	CreateComment(comment entity.Comments) error
	GetCommentById(id int) (entity.Comments, error)
	GetCommentsByPostId(id int) ([]entity.Comments, error)
	LikeComment(commentId int, username string) error
	DislikeComment(commentId int, username string) error
	RemoveLikeFromComment(commentId int, username string) error
	RemoveDislikeFromComment(commentId int, username string) error
	CommentLiked(commentId int, username string) error
	CommentDisliked(commentId int, username string) error
	GetCommentDislikes(postId int) (map[int][]string, error)
	GetCommentLikes(postId int) (map[int][]string, error)
}

type CommentsRepository struct {
	db  *sql.DB
	cnf *config.Config
}

func NewCommentsRepostiry(db *sql.DB, config *config.Config) *CommentsRepository {
	return &CommentsRepository{
		db,
		config,
	}
}

func (c *CommentsRepository) CreateComment(comment entity.Comments) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.cnf.CtxTimeout)*time.Second)
	defer cancel()

	query := `INSERT INTO comments (postId, author, content) VALUES ($1, $2, $3);`
	_, err := c.db.ExecContext(ctx, query, comment.PostId, comment.Author, comment.Content)
	if err != nil {
		return fmt.Errorf("repository: create comment: %w", err)
	}

	return nil
}

func (c *CommentsRepository) GetCommentById(id int) (entity.Comments, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.cnf.CtxTimeout)*time.Second)
	defer cancel()

	query := `SELECT commentsId, postId, author, content FROM comments WHERE commentsId = $1;`
	var comment entity.Comments
	if err := c.db.QueryRowContext(ctx, query, id).Scan(&comment.CommentId, &comment.PostId, &comment.Author, &comment.Content); err != nil {
		return entity.Comments{}, fmt.Errorf("repository: get comment by id: %w", err)
	}

	return comment, nil
}

func (c *CommentsRepository) GetCommentsByPostId(id int) ([]entity.Comments, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.cnf.CtxTimeout)*time.Second)
	defer cancel()

	query := `SELECT commentsId, postId, author, content, likes, dislikes FROM comments WHERE postId = $1;`
	rows, err := c.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("repository:comments:createcomment %w", err)
	}

	defer rows.Close()

	var comments []entity.Comments
	for rows.Next() {
		var comment entity.Comments
		if err := rows.Scan(&comment.CommentId, &comment.PostId, &comment.Author, &comment.Content, &comment.Likes, &comment.Dislikes); err != nil {
			return nil, fmt.Errorf("repository:comments:createcomment:scan %w", err)
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository:comments:getcommentbyPostId: rows error %w", err)
	}

	return comments, nil
}

func (c *CommentsRepository) LikeComment(commentId int, username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.cnf.CtxTimeout)*time.Second)
	defer cancel()

	query := `INSERT INTO likes (commentsId, username) VALUES ($1,$2);`
	_, err := c.db.ExecContext(ctx, query, commentId, username)
	if err != nil {
		return fmt.Errorf("repository:comments:likecomment:query1 %w", err)
	}

	query = `UPDATE comments SET likes = likes +1 WHERE commentsId = $1;`
	_, err = c.db.ExecContext(ctx, query, commentId)
	if err != nil {
		return fmt.Errorf("repository:comments:likecomment:query2 %w", err)
	}

	return nil
}

func (c *CommentsRepository) DislikeComment(commentId int, username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.cnf.CtxTimeout)*time.Second)
	defer cancel()

	query := `INSERT INTO dislikes (commentsId, username) VALUES ($1,$2);`
	_, err := c.db.ExecContext(ctx, query, commentId, username)
	if err != nil {
		return fmt.Errorf("repository:comments:dislikecomment:query1 %w", err)
	}

	query = `UPDATE comments SET dislikes = dislikes +1 WHERE commentsId = $1;`
	_, err = c.db.ExecContext(ctx, query, commentId)
	if err != nil {
		return fmt.Errorf("repository:comments:dislikecomment:query2 %w", err)
	}

	return nil
}

func (c *CommentsRepository) RemoveLikeFromComment(commentId int, username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.cnf.CtxTimeout)*time.Second)
	defer cancel()

	query := `DELETE FROM likes WHERE commentsId = $1 AND username =$2;`
	_, err := c.db.ExecContext(ctx, query, commentId, username)
	if err != nil {
		return fmt.Errorf("repository:comments:removelikefromcomment:query1 %w", err)
	}

	query = `UPDATE comments SET likes = likes -1 WHERE commentsId = $1;`
	_, err = c.db.ExecContext(ctx, query, commentId)
	if err != nil {
		return fmt.Errorf("repository:comments:removelikefromcomment:query2 %w", err)
	}

	return nil
}

func (c *CommentsRepository) RemoveDislikeFromComment(commentId int, username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.cnf.CtxTimeout)*time.Second)
	defer cancel()

	query := `DELETE FROM dislikes WHERE commentsId =$1 AND username =$2;`
	_, err := c.db.ExecContext(ctx, query, commentId, username)
	if err != nil {
		return fmt.Errorf("repository:comments:removedislikefromcomment:query1 %w", err)
	}

	query = `UPDATE comments SET dislikes = dislikes -1 WHERE commentsId = $1;`
	_, err = c.db.ExecContext(ctx, query, commentId)
	if err != nil {
		return fmt.Errorf("repository:comments:removedislikefromcomment:query2 %w", err)
	}

	return nil
}

func (c *CommentsRepository) CommentLiked(commentId int, username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.cnf.CtxTimeout)*time.Second)
	defer cancel()

	query := `SELECT username FROM likes WHERE commentsId = $1 AND username =$2;`
	var user string
	if err := c.db.QueryRowContext(ctx, query, commentId, username).Scan(&user); err != nil {
		return fmt.Errorf("repository:comments:commentliked:query1 %w", err)
	}

	return nil
}

func (c *CommentsRepository) CommentDisliked(commentId int, username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.cnf.CtxTimeout)*time.Second)
	defer cancel()

	query := `SELECT username FROM dislikes WHERE commentsId = $1 AND username = $2;`
	var user string
	if err := c.db.QueryRowContext(ctx, query, commentId, username).Scan(&user); err != nil {
		return fmt.Errorf("repository:comments:commentdisliked:query1 %w", err)
	}

	return nil
}

func (c *CommentsRepository) GetCommentLikes(postId int) (map[int][]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.cnf.CtxTimeout)*time.Second)
	defer cancel()

	users := make(map[int][]string)
	queryCommentId := `SELECT commentsId  FROM comments WHERE postId = $1;`
	queryUsers := `SELECT username FROM likes WHERE commentsId = $1;`

	rowsComments, err := c.db.QueryContext(ctx, queryCommentId, postId)
	if err != nil {
		return nil, fmt.Errorf("repository:comments:commentlikes:query1 %w", err)
	}
	defer rowsComments.Close()

	for rowsComments.Next() {
		var id int

		if err := rowsComments.Scan(&id); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, fmt.Errorf("repository:comments:commentlikes:scancomment %w", err)
		}

		var usernames []string

		rowsUsers, err := c.db.QueryContext(ctx, queryUsers, postId)
		if err != nil {
			return nil, fmt.Errorf("repository:comments:commentlikes:query2 %w", err)
		}

		defer rowsUsers.Close()

		for rowsUsers.Next() {
			var user string

			if err := rowsUsers.Scan(&user); err != nil {
				return nil, fmt.Errorf("repository:comments:commentlikes:scanuser %w", err)
			}

			usernames = append(usernames, user)
		}

		if err := rowsUsers.Err(); err != nil {
			return nil, fmt.Errorf("repository:comments:commentlikes: rows error %w", err)
		}

		users[id] = usernames
	}

	if err := rowsComments.Err(); err != nil {
		return nil, fmt.Errorf("repository:comments:commentlikes: rows error %w", err)
	}

	return users, nil
}

func (c *CommentsRepository) GetCommentDislikes(postId int) (map[int][]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.cnf.CtxTimeout)*time.Second)
	defer cancel()
	users := make(map[int][]string)

	queryCommentId := `SELECT commentsId  FROM comments WHERE postId = $1;`
	queryUsers := `SELECT username FROM dislikes WHERE commentsId = $1;`

	rowsComments, err := c.db.QueryContext(ctx, queryCommentId, postId)
	if err != nil {
		return nil, fmt.Errorf("repository:comments:commentdislikes:query1 %w", err)
	}
	defer rowsComments.Close()

	for rowsComments.Next() {
		var id int

		if err := rowsComments.Scan(&id); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, fmt.Errorf("repository:comments:commentdislikes:scancomment %w", err)
		}

		var usernames []string
		rowsUsers, err := c.db.QueryContext(ctx, queryUsers, postId)
		if err != nil {
			return nil, fmt.Errorf("repository:comments:commentdislikes:query2 %w", err)
		}

		for rowsUsers.Next() {
			var user string
			if err := rowsUsers.Scan(&user); err != nil {
				return nil, fmt.Errorf("repository:comments:commentdislikes:scanuser %w", err)
			}

			usernames = append(usernames, user)
		}

		users[id] = usernames
	}
	if err := rowsComments.Err(); err != nil {
		return nil, fmt.Errorf("repository:comments:commentdislikes: rows error %w", err)
	}

	return users, nil
}
