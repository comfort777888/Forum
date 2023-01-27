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
	"forum/internal/usecase"
)

func (h *handler) createPost(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(ctxKeyUser)
	user := u.(entity.UserModel)

	if user == (entity.UserModel{}) {
		h.errorHandler(w, http.StatusUnauthorized, "user unauthorized")
		return
	}

	if r.URL.Path != "/post/create" {
		h.errorHandler(w, http.StatusNotFound, "incorrect path")
		return
	}

	switch r.Method {
	case http.MethodGet:
		info := entity.Profile{
			User: user,
		}
		if err := h.execute(w, "ui/template/createPost.html", info); err != nil {
			h.errorHandler(w, http.StatusInternalServerError, err.Error())
		}

	case http.MethodPost:
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			h.errorHandler(w, http.StatusBadRequest, err.Error())
			return
		}

		title, ok := r.Form["title"]
		if !ok {
			h.errorHandler(w, http.StatusBadRequest, "title field not found ")
			return
		}
		content, ok := r.Form["content"]
		if !ok {
			h.errorHandler(w, http.StatusBadRequest, "content field not found")
			return
		}
		category, ok := r.Form["categories"]
		if !ok {
			h.errorHandler(w, http.StatusBadRequest, "empty category field")
			return
		}

		post := entity.Post{
			Title:      title[0],
			Content:    content[0],
			PostAuthor: user.Username,
			Category:   category,
		}

		if err := h.usecase.PostsUsecase.CreatePost(post); err != nil {
			if errors.Is(err, usecase.ErrInvalidContent) ||
				errors.Is(err, usecase.ErrInvalidContentLength) ||
				errors.Is(err, usecase.ErrInvalidTitleLength) ||
				errors.Is(err, usecase.ErrNoContent) ||
				errors.Is(err, usecase.ErrNoTitle) {
				h.errorHandler(w, http.StatusBadRequest, err.Error())
				return
			}

			h.errorHandler(w, http.StatusInternalServerError, err.Error())
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		h.errorHandler(w, http.StatusMethodNotAllowed, "incorrect method")
	}
}

func (h *handler) likePost(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(ctxKeyUser)
	user := u.(entity.UserModel)

	if user == (entity.UserModel{}) {
		h.errorHandler(w, http.StatusUnauthorized, "user unauthorized")
		return
	}

	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/post/like/"))
	if err != nil {
		h.errorHandler(w, http.StatusNotFound, "incorrect path")
		return
	}

	if r.Method != http.MethodPost {
		h.errorHandler(w, http.StatusMethodNotAllowed, "incorrect method")
		return
	}

	if err := h.usecase.PostsVoterUsecase.LikePost(id, user.Username); err != nil {
		h.errorHandler(w, http.StatusInternalServerError, err.Error())
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%d", id), http.StatusSeeOther)
}

func (h *handler) disLikePost(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(ctxKeyUser)
	user := u.(entity.UserModel)

	if user == (entity.UserModel{}) {
		h.errorHandler(w, http.StatusUnauthorized, "user unauthorized")
		return
	}

	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/post/dislike/"))
	if err != nil {
		h.errorHandler(w, http.StatusNotFound, "incorrect path")
		return
	}

	if r.Method != http.MethodPost {
		h.errorHandler(w, http.StatusMethodNotAllowed, "incorrect method")
		return
	}

	if err := h.usecase.PostsVoterUsecase.DisikePost(id, user.Username); err != nil {
		h.errorHandler(w, http.StatusInternalServerError, err.Error())
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%d", id), http.StatusSeeOther)
}

func (h *handler) postPage(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(ctxKeyUser)
	user := u.(entity.UserModel)

	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/post/"))
	if err != nil {
		h.errorHandler(w, http.StatusBadRequest, err.Error())
		return
	}

	post, err := h.usecase.PostsUsecase.GetPostsById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.errorHandler(w, http.StatusNotFound, "incorrect path")
			return
		}

		h.errorHandler(w, http.StatusInternalServerError, err.Error())
		return
	}

	switch r.Method {
	case http.MethodGet:
		comments, err := h.usecase.CommentsUsecase.GetCommentsByPostId(post.PostId)
		if err != nil {
			h.errorHandler(w, http.StatusInternalServerError, err.Error())
			return
		}
		postLikes, err := h.usecase.GetPostLikes(post.PostId)
		if err != nil {
			h.errorHandler(w, http.StatusInternalServerError, err.Error())
			return
		}
		postDislikes, err := h.usecase.GetPostDislikes(post.PostId)
		if err != nil {
			h.errorHandler(w, http.StatusInternalServerError, err.Error())
			return
		}
		commentsLikes, err := h.usecase.CommentsLikes(post.PostId)
		if err != nil {
			h.errorHandler(w, http.StatusInternalServerError, err.Error())
			return
		}
		commentsDislikes, err := h.usecase.CommentsDislikes(post.PostId)
		if err != nil {
			h.errorHandler(w, http.StatusInternalServerError, err.Error())
			return
		}

		info := entity.Profile{
			Post:             post,
			PostLikes:        postLikes,
			PostDislikes:     postDislikes,
			User:             user,
			Comments:         comments,
			CommentsLikes:    commentsLikes,
			CommentsDislikes: commentsDislikes,
		}
		if err := h.execute(w, "ui/template/post.html", info); err != nil {
			h.errorHandler(w, http.StatusInternalServerError, err.Error())
		}

	case http.MethodPost:
		if user == (entity.UserModel{}) {
			log.Printf("error: user is not received")
			h.errorHandler(w, http.StatusUnauthorized, "user unauthorized")
			return
		}

		if err := r.ParseForm(); err != nil {
			h.errorHandler(w, http.StatusBadRequest, err.Error())
			return
		}

		comment, ok := r.Form["comment"]
		if !ok {
			h.errorHandler(w, http.StatusBadRequest, "comment field not found")
			return
		}

		nComment := entity.Comments{
			Content: comment[0],
			Author:  user.Username,
			PostId:  post.PostId,
		}

		if err := h.usecase.CommentsUsecase.CreateComment(nComment); err != nil {
			if errors.Is(err, usecase.ErrInvalidContentLength) ||
				errors.Is(err, usecase.ErrInvalidCharacter) {
				h.errorHandler(w, http.StatusBadRequest, err.Error())
				return
			}

			h.errorHandler(w, http.StatusInternalServerError, err.Error())
			return
		}
		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)

	default:
		h.errorHandler(w, http.StatusMethodNotAllowed, "incorrect method")
	}
}
