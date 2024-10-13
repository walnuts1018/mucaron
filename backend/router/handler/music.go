package handler

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) GetMusics(c *gin.Context) {}

func (h *Handler) GetMusic(c *gin.Context) {}

func (h *Handler) UpdateMusicMetadata(c *gin.Context) {}

type DeleteMusicsRequest struct {
	IDs []string `json:"ids"`
}

func (h *Handler) DeleteMusics(c *gin.Context) {
	var req DeleteMusicsRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request",
		})
		return
	}

	if len(req.IDs) == 0 {
		c.JSON(400, gin.H{
			"error": "ids is required",
		})
		return
	}

	uuids := make([]uuid.UUID, 0, len(req.IDs))
	for _, id := range req.IDs {
		u, err := uuid.Parse(id)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "invalid id",
			})
		}
		uuids = append(uuids, u)
	}

	if err := h.usecase.DeleteMusics(uuids); err != nil {
		slog.Error("failed to delete musics", slog.Any("error", err))
		c.JSON(500, gin.H{
			"error": "failed to delete musics",
		})
		return
	}

	c.JSON(200, gin.H{})
}
