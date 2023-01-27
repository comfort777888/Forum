package controller

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"forum/internal/entity"
)

type ctxKey int8

const ctxKeyUser ctxKey = iota

func (h *handler) verification(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_cookie")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				// log.Printf("error: no cookie\n")
				next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, entity.UserModel{})))
				return
			}
			h.errorHandler(w, http.StatusBadRequest, err.Error())
			return
		}

		user, err := h.usecase.ParseToken(cookie.Value)
		if err != nil {
			// log.Printf("error: parse token %v\n", err)
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, entity.UserModel{})))
			return
		}

		if user.ExpirationTime.Before(time.Now()) {
			if err = h.usecase.DeleteToken(cookie.Value); err != nil {
				log.Printf("middleware: delete token: %v\n", err)
				h.errorHandler(w, http.StatusInternalServerError, err.Error())
				return
			}

			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, entity.UserModel{})))
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, user)))
	}
}
