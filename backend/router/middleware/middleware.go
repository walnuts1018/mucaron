package middleware

import (
	"errors"
	"log/slog"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/walnuts1018/mucaron/backend/domain"
	"github.com/walnuts1018/mucaron/backend/usecase"
)

const (
	UserKey = "user"
)

type Middleware struct {
	usecase *usecase.Usecase
}

func NewMiddleware(usecase *usecase.Usecase) Middleware {
	return Middleware{
		usecase,
	}
}

func (m Middleware) CheckUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		slog.Debug("check user middleware")
		session := sessions.Default(c)
		userID, ok := session.Get("user_id").(string)
		if !ok || userID == "" {
			slog.Info("user_id is nil",
				slog.String("path", c.Request.URL.Path),
				slog.String("method", c.Request.Method),
				slog.String("client_ip", c.ClientIP()),
			)
			c.Redirect(302, "/login")
			c.Abort()
			return
		}

		user, err := m.usecase.GetUserByID(userID)
		if err != nil {
			if errors.Is(err, domain.ErrNotFount) {
				slog.Info("failed to get user by id", slog.String("user_id", userID), slog.Any("error", err))
				c.Redirect(302, "/login")
				c.Abort()
			}
			slog.Error("failed to get user by id", slog.String("user_id", userID), slog.Any("error", err))
			c.JSON(500, gin.H{
				"error": "failed to get user",
			})
			c.Abort()
			return
		}

		c.Set(UserKey, user)
		c.Next()
	}
}
