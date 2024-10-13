package handler

import (
	"errors"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/walnuts1018/mucaron/backend/domain"
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

	fh, err := c.FormFile("file")
	if err != nil {
		slog.Error("failed to get file", slog.Any("error", err))
		c.JSON(400, gin.H{
			"error": "file is required",
		})
		return
	}
	f, err := fh.Open()
	if err != nil {
		slog.Error("failed to open file", slog.Any("error", err))
		c.JSON(500, gin.H{
			"error": "failed to open file",
		})
		return
	}

	if err := h.usecase.UploadMusic(c.Request.Context(), user, f, fh.Filename); err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			c.JSON(400, gin.H{
				"error": "music already exists",
			})
			return
		}
		slog.Error("failed to upload music", slog.Any("error", err))
		c.JSON(500, gin.H{
			"error": "failed to upload music",
		})
		return
	}
}