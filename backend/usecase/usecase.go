package usecase

import (
	"github.com/google/uuid"
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

type Usecase struct {
	albumRepository    AlbumRepository
	artistRepository   ArtistRepository
	genreRepository    GenreRepository
	MusicRepository    MusicRepository
	playlistRepository PlaylistRepository
	userRepository     UserRepository
	encoder            Encoder
}

func NewUsecase(
	albumRepository AlbumRepository,
	artistRepository ArtistRepository,
	genreRepository GenreRepository,
	MusicRepository MusicRepository,
	playlistRepository PlaylistRepository,
	userRepository UserRepository,
	encoder Encoder,
) Usecase {
	return Usecase{
		albumRepository,
		artistRepository,
		genreRepository,
		MusicRepository,
		playlistRepository,
		userRepository,
		encoder,
	}
}
