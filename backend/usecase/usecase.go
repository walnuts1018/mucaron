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

	CreateMusicWithDependencies(m entity.Music, album *entity.Album, artist *entity.Artist, genre *entity.Genre) error
}

type musicRepository interface {
	CreateMusic(m entity.Music) error
	UpdateMusic(m entity.Music) error
	UpdateMusicStatus(musicID uuid.UUID, status entity.MusicStatus) error
	DeleteMusics(musicIDs []uuid.UUID) error
	GetMusicByID(id uuid.UUID) (entity.Music, error)
	GetMusicByIDs(ids []uuid.UUID) ([]entity.Music, error)
	GetMusicsByUserID(userID uuid.UUID) ([]entity.Music, error)
	GetMusicIDsByUserID(userID uuid.UUID) ([]uuid.UUID, error) 
	UpdateMusicStatuses(musicIDs []uuid.UUID, status entity.MusicStatus) error
}

type userRepository interface {
	CreateUser(u entity.User) error
	UpdateUser(u entity.User) error
	DeleteUser(u entity.User) error
	GetUserByIDs(userIDs []uuid.UUID) ([]entity.User, error)
	GetUserByID(userID uuid.UUID) (entity.User, error)
	GetUserByName(userName string) (entity.User, error)
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
