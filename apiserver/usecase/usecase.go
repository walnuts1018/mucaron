package usecase

type AlbumRepository interface{}
type ArtistRepository interface{}
type GenreRepository interface{}
type MusicRepository interface{}
type PlaylistRepository interface{}
type UserRepository interface{}

type Usecase struct {
	albumRepository    AlbumRepository
	artistRepository   ArtistRepository
	genreRepository    GenreRepository
	MusicRepository    MusicRepository
	playlistRepository PlaylistRepository
	userRepository     UserRepository
}

func NewUsecase(
	albumRepository AlbumRepository,
	artistRepository ArtistRepository,
	genreRepository GenreRepository,
	MusicRepository MusicRepository,
	playlistRepository PlaylistRepository,
	userRepository UserRepository,
) Usecase {
	return Usecase{
		albumRepository,
		artistRepository,
		genreRepository,
		MusicRepository,
		playlistRepository,
		userRepository,
	}
}
