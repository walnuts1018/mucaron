package usecase

import (
	"context"
	"io"
	"net/url"
	"sync"

	"github.com/google/uuid"
	"github.com/walnuts1018/mucaron/backend/config"
	"github.com/walnuts1018/mucaron/backend/domain/entity"
)

//go:generate go run -mod=mod go.uber.org/mock/mockgen -source=usecase.go -destination=mock_usecase/mock.go

type EntityRepository interface {
	musicRepository
	userRepository
	artistRepository
	genericRepository
	albumRepository

	Transaction(ctx context.Context, f func(ctx context.Context) error) error
}

type musicRepository interface {
	CreateMusic(ctx context.Context, m entity.Music) error
	UpdateMusic(ctx context.Context, m entity.Music) error
	UpdateMusicStatus(ctx context.Context, musicID uuid.UUID, status entity.MusicStatus) error
	DeleteMusics(ctx context.Context, musicIDs []uuid.UUID) error
	HardDeleteMusic(ctx context.Context, music entity.Music) error
	GetMusicByID(ctx context.Context, id uuid.UUID) (entity.Music, error)
	GetMusicByIDs(ctx context.Context, ids []uuid.UUID) ([]entity.Music, error)
	GetMusicsByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Music, error)
	GetMusicByFileHash(ctx context.Context, userID uuid.UUID, fileHash string, m *entity.Music) error
	GetMusicIDsByUserID(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
	UpdateMusicStatuses(ctx context.Context, musicIDs []uuid.UUID, status entity.MusicStatus) error
}

type userRepository interface {
	CreateUser(ctx context.Context, u entity.User) error
	UpdateUser(ctx context.Context, u entity.User) error
	DeleteUser(ctx context.Context, u entity.User) error
	GetUserByIDs(ctx context.Context, userIDs []uuid.UUID) ([]entity.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (entity.User, error)
	GetUserByName(ctx context.Context, userName string) (entity.User, error)
}

type artistRepository interface {
	CreateArtist(ctx context.Context, a entity.Artist) error
	UpdateArtist(ctx context.Context, a entity.Artist) error
	DeleteArtist(ctx context.Context, a entity.Artist) error
	GetArtistByIDs(ctx context.Context, ids []uuid.UUID) ([]entity.Artist, error)
	GetArtistByID(ctx context.Context, id uuid.UUID) (entity.Artist, error)
	GetArtistByName(ctx context.Context, ownerID uuid.UUID, name string) (entity.Artist, error)
}

type genericRepository interface {
	CreateGenre(ctx context.Context, g entity.Genre) error
	UpdateGenre(ctx context.Context, g entity.Genre) error
	DeleteGenre(ctx context.Context, g entity.Genre) error
	GetGenreByIDs(ctx context.Context, ids []uuid.UUID) ([]entity.Genre, error)
	GetGenreByID(ctx context.Context, id uuid.UUID) (entity.Genre, error)
	GetGenreByName(ctx context.Context, ownerID uuid.UUID, name string) (entity.Genre, error)
}

type albumRepository interface {
	CreateAlbum(ctx context.Context, a entity.Album) error
	UpdateAlbum(ctx context.Context, a entity.Album) error
	DeleteAlbums(ctx context.Context, a []entity.Album) error
	GetAlbumByIDs(ctx context.Context, ids []uuid.UUID) ([]entity.Album, error)
	GetAlbumByID(ctx context.Context, id uuid.UUID) (entity.Album, error)
	GetAlbumsByNameAndArtist(ctx context.Context, ownerID uuid.UUID, albumName string, artist entity.Artist) ([]entity.Album, error)
}

type ObjectStorage interface {
	GetObjectURL(ctx context.Context, objectName string, cacheControl string) (*url.URL, error)
	// GetObjectURLs(ctx context.Context, objectName []string, cacheControl string) ([]*url.URL, error)

	GetObject(ctx context.Context, objectName string) (io.ReadCloser, error)
	UploadObject(ctx context.Context, objectName string, data io.Reader, size int64) error
	UploadDirectory(ctx context.Context, objectBaseDir string, localDir string) error
	DeleteObject(ctx context.Context, objectName string) error
}

type Encoder interface {
	Encode(id string, path string, audioOnly bool) (string, error)
	GetOutDirPrefix() string
}

type MetadataReader interface {
	GetMetadata(ctx context.Context, path string) (entity.RawMusicMetadata, error)
}

type Usecase struct {
	cfg              config.Config
	entityRepository EntityRepository
	encoder          Encoder
	metadataReader   MetadataReader
	objectStorage    ObjectStorage

	encodeMutex sync.Mutex
}

func NewUsecase(
	cfg config.Config,
	entityRepository EntityRepository,
	encoder Encoder,
	metadataReader MetadataReader,
	objectStorage ObjectStorage,
) (*Usecase, error) {
	u := Usecase{
		cfg,
		entityRepository,
		encoder,
		metadataReader,
		objectStorage,
		sync.Mutex{},
	}

	if err := u.EncodeSuspended(context.Background()); err != nil {
		return nil, err
	}

	return &u, nil
}
