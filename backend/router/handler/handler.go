package handler

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/walnuts1018/mucaron/backend/config"
	"github.com/walnuts1018/mucaron/backend/domain"
	"github.com/walnuts1018/mucaron/backend/domain/entity"
	"github.com/walnuts1018/mucaron/backend/usecase"
)

const (
	UserIDSessionKey = "user_id"
)

var (
	ErrLoginRequired = errors.New("login required")
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

func (h *Handler) getUser(c *gin.Context) (entity.User, error) {
	session := sessions.Default(c)
	userIDStr, ok := session.Get(UserIDSessionKey).(string)
	if !ok || userIDStr == "" {
		slog.Info("user not found in session",
			slog.String("path", c.Request.URL.Path),
			slog.String("method", c.Request.Method),
			slog.String("client_ip", c.ClientIP()),
		)
		return entity.User{}, ErrLoginRequired
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return entity.User{}, err
	}

	user, err := h.usecase.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			slog.Info("failed to get user by id",
				slog.String("path", c.Request.URL.Path),
				slog.String("method", c.Request.Method),
				slog.String("client_ip", c.ClientIP()),
			)
			return entity.User{}, ErrLoginRequired
		}
		return entity.User{}, fmt.Errorf("failed to get user by id: %w", err)
	}

	return user, nil
}
