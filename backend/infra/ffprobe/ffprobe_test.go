package ffprobe

import (
	"reflect"
	"testing"
	"time"

	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/walnuts1018/mucaron/backend/domain/entity"
	"gopkg.in/vansante/go-ffprobe.v2"
)

func Test_newRawMusicTags(t *testing.T) {
	type args struct {
		data *ffprobe.ProbeData
	}
	tests := []struct {
		name    string
		args    args
		want    entity.RawMusicMetadata
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				data: &ffprobe.ProbeData{
					Format: &ffprobe.Format{
						Filename:        "filename",
						DurationSeconds: 124.273197,
						TagList: map[string]any{
							"title":         "title",
							"sort_name":     "sort_name",
							"artist":        "artist",
							"sort_artist":   "sort_artist",
							"album_artist":  "album_artist",
							"composer":      "composer",
							"album":         "album",
							"sort_album":    "sort_album",
							"genre":         "genre",
							"track":         "1/15",
							"disc":          "1/1",
							"creation_time": "2020-06-18T11:30:22.000000Z",
						},
					},
				},
			},
			want: entity.RawMusicMetadata{
				Title:            "title",
				SortTitle:        "sort_name",
				Artist:           "artist",
				SortArtist:       "sort_artist",
				AlbumArtist:      "album_artist",
				Composer:         "composer",
				Album:            "album",
				SortAlbum:        "sort_album",
				Genre:            "genre",
				TrackNumber:      1,
				TrackTotal:       15,
				DiscNumber:       1,
				DiscTotal:        1,
				CreationDatetime: synchro.In[tz.AsiaTokyo](time.Date(2020, 6, 18, 11, 30, 22, 0, time.UTC)),
				Duration:         124*time.Second + 273*time.Millisecond + 197*time.Microsecond,
				TagList: []entity.RawMusicMetadataTag{
					{
						Key:   "album",
						Value: "album",
					},
					{
						Key:   "album_artist",
						Value: "album_artist",
					},
					{
						Key:   "artist",
						Value: "artist",
					},
					{
						Key:   "composer",
						Value: "composer",
					},
					{
						Key:   "creation_time",
						Value: "2020-06-18T11:30:22.000000Z",
					},
					{
						Key:   "disc",
						Value: "1/1",
					},
					{
						Key:   "genre",
						Value: "genre",
					},
					{
						Key:   "sort_album",
						Value: "sort_album",
					},
					{
						Key:   "sort_artist",
						Value: "sort_artist",
					},
					{
						Key:   "sort_name",
						Value: "sort_name",
					},
					{
						Key:   "title",
						Value: "title",
					},
					{
						Key:   "track",
						Value: "1/15",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "no tags",
			args: args{
				data: &ffprobe.ProbeData{
					Format: &ffprobe.Format{
						Filename:        "filename",
						DurationSeconds: 124.273197,
						TagList:         map[string]any{},
					},
				},
			},
			want: entity.RawMusicMetadata{
				Title:            "",
				SortTitle:        "",
				Artist:           "",
				SortArtist:       "",
				AlbumArtist:      "",
				Composer:         "",
				Album:            "",
				SortAlbum:        "",
				Genre:            "",
				TrackNumber:      0,
				TrackTotal:       0,
				DiscNumber:       0,
				DiscTotal:        0,
				CreationDatetime: synchro.Time[tz.AsiaTokyo]{},
				Duration:         124*time.Second + 273*time.Millisecond + 197*time.Microsecond,
				TagList:          []entity.RawMusicMetadataTag{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newRawMusicMetadata(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("newRawMusicTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newRawMusicTags() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
