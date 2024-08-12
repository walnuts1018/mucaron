//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/walnuts1018/mucaron/apiserver/router"
	"github.com/walnuts1018/mucaron/apiserver/router/handler"
	"github.com/walnuts1018/mucaron/apiserver/usecase"
	"github.com/walnuts1018/mucaron/common/config"
)

func CreateRouter(
	cfg config.Config,
) (*gin.Engine, error) {
	wire.Build(
		usecase.Set,
		usecase.NewUsecase,
		router.NewRouter,
		handler.NewHandler,
	)

	return &gin.Engine{}, nil
}
