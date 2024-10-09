package handler

import (
	"github.com/walnuts1018/mucaron/backend/config"
	"github.com/walnuts1018/mucaron/backend/usecase"
)

type Handler struct {
	config  config.Config
	usecase *usecase.Usecase
}

func NewHandler(config config.Config, usecase *usecase.Usecase) (Handler, error) {
	return Handler{
		config,
		usecase,
	}, nil
}
