package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"forum/config"
	"forum/internal/entity"
)

type Authorization interface {
	CreateUser(user entity.UserModel) error
	GetUserByUsername(username string) (entity.UserModel, error)
	GetUserByEmail(email string) (entity.UserModel, error)
	SaveToken(user entity.UserModel) error
	GetUserByToken(token string) (entity.UserModel, error)
	DeleteToken(token string) error
}

type AuthRepository struct {
	db     *sql.DB
	config *config.Config
}

// NewUser lower layer
// that implements interface UserRepo
func NewAuthRepostiry(db *sql.DB, config *config.Config) *AuthRepository {
	return &AuthRepository{
		db:     db,
		config: config,
	}
}

// CreateUser errors for repo
// goes directly to database and create new user in DB
func (u *AuthRepository) CreateUser(user entity.UserModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(u.config.CtxTimeout)*time.Second)
	defer cancel()

	query := `INSERT INTO user (username, password, email) VALUES ($1, $2, $3);`
	_, err := u.db.ExecContext(ctx, query, user.Username, user.Password, user.Email)
	if err != nil {
		return fmt.Errorf("repository: create user :%w", err)
	}

	return nil
}

// get user info
func (u *AuthRepository) GetUserByUsername(username string) (entity.UserModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(u.config.CtxTimeout)*time.Second)
	defer cancel()

	var user entity.UserModel
	selectQuery := `SELECT email, userId, username, password, creationDate FROM user WHERE username = $1;`
	if err := u.db.QueryRowContext(ctx, selectQuery, username).Scan(&user.Email, &user.UserId, &user.Username, &user.Password, &user.CreatedAt); err != nil {
		return entity.UserModel{}, fmt.Errorf("repository: get user by username: %w", err)
	}

	return user, nil
}

func (u *AuthRepository) GetUserByEmail(email string) (entity.UserModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(u.config.CtxTimeout)*time.Second)
	defer cancel()

	var user entity.UserModel
	selectQuery := `SELECT email, userId, username, password, creationDate FROM user WHERE email = $1;`
	if err := u.db.QueryRowContext(ctx, selectQuery, email).Scan(&user.Email, &user.UserId, &user.Username, &user.Password, &user.CreatedAt); err != nil {
		return entity.UserModel{}, fmt.Errorf("repository: get user :%w", err)
	}

	return user, nil
}

func (u *AuthRepository) SaveToken(user entity.UserModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(u.config.CtxTimeout)*time.Second)
	defer cancel()

	query := `UPDATE user SET token = $1, expiresAt = $2 WHERE username = $3;`
	_, err := u.db.ExecContext(ctx, query, user.Token, user.ExpirationTime, user.Username)
	if err != nil {
		fmt.Printf("error SaveToken in repo -%v", err)

		return fmt.Errorf("repository: user: save token :%w", err)
	}

	return nil
}

// reeives user info by token
func (u *AuthRepository) GetUserByToken(token string) (entity.UserModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(u.config.CtxTimeout)*time.Second)
	defer cancel()

	var user entity.UserModel
	query := `SELECT userId, username, email, password, creationDate, token, expiresAt FROM user WHERE token = $1;`
	if err := u.db.QueryRowContext(ctx, query, token).Scan(&user.UserId, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.Token, &user.ExpirationTime); err != nil {
		return entity.UserModel{}, fmt.Errorf("repository: get user by token: %w", err)
	}

	return user, nil
}

// deletes token
func (u *AuthRepository) DeleteToken(token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(u.config.CtxTimeout)*time.Second)
	defer cancel()

	query := `UPDATE user SET token = NULL, expiresAt = NULL WHERE token = $1;`
	_, err := u.db.ExecContext(ctx, query, token)
	if err != nil {
		return fmt.Errorf("repository: delete token: %w", err)
	}

	return nil
}

// updates user password
// func (u *AuthRepository) UpdateUser(username, password string) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(u.config.CtxTimeout)*time.Second)
// 	defer cancel()

// 	updateQuery := `UPDATE user SET password = $1 WHERE username = $2;`
// 	_, err := u.db.ExecContext(ctx, updateQuery, password, username)
// 	if err != nil {
// 		return fmt.Errorf("repository: update user:%w", err)
// 	}

// 	return nil
// }

// // deletes user
// func (u *AuthRepository) DeleteUser(username string) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(u.config.CtxTimeout)*time.Second)
// 	defer cancel()

// 	deleteQuery := `DELETE FROM user WHERE username = $1;`
// 	_, err := u.db.ExecContext(ctx, deleteQuery, username)
// 	if err != nil {
// 		return fmt.Errorf("repository: delete user :%w", err)
// 	}

// 	return nil
// }
