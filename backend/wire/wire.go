//go:build wireinject
// +build wireinject

package wire

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/walnuts1018/mucaron/backend/config"
	"github.com/walnuts1018/mucaron/backend/infra/ffmpeg"
	"github.com/walnuts1018/mucaron/backend/infra/ffprobe"
	"github.com/walnuts1018/mucaron/backend/infra/minio"
	"github.com/walnuts1018/mucaron/backend/infra/postgres"
	"github.com/walnuts1018/mucaron/backend/infra/redis"
	"github.com/walnuts1018/mucaron/backend/router"
	"github.com/walnuts1018/mucaron/backend/router/handler"
	"github.com/walnuts1018/mucaron/backend/usecase"
)

func CreateRouter(
	ctx context.Context,
	cfg config.Config,
) (*gin.Engine, error) {
	wire.Build(
		postgresSet,
		ffmpegSet,
		ffprobeSet,
		minioSet,
		usecase.NewUsecase,
		handler.NewHandler,
		redis.NewSessionStore,
		router.NewRouter,
	)

	return &gin.Engine{}, nil
}

var postgresSet = wire.NewSet(
	postgres.NewPostgres,
	wire.Bind(new(usecase.EntityRepository), new(*postgres.PostgresClient)),
)

var ffmpegSet = wire.NewSet(
	ffmpeg.NewFFMPEG,
	wire.Bind(new(usecase.Encoder), new(*ffmpeg.FFMPEG)),
)

var ffprobeSet = wire.NewSet(
	ffprobe.NewFFProbe,
	wire.Bind(new(usecase.MetadataReader), new(ffprobe.FFProbe)),
)

var minioSet = wire.NewSet(
	minio.NewMinIOClient,
	wire.Bind(new(usecase.ObjectStorage), new(*minio.MinIO)),
)
