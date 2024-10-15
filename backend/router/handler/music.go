package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/walnuts1018/mucaron/backend/domain"
)

func (h *Handler) GetMusics(c *gin.Context) {
	user, err := h.getUser(c)
	if err != nil {
		if errors.Is(err, ErrLoginRequired) {
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

	musics, err := h.usecase.GetMusics(c.Request.Context(), user)
	if err != nil {
		slog.Error("failed to get musics", slog.Any("error", err))
		c.JSON(500, gin.H{
			"error": "failed to get musics",
		})
		return
	}

	c.JSON(200, gin.H{
		"musics": musics,
	})
}

func (h *Handler) GetMusic(c *gin.Context) {}

func (h *Handler) UpdateMusicMetadata(c *gin.Context) {}

type DeleteMusicsRequest struct {
	IDs []string `json:"ids"`
}

func (h *Handler) DeleteMusics(c *gin.Context) {
	user, err := h.getUser(c)
	if err != nil {
		if errors.Is(err, ErrLoginRequired) {
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

	if err := h.usecase.DeleteMusics(c.Request.Context(), user, uuids); err != nil {
		if errors.Is(err, domain.ErrAccessDenied) {
			c.JSON(403, gin.H{
				"error": "access denied",
			})
			return
		}
		slog.Error("failed to delete musics", slog.Any("error", err))
		c.JSON(500, gin.H{
			"error": "failed to delete musics",
		})
		return
	}

	c.JSON(200, gin.H{})
}

func (h *Handler) RedirectMusicPrimaryStream(c *gin.Context) {
	musicID := c.Param("id")
	if musicID == "" {
		c.JSON(400, gin.H{
			"error": "id is required",
		})
		return
	}

	musicUUID, err := uuid.Parse(musicID)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid id",
		})
		return
	}

	user, err := h.getUser(c)
	if err != nil {
		if errors.Is(err, ErrLoginRequired) {
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

	url, err := h.usecase.GetPrimaryStreamM3U8URL(c.Request.Context(), user, musicUUID)
	if err != nil {
		if errors.Is(err, domain.ErrAccessDenied) {
			c.JSON(403, gin.H{
				"error": "access denied",
			})
			return
		}
		slog.Error("failed to get primary stream", slog.Any("error", err))
		c.JSON(500, gin.H{
			"error": "failed to get primary stream",
		})
		return
	}

	c.Redirect(http.StatusFound, url.String())
}

func (h *Handler) GetMusicStream(c *gin.Context) {
	musicID := c.Param("id")
	if musicID == "" {
		c.JSON(400, gin.H{
			"error": "id is required",
		})
		return
	}

	streamID := c.Param("stream_id")
	if streamID == "" {
		c.JSON(400, gin.H{
			"error": "stream_id is required",
		})
		return
	}

	user, err := h.getUser(c)
	if err != nil {
		if errors.Is(err, ErrLoginRequired) {
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

	musicUUID, err := uuid.Parse(musicID)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid id",
		})
		return
	}

	body, err := h.usecase.GetStreamM3U8(c.Request.Context(), user, musicUUID, streamID)
	if err != nil {
		if errors.Is(err, domain.ErrAccessDenied) {
			c.JSON(403, gin.H{
				"error": "access denied",
			})
			return
		}
		slog.Error("failed to get stream", slog.Any("error", err))
		c.JSON(500, gin.H{
			"error": "failed to get stream",
		})
		return
	}

	c.Data(200, "application/vnd.apple.mpegurl", []byte(body))
}
