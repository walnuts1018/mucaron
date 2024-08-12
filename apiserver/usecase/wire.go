package usecase

import (
	"github.com/google/wire"
	"github.com/walnuts1018/mucaron/common/infra/postgres"
)

var Set = wire.NewSet(
	postgres.NewPostgres,
	wire.Bind(new(AlbumRepository), new(*postgres.PostgresClient)),
	wire.Bind(new(ArtistRepository), new(*postgres.PostgresClient)),
	wire.Bind(new(GenreRepository), new(*postgres.PostgresClient)),
	wire.Bind(new(MusicRepository), new(*postgres.PostgresClient)),
	wire.Bind(new(PlaylistRepository), new(*postgres.PostgresClient)),
	wire.Bind(new(UserRepository), new(*postgres.PostgresClient)),
)
