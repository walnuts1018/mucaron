// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

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

// Injectors from wire.go:

func CreateRouter(ctx context.Context, cfg config.Config) (*gin.Engine, error) {
	postgresClient, err := postgres.NewPostgres(ctx, cfg)
	if err != nil {
		return nil, err
	}
	ffmpegFFMPEG, err := ffmpeg.NewFFMPEG(cfg)
	if err != nil {
		return nil, err
	}
	ffProbe := ffprobe.NewFFProbe()
	minIO, err := minio.NewMinIOClient(cfg)
	if err != nil {
		return nil, err
	}
	usecaseUsecase, err := usecase.NewUsecase(cfg, postgresClient, ffmpegFFMPEG, ffProbe, minIO)
	if err != nil {
		return nil, err
	}
	handlerHandler, err := handler.NewHandler(cfg, usecaseUsecase)
	if err != nil {
		return nil, err
	}
	store, err := redis.NewSessionStore(cfg)
	if err != nil {
		return nil, err
	}
	engine, err := router.NewRouter(cfg, handlerHandler, store)
	if err != nil {
		return nil, err
	}
	return engine, nil
}

// wire.go:

var postgresSet = wire.NewSet(postgres.NewPostgres, wire.Bind(new(usecase.EntityRepository), new(*postgres.PostgresClient)))

var ffmpegSet = wire.NewSet(ffmpeg.NewFFMPEG, wire.Bind(new(usecase.Encoder), new(*ffmpeg.FFMPEG)))

var ffprobeSet = wire.NewSet(ffprobe.NewFFProbe, wire.Bind(new(usecase.MetadataReader), new(ffprobe.FFProbe)))

var minioSet = wire.NewSet(minio.NewMinIOClient, wire.Bind(new(usecase.ObjectStorage), new(*minio.MinIO)))
