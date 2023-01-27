package controller

import (
	"net/http"
)

func SetupRouter(h *handler) *http.ServeMux {
	router := http.NewServeMux()
	styles := http.FileServer(http.Dir("ui/static/"))

	router.Handle("/static/", http.StripPrefix("/static/", styles))

	router.HandleFunc("/", h.verification(h.Home))

	router.HandleFunc("/auth/sign-up", h.verification(h.signUp))
	router.HandleFunc("/auth/sign-in", h.verification(h.signIn))

	router.HandleFunc("/auth/logout", h.verification(h.logout))
	router.HandleFunc("/post/create", h.verification(h.createPost))
	router.HandleFunc("/post/like/", h.verification(h.likePost))
	router.HandleFunc("/post/dislike/", h.verification(h.disLikePost))
	router.HandleFunc("/post/", h.verification(h.postPage))

	router.HandleFunc("/comment/like/", h.verification(h.likeComment))
	router.HandleFunc("/comment/dislike/", h.verification(h.dislikeComment))

	router.HandleFunc("/profile/", h.verification(h.userProfile))

	return router
}
