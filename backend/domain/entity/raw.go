package entity

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/google/uuid"
	"github.com/walnuts1018/mucaron/backend/domain/entity/gormmodel"
)

type RawMusicMetadata struct {
	gormmodel.UUIDModel
	MusicID uuid.UUID

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
	Duration         time.Duration
	CreationDatetime synchro.Time[tz.AsiaTokyo]

	TagList []RawMusicMetadataTag
}

func (r RawMusicMetadata) ToEntity(Owner User) (music Music, album *Album, artist *Artist, genre *Genre) {
	if r.Artist != "" {
		artist = &Artist{
			Owner:    Owner,
			Name:     r.Artist,
			SortName: r.SortArtist,
		}
	}

	if r.Genre != "" {
		genre = &Genre{
			Owner: Owner,
			Name:  r.Genre,
		}
	}

	var musicTitle string
	if r.Title == "" {
		musicTitle = r.FileName
	} else {
		musicTitle = r.Title
	}

	music = Music{
		Owner:            Owner,
		Name:             musicTitle,
		SortName:         r.SortTitle,
		AlbumTrackNumber: int64(r.TrackNumber),
		Duration:         r.Duration,
		RawMetaData:      r,
		Status:           MetadataParsed,
	}

	if r.Album != "" {
		album = &Album{
			Owner:    Owner,
			Name:     r.Album,
			SortName: r.SortAlbum,
		}
	}

	return
}

type RawMusicMetadataTag struct {
	gormmodel.UUIDModel
	RawMusicMetadataID uuid.UUID

	Key   string
	Value string
}

func NewRawMusicMetadataTags(m map[string]any) []RawMusicMetadataTag {
	tags := []RawMusicMetadataTag{}
	for k, v := range m {
		tags = append(tags, RawMusicMetadataTag{
			Key:   k,
			Value: fmt.Sprintf("%v", v),
		})
	}
	slices.SortFunc(tags, func(a, b RawMusicMetadataTag) int {
		return strings.Compare(a.Key, b.Key)
	})
	return tags
}
