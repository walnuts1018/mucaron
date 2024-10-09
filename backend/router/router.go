package router

import (
	"log/slog"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"github.com/walnuts1018/mucaron/backend/config"
	"github.com/walnuts1018/mucaron/backend/consts"
	"github.com/walnuts1018/mucaron/backend/router/handler"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func NewRouter(config config.Config, handler handler.Handler, sessionStore sessions.Store) (*gin.Engine, error) {
	if config.LogLevel != slog.LevelDebug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.ContextWithFallback = true
	r.Use(gin.Recovery())
	r.Use(sloggin.New(slog.Default()))
	r.Use(otelgin.Middleware(consts.ApplicationName))

	r.Use(sessions.Sessions("mysession", sessionStore))

	r.GET("/healthz", handler.Health)
	apiv1 := r.Group("/api/v1")
	{
		apiv1.POST("/upload", handler.UploadMusic).Use(checkUserMiddleware())

		music := apiv1.Group("/music")
		{
			music.GET("/", handler.GetMusics)
			music.GET("/:id", handler.GetMusic)
			music.PATCH("/metadata/:id", handler.UpdateMusicMetadata)
			music.POST("/delete", handler.DeleteMusics)
		}

		playlist := apiv1.Group("/playlist")
		{
			playlist.GET("/", handler.GetPlaylists)
			playlist.GET("/:id", handler.GetPlaylist)
			playlist.POST("/", handler.CreatePlaylist)
			playlist.POST("/add", handler.AddMusicToPlaylist)
			playlist.PATCH("/:id", handler.UpdatePlaylist)
			playlist.POST("/delete", handler.DeletePlaylists)
		}
	}

	return r, nil
}
