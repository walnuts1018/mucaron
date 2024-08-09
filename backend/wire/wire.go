//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/walnuts1018/mucaron/config"
	"github.com/walnuts1018/mucaron/infra/postgres"
	"github.com/walnuts1018/mucaron/router"
	"github.com/walnuts1018/mucaron/router/handler"
)

func CreateRouter(
	cfg config.Config,
) (*gin.Engine, error) {
	wire.Build(
		postgres.Set,
		router.NewRouter,
		handler.NewHandler,
	)

	return &gin.Engine{}, nil
}
