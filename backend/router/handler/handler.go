package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/walnuts1018/mucaron/backend/config"
	"github.com/walnuts1018/mucaron/backend/domain/entity"
	"github.com/walnuts1018/mucaron/backend/router/middleware"
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

func getUser(c *gin.Context) (entity.User, error) {
	user, ok := c.Get(middleware.UserKey)
	if !ok {
		return entity.User{}, fmt.Errorf("failed to get user")
	}

	u, ok := user.(entity.User)
	if !ok {
		return entity.User{}, fmt.Errorf("failed to assert user")
	}

	return u, nil
}
