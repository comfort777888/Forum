package entity

import "time"

type UserModel struct {
	UserId          int
	Email           string
	Username        string
	Password        string
	ConfirmPassword string
	Posts           int
	CreatedAt       time.Time

	Token          string
	ExpirationTime time.Time
	Emailcheck     string
	Usernamecheck  string
	// Passwordcheck        string
	// ConfirmPasswordcheck string
}
