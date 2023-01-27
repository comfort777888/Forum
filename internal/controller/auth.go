package controller

import (
	"errors"
	"log"
	"net/http"
	"time"

	"forum/internal/entity"
	"forum/internal/usecase"
)

func (h *handler) Home(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(ctxKeyUser)
	user := u.(entity.UserModel)

	if r.URL.Path != "/" {
		h.errorHandler(w, http.StatusNotFound, "incorrect path")
		log.Printf("error incorrect path\n")
		return
	}
	if r.Method != http.MethodGet {
		h.errorHandler(w, http.StatusMethodNotAllowed, "incorrect method")
		log.Printf("error incorrect method\n")
		return
	}

	var posts []entity.Post
	var err error
	if len(r.URL.Query()) == 0 {
		posts, err = h.usecase.PostsUsecase.GetAllPosts()
		if err != nil {
			log.Printf("error home page: query len: %v\n", err)
			h.errorHandler(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		// get all by filter
		posts, err = h.usecase.PostsUsecase.GetAllPostsFromFilter(user, r.URL.Query())
		if err != nil {
			if errors.Is(err, usecase.ErrUserNotFound) {
				log.Printf("error: %v\n", err)
				h.errorHandler(w, http.StatusUnauthorized, err.Error())
				return
			}
			log.Printf("error %v\n", err)
			h.errorHandler(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	info := entity.Profile{
		Posts: posts,
		User:  user,
	}
	if err := h.execute(w, "ui/template/index.html", info); err != nil {
		h.errorHandler(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// sign-up handler
func (h *handler) signUp(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(ctxKeyUser)
	user := u.(entity.UserModel)

	if user != (entity.UserModel{}) {
		log.Printf("error user unauthorized\n")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	if r.URL.Path != "/auth/sign-up" {
		log.Printf("error incorrect path\n")
		h.errorHandler(w, http.StatusNotFound, "incorrect path")
		return
	}

	switch r.Method {
	case http.MethodGet:
		if err := h.execute(w, "ui/template/register.html", nil); err != nil {
			h.errorHandler(w, http.StatusInternalServerError, err.Error())
			return
		}
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			h.errorHandler(w, http.StatusBadRequest, err.Error())
			return
		}

		username, ok := r.Form["username"]
		if !ok {
			h.errorHandler(w, http.StatusBadRequest, "username field not found")
			return
		}

		email, ok := r.Form["email"]
		if !ok {
			h.errorHandler(w, http.StatusBadRequest, "email field not found")
			return
		}

		password, ok := r.Form["password"]
		if !ok {
			h.errorHandler(w, http.StatusBadRequest, "password field not found")
			return
		}

		confirmPassword, ok := r.Form["confirm_password"]
		if !ok {
			h.errorHandler(w, http.StatusBadRequest, "confirm_password field not found")
			return
		}

		user := entity.UserModel{
			Username:        username[0],
			Email:           email[0],
			Password:        password[0],
			ConfirmPassword: confirmPassword[0],
			Emailcheck:      email[0],
			Usernamecheck:   username[0],
		}

		userCheck := entity.UserModel{
			Emailcheck:    email[0],
			Usernamecheck: username[0],
		}

		if err := h.usecase.CreateUserandValidate(user); err != nil {

			switch {
			case errors.Is(err, usecase.ErrInvalidEmail):
				log.Printf("error: %v", err)
				userCheck.Email = "Enter a valid email address"

			case errors.Is(err, usecase.ErrUserExist):
				log.Printf("error: %v", err)
				userCheck.Username = "Username is already taken"

			case errors.Is(err, usecase.ErrInvalidCharacter):
				log.Printf("error: %v", err)
				userCheck.Username = "Username is not correct. You can use only latin letters, numbers, undescore, dot from 4 to 15 symbols without spaces"

			case errors.Is(err, usecase.ErrEmailExist):
				log.Printf("error: %v", err)
				userCheck.Email = "User with this email already exists"

			case errors.Is(err, usecase.ErrInvalidPassword):
				log.Printf("error: %v", err)
				userCheck.Password = "Your password must consist of english letters, at least 1 Upper case and special symbol, 6 to 20 long without spaces"

			case errors.Is(err, usecase.ErrConfirmPassword):
				log.Printf("error: %v", err)
				userCheck.ConfirmPassword = "password not the same"

			case errors.Is(err, usecase.ErrHashPassword):
				log.Printf("error: %v", err)
				userCheck.Password = "error with password, try another password"

			default:
				log.Println("error create user")
				h.errorHandler(w, http.StatusInternalServerError, err.Error())
				return
			}

			if err := h.execute(w, "ui/template/register.html", userCheck); err != nil {
				h.errorHandler(w, http.StatusInternalServerError, err.Error())
				return
			}
			return
		}

		http.Redirect(w, r, "/auth/sign-in", http.StatusSeeOther)
	default:
		h.errorHandler(w, http.StatusMethodNotAllowed, "incorrect method")
		return
	}
}

func (h *handler) signIn(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(ctxKeyUser).(entity.UserModel)

	if user != (entity.UserModel{}) {
		log.Printf("user not found")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if r.URL.Path != "/auth/sign-in" {
		log.Printf("error incorrect path\n")
		h.errorHandler(w, http.StatusNotFound, "incorrect path")
		return
	}

	switch r.Method {
	case http.MethodGet:
		if err := h.execute(w, "ui/template/login.html", nil); err != nil {
			h.errorHandler(w, http.StatusInternalServerError, err.Error())
			return
		}

	case http.MethodPost:
		userCheck := entity.UserModel{}
		if err := r.ParseForm(); err != nil {
			h.errorHandler(w, http.StatusBadRequest, err.Error())
			return
		}

		username, ok := r.Form["username"]
		if !ok {
			h.errorHandler(w, http.StatusBadRequest, "username field not found")
			return
		}

		password, ok := r.Form["password"]
		if !ok {
			h.errorHandler(w, http.StatusBadRequest, "password field not found")
			return
		}

		user, err := h.usecase.AuthorizationUsecase.CreateToken(username[0], password[0])
		if err != nil {
			switch {
			case errors.Is(err, usecase.ErrUserNotFound):
				log.Printf("error: %v", err)
				userCheck.Username = "No such user, please register"
			case errors.Is(err, usecase.ErrInvalidPassword):
				log.Printf("error: %v", err)
				userCheck.Password = "Password is incorrect"
			default:
				log.Printf("error cannot create token\n")
				h.errorHandler(w, http.StatusInternalServerError, err.Error())
				return
			}

			if err := h.execute(w, "ui/template/login.html", userCheck); err != nil {
				h.errorHandler(w, http.StatusInternalServerError, err.Error())
				return
			}
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "session_cookie",
			Value: user.Token,
			Path:  "/",
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		h.errorHandler(w, http.StatusMethodNotAllowed, "incorrect method")
	}
}

func (h *handler) logout(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/logout" {
		log.Printf("error incorrect path\n")
		h.errorHandler(w, http.StatusNotFound, "incorrect path")
		return
	}
	if r.Method != http.MethodGet {
		log.Printf("error incorrect method\n")
		h.errorHandler(w, http.StatusMethodNotAllowed, "incorrect method")
		return
	}

	cookie, err := r.Cookie("session_cookie")
	if err != nil {
		log.Printf("error in cookie\n")
		h.errorHandler(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.usecase.AuthorizationUsecase.DeleteToken(cookie.Value); err != nil {
		h.errorHandler(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_cookie",
		Value:   "",
		Expires: time.Time{},
		Path:    "/",
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
