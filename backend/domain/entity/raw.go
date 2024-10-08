package entity

import (
	"time"

	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/walnuts1018/mucaron/backend/domain/entity/gormmodel"
)

type RawMusicMetadata struct {
	gormmodel.UUIDModel
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

	TagList map[string]string
}

func (r RawMusicMetadata) ToEntity(Owner User) Music {
	artist := Artist{
		Owner:    Owner,
		Name:     r.Artist,
		SortName: r.SortArtist,
	}

	album := Album{
		Owner:    Owner,
		Name:     r.Album,
		SortName: r.SortAlbum,
	}

	genre := Genre{
		Owner: Owner,
		Name:  r.Genre,
	}

	var musicTitle string
	if r.Title == "" {
		musicTitle = r.FileName
	} else {
		musicTitle = r.Title
	}

	music := Music{
		Name:             musicTitle,
		SortName:         r.SortTitle,
		Album:            album,
		AlbumTrackNumber: int64(r.TrackNumber),
		Artists:          []Artist{artist},
		Duration:         r.Duration,
		Genre:            genre,
		RawMetaData:      r,
		Status:           MetadataParsed,
	}

	return music
}
