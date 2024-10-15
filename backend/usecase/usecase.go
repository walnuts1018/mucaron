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

type EntityRepository interface {
	musicRepository
	userRepository

	CreateMusicWithDependencies(ctx context.Context, m entity.Music, album *entity.Album, artist *entity.Artist, genre *entity.Genre) error
}

type musicRepository interface {
	CreateMusic(ctx context.Context, m entity.Music) error
	UpdateMusic(ctx context.Context, m entity.Music) error
	UpdateMusicStatus(ctx context.Context, musicID uuid.UUID, status entity.MusicStatus) error
	DeleteMusics(ctx context.Context, musicIDs []uuid.UUID) error
	GetMusicByID(ctx context.Context, id uuid.UUID) (entity.Music, error)
	GetMusicByIDs(ctx context.Context, ids []uuid.UUID) ([]entity.Music, error)
	GetMusicsByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Music, error)
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
