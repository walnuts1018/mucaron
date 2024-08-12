package handler

import (
	"github.com/walnuts1018/mucaron/apiserver/usecase"
	"github.com/walnuts1018/mucaron/common/config"
)

type Handler struct {
	config  config.Config
	usecase usecase.Usecase
}

func NewHandler(config config.Config, usecase usecase.Usecase) (Handler, error) {
	return Handler{
		config,
		usecase,
	}, nil
}
