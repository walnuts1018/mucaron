package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/walnuts1018/mucaron/backend/domain/entity"
	"github.com/walnuts1018/mucaron/backend/usecase"
)

func (h *Handler) CreateUser(c *gin.Context) {
	userName := c.PostForm("user_name")
	inputPassword := c.PostForm("password")
	user, err := h.usecase.CreateUser(userName, entity.RawPassword(inputPassword))
	if err != nil {
		if errors.Is(err, usecase.ErrUserExists) {
			c.JSON(400, gin.H{
				"error": "user already exists",
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
