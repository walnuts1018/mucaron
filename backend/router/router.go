package router

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"github.com/walnuts1018/mucaron/config"
	"github.com/walnuts1018/mucaron/router/handler"
)

func NewRouter(config config.Config, handler handler.Handler) (*gin.Engine, error) {
	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(sloggin.New(slog.Default()))

	if config.LogLevel != slog.LevelDebug {
		gin.SetMode(gin.ReleaseMode)
	}

	r.GET("/healthz", handler.Health)

	apiv1 := r.Group("/api/v1")
	{
		music := apiv1.Group("/music")
		{
			music.GET("/", handler.GetMusics)
			music.GET("/:id", handler.GetMusic)
			music.POST("/upload", handler.UploadMusic)
			music.PATCH("/metadata/:id", handler.UpdateMusicMetadata)
			music.POST("/", handler.DeleteMusics)
		}

		playlist := apiv1.Group("/playlist")
		{
			playlist.GET("/", handler.GetPlaylists)
			playlist.GET("/:id", handler.GetPlaylist)
			playlist.POST("/", handler.CreatePlaylist)
			playlist.POST("/add", handler.AddMusicToPlaylist)
			playlist.PATCH("/:id", handler.UpdatePlaylist)
			playlist.POST("/", handler.DeletePlaylists)
		}
	}

	return r, nil
}
