package postgres

import (
	"github.com/google/uuid"
	"github.com/walnuts1018/mucaron/domain/entity"
)

type music struct {
	ID               uuid.UUID `db:"id"`
	Title            string    `db:"title"`
	AlbumID          uuid.UUID `db:"album_id"`
	AlbumTruckNumber int64     `db:"album_truck_number"`
	IsCover          bool      `db:"is_cover"`
	OriginalMusicID  uuid.UUID `db:"original_music_id"`
	Score            int64     `db:"score"`
}

type musics_and_artists struct {
	id       int64     `db:"id"`
	MusicID  uuid.UUID `db:"music_id"`
	ArtistID uuid.UUID `db:"artist_id"`
}

func (m *music) ToEntity() entity.Music {
	entity := entity.Music{
		ID:               m.ID,
		Title:            m.Title,
		AlbumID:          m.AlbumID,
		AlbumTruckNumber: m.AlbumTruckNumber,
		CoverInfo: entity.CoverInfo{
			IsCover:         m.IsCover,
			OriginalMusicID: m.OriginalMusicID,
		},
		Score: m.Score,
	}
	return entity
}

func (m *music) FromEntity(e entity.Music) {
	m.ID = e.ID
	m.Title = e.Title
	m.AlbumID = e.AlbumID
	m.AlbumTruckNumber = e.AlbumTruckNumber
	m.IsCover = e.CoverInfo.IsCover
	m.OriginalMusicID = e.CoverInfo.OriginalMusicID
	m.Score = e.Score
}

func (p *PostgresClient) initMusic() error {
	_, err := p.db.Exec(`CREATE TABLE IF NOT EXISTS musics (
		id UUID PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		album_id INT NOT NULL,
		album_truck_number INT NOT NULL,
		is_cover BOOLEAN NOT NULL,
		original_music_id INT NOT NULL,
		score INT NOT NULL
	)`)
	return err
}

func (p *PostgresClient) initMusicsAndArtists() error {
	_, err := p.db.Exec(`CREATE TABLE IF NOT EXISTS musics_and_artists (
		id SERIAL PRIMARY KEY,
		music_id INT NOT NULL,
		artist_id INT NOT NULL,
		FOREIGN KEY (music_id) REFERENCES musics(id)
	)`)
	return err
}

func (p *PostgresClient) CreateMusic(e entity.Music) error {
	var m music
	m.FromEntity(e)

	if _, err := p.db.Exec(`INSERT INTO musics_and_Artists (
			music_id, artist_id
		) VALUES (
			$1, $2
		)`,
		e.ArtistIDs); err != nil {
		return err
	}

	if _, err := p.db.NamedExec(`INSERT INTO musics (
			id, title, album_id, album_truck_number, is_cover, original_music_id, score
		) VALUES (
			:id, :title, :album_id, :album_truck_number, :is_cover, :original_music_id, :score
		)`,
		m); err != nil {
		return err
	}

	return nil
}
