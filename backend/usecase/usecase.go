package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/walnuts1018/mucaron/backend/config"
	"github.com/walnuts1018/mucaron/backend/domain/entity"
)

type AlbumRepository interface{}
type ArtistRepository interface{}
type GenreRepository interface{}
type MusicRepository interface {
	CreateMusic(m entity.Music) error
}
type PlaylistRepository interface{}
type UserRepository interface{}

type Encoder interface {
	Encode(id uuid.UUID, path string, audioOnly bool) (string, error)
}

type MetadataReader interface {
	GetMetadata(ctx context.Context, path string) (entity.RawMusicMetadata, error)
}

type Usecase struct {
	cfg                config.Config
	albumRepository    AlbumRepository
	artistRepository   ArtistRepository
	genreRepository    GenreRepository
	MusicRepository    MusicRepository
	playlistRepository PlaylistRepository
	userRepository     UserRepository
	encoder            Encoder
	metadataReader     MetadataReader
}

func NewUsecase(
	cfg config.Config,
	albumRepository AlbumRepository,
	artistRepository ArtistRepository,
	genreRepository GenreRepository,
	MusicRepository MusicRepository,
	playlistRepository PlaylistRepository,
	userRepository UserRepository,
	encoder Encoder,
	metadataReader MetadataReader,
) Usecase {
	return Usecase{
		cfg,
		albumRepository,
		artistRepository,
		genreRepository,
		MusicRepository,
		playlistRepository,
		userRepository,
		encoder,
		metadataReader,
	}
}
