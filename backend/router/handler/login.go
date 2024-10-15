package handler

import (
	"log/slog"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/walnuts1018/mucaron/backend/domain/entity"
)

func (h *Handler) Login(c *gin.Context) {
	userName := c.PostForm("user_name")
	inputPassword := c.PostForm("password")
	user, err := h.usecase.Login(c.Request.Context(), userName, entity.RawPassword(inputPassword))
	if err != nil {
		slog.Error("failed to login", slog.Any("error", err), slog.String("user_name", userName))
		c.JSON(401, gin.H{
			"error": "user_name or password is incorrect",
		})
		return
	}

	slog.Debug("login", slog.Any("user_id", user.ID))

	session := sessions.Default(c)
	session.Set(UserIDSessionKey, user.ID.String())
	if err := session.Save(); err != nil {
		slog.Error("failed to save session", slog.Any("error", err))
		c.JSON(500, gin.H{
			"error": "failed to save session",
		})
		return
	}

	c.JSON(200, gin.H{
		"user_id": user.ID,
	})
}
