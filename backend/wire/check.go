package wire

import (
	"github.com/walnuts1018/mucaron/backend/infra/ffmpeg"
	"github.com/walnuts1018/mucaron/backend/infra/ffprobe"
	"github.com/walnuts1018/mucaron/backend/infra/minio"
	"github.com/walnuts1018/mucaron/backend/infra/postgres"
	"github.com/walnuts1018/mucaron/backend/usecase"
)

var _ usecase.AlbumRepository = &postgres.PostgresClient{}
var _ usecase.ArtistRepository = &postgres.PostgresClient{}
var _ usecase.GenreRepository = &postgres.PostgresClient{}
var _ usecase.MusicRepository = &postgres.PostgresClient{}
var _ usecase.PlaylistRepository = &postgres.PostgresClient{}
var _ usecase.UserRepository = &postgres.PostgresClient{}

var _ usecase.Encoder = &ffmpeg.FFMPEG{}

var _ usecase.MetadataReader = ffprobe.FFProbe{}

var _ usecase.ObjectStorage = &minio.MinIO{}
