package entity

import (
	"time"

	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
)

type RawMusicMetadata struct {
	FileName         string
	Title            string `json:"title"`
	SortTitle        string `json:"sort_name"`
	Artist           string `json:"artist"`
	SortArtist       string `json:"sort_artist"`
	AlbumArtist      string `json:"album_artist"`
	Composer         string `json:"composer"`
	Album            string `json:"album"`
	SortAlbum        string `json:"sort_album"`
	Genre            string `json:"genre"`
	TrackNumber      int64
	TrackTotal       int64
	DiscNumber       int64
	DiscTotal        int64
	CreationDatetime synchro.Time[tz.AsiaTokyo]
	Duration         time.Duration
}

func (r RawMusicMetadata) ToEntity() (Music, Artist, Album, Genre) {
	artist := Artist{
		Name:     r.Artist,
		SortName: r.SortArtist,
	}

	album := Album{
		Name:     r.Album,
		SortName: r.SortAlbum,
	}

	genre := Genre{
		Name: r.Genre,
	}

	m := Music{
		Name:             r.Title,
		SortName:         r.SortTitle,
		Album:            album,
		AlbumTrackNumber: int64(r.TrackNumber),
		Artists:          []Artist{artist},
		Duration:         r.Duration,
		Genre:            genre,
		RawMetaData:      r,
		Status:           MetadataParsed,
	}

	return m, artist, album, genre
}
