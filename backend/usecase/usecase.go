package usecase

import (
	"context"
	"io"

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

type ObjectStorage interface {
	GetObjectURL(ctx context.Context, path string) (string, error)
	UploadObject(ctx context.Context, objectName string, data io.Reader, size int64) error
	UploadDirectory(ctx context.Context, objectBaseDir string, localDir string) error
	DeleteObject(ctx context.Context, objectName string) error
}

type Encoder interface {
	Encode(id uuid.UUID, path string, audioOnly bool) (string, error)
}

type MetadataReader interface {
	GetMetadata(ctx context.Context, path string) (entity.RawMusicMetadata, error)
}

type queueItem struct {
	ID        uuid.UUID
	Path      string
	AudioOnly bool
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
	objectStorage      ObjectStorage

	encodeQueue chan queueItem
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
	objectStorage ObjectStorage,
) Usecase {
	encodeQueue := make(chan queueItem)

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
		objectStorage,
		encodeQueue,
	}
}
