package router

import (
	"log/slog"
	"net/http"

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

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(sloggin.NewWithConfig(slog.Default(), sloggin.Config{
		DefaultLevel:     config.LogLevel,
		ClientErrorLevel: slog.LevelWarn,
		ServerErrorLevel: slog.LevelError,

		WithUserAgent:      false,
		WithRequestID:      true,
		WithRequestBody:    false,
		WithRequestHeader:  false,
		WithResponseBody:   false,
		WithResponseHeader: false,
		WithSpanID:         true,
		WithTraceID:        true,

		Filters: []sloggin.Filter{
			sloggin.IgnorePath("/healthz"),
		},
	}))

        sessionStore.Options(sessions.Options{SameSite: http.SameSiteStrictMode})
	r.Use(sessions.Sessions("default", sessionStore))

	r.GET("/healthz", handler.Health)
	apiv1 := r.Group("/api/v1")
	apiv1.Use(otelgin.Middleware(consts.ApplicationName))
	{
		apiv1.POST("/create_user", handler.CreateUser)
		apiv1.POST("/login", handler.Login)
		apiv1.POST("/upload", handler.Upload)

		music := apiv1.Group("/music")
		{
			music.GET("/", handler.GetMusics)
			music.GET("/:id", handler.GetMusic)
			music.GET("/:id/stream/:stream_id", handler.GetMusicStream)
			music.GET("/:id/primary_stream", handler.RedirectMusicPrimaryStream)
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
