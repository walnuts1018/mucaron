package router

import (
	"log/slog"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func checkUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID == nil {
			slog.Info("user_id is nil",
				slog.String("path", c.Request.URL.Path),
				slog.String("method", c.Request.Method),
				slog.String("client_ip", c.ClientIP()),
			)
			c.Redirect(302, "/login")
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
