package redis

import (
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/walnuts1018/mucaron/backend/config"
	"golang.org/x/exp/slog"
)

func NewSessionStore(cfg config.Config) (sessions.Store, error) {
	slog.Debug("creating redis store")
	store, err := redis.NewStoreWithDB(10, "tcp", fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort), cfg.RedisPassword, fmt.Sprintf("%d", cfg.RedisDB), []byte(cfg.SessionSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to create redis store: %w", err)
	}

	err, rs := redis.GetRedisStore(store)
	if err != nil {
		return nil, fmt.Errorf("failed to get redis store: %w", err)
	}

	if rs == nil {
		return nil, fmt.Errorf("failed to get redis store")
	}

	rs.Options.SameSite = cfg.SessionOptions.SameSite
	rs.Options.HttpOnly = cfg.SessionOptions.HttpOnly
	rs.Options.Secure = cfg.SessionOptions.Secure

	slog.Debug("created redis store")

	return store, nil
}
