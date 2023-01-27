package controller

import (
	"errors"
	"net/http"
	"strings"

	"forum/internal/entity"
	"forum/internal/usecase"
)

func (h *handler) userProfile(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(ctxKeyUser)
	user := u.(entity.UserModel)

	username := strings.TrimPrefix(r.URL.Path, "/profile/")

	userP, err := h.usecase.UsersUsecase.GetUserByName(username)
	if err != nil {
		h.errorHandler(w, http.StatusNotFound, err.Error())
		return
	}

	if r.Method != http.MethodGet {
		h.errorHandler(w, http.StatusMethodNotAllowed, "invalid method")
		return
	}
	info := entity.Profile{}
	if len(r.URL.Query()) == 0 {
		info = entity.Profile{
			User:        user,
			ProfileUser: userP,
		}
	} else {
		posts, err := h.usecase.UsersUsecase.GetPostsByName(userP.Username, r.URL.Query())
		if err != nil {
			if errors.Is(err, usecase.ErrInvalidQuery) {
				h.errorHandler(w, http.StatusBadRequest, err.Error())
				return
			}
			h.errorHandler(w, http.StatusInternalServerError, err.Error())
			return
		}
		info = entity.Profile{
			User:        user,
			ProfileUser: userP,
			Posts:       posts,
		}
	}

	if err := h.execute(w, "ui/template/profile.html", info); err != nil {
		h.errorHandler(w, http.StatusInternalServerError, err.Error())
	}
}
