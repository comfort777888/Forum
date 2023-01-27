package controller

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"forum/internal/entity"
)

func (h *handler) likeComment(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(ctxKeyUser)
	user := u.(entity.UserModel)

	if user == (entity.UserModel{}) {
		h.errorHandler(w, http.StatusUnauthorized, "user unauthorized")
		return
	}

	if r.Method != http.MethodPost {
		log.Printf("error incorrect method\n")
		h.errorHandler(w, http.StatusMethodNotAllowed, "incorrect method")
		return
	}

	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/comment/like/"))
	if err != nil {
		h.errorHandler(w, http.StatusNotFound, "incorrect path")
		return
	}

	comment, err := h.usecase.CommentsUsecase.GetCommentById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.errorHandler(w, http.StatusNotFound, "incorrect path")
			return
		}

		h.errorHandler(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.usecase.CommentsUsecase.LikeComment(id, user.Username); err != nil {
		h.errorHandler(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%d", comment.PostId), http.StatusSeeOther)
}

func (h *handler) dislikeComment(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(ctxKeyUser)
	user := u.(entity.UserModel)

	if user == (entity.UserModel{}) {
		h.errorHandler(w, http.StatusUnauthorized, "user unauthorized")
		return
	}
	if r.Method != http.MethodPost {
		h.errorHandler(w, http.StatusMethodNotAllowed, "incorrect method")
		return
	}

	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/comment/dislike/"))
	if err != nil {
		h.errorHandler(w, http.StatusNotFound, "incorrect path")
		return
	}

	comment, err := h.usecase.CommentsUsecase.GetCommentById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.errorHandler(w, http.StatusNotFound, "incorrect path")
			return
		}
		h.errorHandler(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.usecase.CommentsUsecase.DislikeComment(id, user.Username); err != nil {
		h.errorHandler(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%d", comment.PostId), http.StatusSeeOther)
}
