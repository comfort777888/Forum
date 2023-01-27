package controller

import (
	"forum/internal/usecase"
)

type handler struct {
	usecase *usecase.UseCase
}

func NewHandler(u *usecase.UseCase) *handler {
	return &handler{
		usecase: u,
	}
}
