package handler

import (
	"errors"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/walnuts1018/mucaron/backend/domain/entity"
	"github.com/walnuts1018/mucaron/backend/usecase"
)

func (h *Handler) CreateUser(c *gin.Context) {
	userName := c.PostForm("user_name")
	if userName == "" {
		slog.Error("user_name is required")
		c.JSON(400, gin.H{
			"error": "user_name is required",
		})
		return
	}

	inputPassword := c.PostForm("password")
	if inputPassword == "" {
		slog.Error("password is required")
		c.JSON(400, gin.H{
			"error": "password is required",
		})
		return
	}

	user, err := h.usecase.CreateUser(c, userName, entity.RawPassword(inputPassword))
	if err != nil {
		if errors.Is(err, usecase.ErrUserExists) {
			c.JSON(400, gin.H{
				"error": "user already exists",
			})
			return
		}
		if errors.Is(err, entity.ErrInvalidPassword) {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(400, gin.H{
			"error": "failed to create user",
		})
		return
	}

	c.JSON(200, gin.H{
		"user_id": user.ID,
	})
}
