package ffprobe

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/pkg/errors"
	"github.com/walnuts1018/mucaron/backend/domain/entity"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type FFProbe struct{}

func NewFFProbe() FFProbe {
	return FFProbe{}
}

func (p FFProbe) GetMetadata(ctx context.Context, path string) (entity.RawMusicMetadata, error) {
	data, err := ffprobe.ProbeURL(ctx, path)
	if err != nil {
		return entity.RawMusicMetadata{}, err
	}
	return newRawMusicMetadata(data)
}

func newRawMusicMetadata(data *ffprobe.ProbeData) (entity.RawMusicMetadata, error) {
	var m entity.RawMusicMetadata
	m.Title = toString(data.Format.TagList, "title")
	m.SortTitle = toString(data.Format.TagList, "sort_name")
	m.Artist = toString(data.Format.TagList, "artist")
	m.SortArtist = toString(data.Format.TagList, "sort_artist")
	m.AlbumArtist = toString(data.Format.TagList, "album_artist")
	m.Composer = toString(data.Format.TagList, "composer")
	m.Album = toString(data.Format.TagList, "album")
	m.SortAlbum = toString(data.Format.TagList, "sort_album")
	m.Genre = toString(data.Format.TagList, "genre")

	track := toString(data.Format.TagList, "track")
	if track != "" {
		s := strings.Split(track, "/")
		if len(s) == 2 {
			var err error
			m.TrackNumber, err = strconv.ParseInt(s[0], 10, 64)
			if err != nil {
				m.TrackNumber = 0
			}
			m.TrackTotal, err = strconv.ParseInt(s[1], 10, 64)
			if err != nil {
				m.TrackTotal = 0
			}
		}
	}

	disc := toString(data.Format.TagList, "disc")
	if disc != "" {
		s := strings.Split(disc, "/")
		if len(s) == 2 {
			var err error
			m.DiscNumber, err = strconv.ParseInt(s[0], 10, 64)
			if err != nil {
				m.DiscNumber = 0
			}
			m.DiscTotal, err = strconv.ParseInt(s[1], 10, 64)
			if err != nil {
				m.DiscTotal = 0
			}
		}
	}

	var err error
	m.CreationDatetime, err = synchro.ParseISO[tz.AsiaTokyo](toString(data.Format.TagList, "creation_time"))
	if err != nil {
		m.CreationDatetime = synchro.Time[tz.AsiaTokyo]{}
	}

	duration, err := time.ParseDuration(fmt.Sprintf("%fs", data.Format.DurationSeconds))
	if err != nil {
		return entity.RawMusicMetadata{}, errors.Wrap(err, "failed to parse duration")
	}

	m.Duration = duration

	m.TagList = entity.NewRawMusicMetadataTags(data.Format.TagList)

	return m, nil
}

func toString(m map[string]any, key string) string {
	v, ok := m[key]
	if !ok {
		return ""
	}
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}
