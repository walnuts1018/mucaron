package handler

import (
	"errors"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Upload(c *gin.Context) {
	slog.Debug("upload music")
	user, err := h.getUser(c)
	if err != nil {
		if errors.Is(err, ErrNeedLogin) {
			c.JSON(401, gin.H{
				"error": "need login",
			})
			return
		}
		slog.Error("failed to get user", slog.Any("error", err))
		c.JSON(500, gin.H{
			"error": "failed to get user",
		})
		return
	}

	defer c.Request.Body.Close()

	if err := h.usecase.UploadMusic(c.Request.Context(), user, c.Request.Body); err != nil {
		slog.Error("failed to upload music", slog.Any("error", err))
		c.JSON(500, gin.H{
			"error": "failed to upload music",
		})
		return
	}
}
