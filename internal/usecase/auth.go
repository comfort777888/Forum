package usecase

import (
	"errors"
	"fmt"
	"log"
	"net/mail"
	"regexp"
	"time"
	"unicode"

	"forum/internal/entity"
	"forum/internal/repository"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidNameLength = errors.New("invalid username length")
	ErrUserExist         = errors.New("username is already exist")
	ErrEmailExist        = errors.New("email is already exist")
	ErrUserNotFound      = errors.New("user is not found")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrInvalidEmail      = errors.New("invalid email address")
	ErrInvalidCharacter  = errors.New("invalid character")
	ErrHashPassword      = errors.New("cannot hash password")
	ErrConfirmPassword   = errors.New("password not the same")
)

type AuthorizationUsecase interface {
	CreateUserandValidate(user entity.UserModel) error
	CreateToken(username, password string) (entity.UserModel, error)
	ParseToken(token string) (entity.UserModel, error)
	DeleteToken(token string) error
}

type AuthUserUse struct {
	Repository repository.Authorization
}

func NewAuthUseCase(ar repository.Authorization) *AuthUserUse {
	return &AuthUserUse{
		Repository: ar,
	}
}

// CreateUserandValidate creates user if not exist in db, checks user's info
// hashes password
func (u *AuthUserUse) CreateUserandValidate(user entity.UserModel) error {
	if err := checkUser(user); err != nil {
		return err
	}

	if _, err := u.Repository.GetUserByUsername(user.Username); err == nil {
		return fmt.Errorf("usecase: create and validate: %w", ErrUserExist)
	}

	if _, err := u.Repository.GetUserByEmail(user.Email); err == nil {
		return fmt.Errorf("usecase: create and validate: %w", ErrEmailExist)
	}

	var err error
	user.CreatedAt = time.Now()

	user.Password, err = generateHashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("usecase: cannot generate hash: %w", ErrHashPassword)
	}

	return u.Repository.CreateUser(user)
}

func (u *AuthUserUse) CreateToken(username, password string) (entity.UserModel, error) {
	user, err := u.Repository.GetUserByUsername(username)
	if err != nil {
		return entity.UserModel{}, fmt.Errorf("usecase: create token: %w", ErrUserNotFound)
	}

	if err := checkPasswordHash(password, user.Password); err != nil {
		return entity.UserModel{}, fmt.Errorf("usecase: check password hash: %w", ErrInvalidPassword)
	}

	token, err := uuid.NewV4()
	if err != nil {
		return entity.UserModel{}, fmt.Errorf("error create token uuid: %w", err)
	}
	user.Token = token.String()
	user.ExpirationTime = time.Now().Add(12 * time.Hour)

	if err := u.Repository.SaveToken(user); err != nil {
		fmt.Printf("error create token - SaveToken - %v", err)
		return entity.UserModel{}, fmt.Errorf("usercase: create token: %w", err)
	}

	return user, nil
}

func (u *AuthUserUse) ParseToken(token string) (entity.UserModel, error) {
	user, err := u.Repository.GetUserByToken(token)
	if err != nil {
		return entity.UserModel{}, fmt.Errorf("usecase: parse token: %w", err)
	}

	return user, nil
}

func (u *AuthUserUse) DeleteToken(token string) error {
	return u.Repository.DeleteToken(token)
}

// chicking user's given information
func checkUser(user entity.UserModel) error {
	if _, err := mail.ParseAddress(user.Emailcheck); err != nil {
		return fmt.Errorf("usecase: check user :%w", ErrInvalidEmail)
	}

	if m, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`, user.Emailcheck); !m {
		return fmt.Errorf("usecase: check user :%w", ErrInvalidEmail)
	}

	if user.Password != user.ConfirmPassword {
		log.Println("error - password not the same - createuserandvalidate")
		return fmt.Errorf("usecase: check user: %w", ErrConfirmPassword)
	}

	if len(user.Username) < 4 || len(user.Username) > 15 {
		return fmt.Errorf("usecase: check user: %w", ErrInvalidCharacter)
	}

	for _, w := range user.Username {
		// underscore and dot are allowed
		if w == 46 || w == 95 {
			continue
		}

		if (w < 48 || w > 57) && (w < 65 || w > 90) && (w < 97 || w > 122) {
			return fmt.Errorf("usecase: checkUser err: %w", ErrInvalidCharacter)
		}
	}

	if vP := checkPassword(user.Password); !vP {
		return fmt.Errorf("usecase: checkUser err: %w", ErrInvalidPassword)
	}

	return nil
}

// create new hashed password upon signup
func generateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// checks login in the process of login
func checkPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func checkPassword(s string) bool {
	count := 0
	var passwordlen, number, upperCase, specialSym bool
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
			count++
		case unicode.IsUpper(c):
			upperCase = true
			count++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			specialSym = true
		case unicode.IsLetter(c):
			count++
		default:
			return false
		}
	}
	passwordlen = (count >= 6 && count <= 20)
	if passwordlen && number && upperCase && specialSym {
		return true
	}
	return false
}
