package handler

import (
	"github.com/walnuts1018/mucaron/config"
)

type Handler struct {
	config       config.Config
}

func NewHandler(config config.Config) (Handler, error) {
	return Handler{
		config:       config,
	}, nil
}
