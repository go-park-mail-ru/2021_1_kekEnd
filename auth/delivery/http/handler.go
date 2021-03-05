package http

import "github.com/go-park-mail-ru/2021_1_kekEnd/auth"

type Handler struct {
	useCase auth.UseCase
}

func NewHandler(useCase auth.UseCase) *Handler {
	return &Handler {
		useCase: useCase,
	}
}
