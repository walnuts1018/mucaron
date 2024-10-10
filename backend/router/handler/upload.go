package handler

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Upload(c *gin.Context) {
	slog.Debug("upload music")
	user, err := getUser(c)
	if err != nil {
		slog.Error("failed to get user", slog.Any("error", err))
		c.JSON(401, gin.H{
			"error": "failed to get user",
		})
		return
	}

	body, err := c.Request.GetBody()
	if err != nil {
		slog.Error("failed to get body", slog.Any("error", err))
		c.JSON(400, gin.H{
			"error": "failed to get body",
		})
		return
	}
	defer body.Close()

	if err := h.usecase.UploadMusic(c.Request.Context(), user, body); err != nil {
		slog.Error("failed to upload music", slog.Any("error", err))
		c.JSON(500, gin.H{
			"error": "failed to upload music",
		})
		return
	}
}
