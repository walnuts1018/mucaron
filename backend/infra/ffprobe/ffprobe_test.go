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
				FileName:         "filename",
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
				TagList: map[string]string{
					"album":         "album",
					"album_artist":  "album_artist",
					"artist":        "artist",
					"composer":      "composer",
					"creation_time": "2020-06-18T11:30:22.000000Z",
					"disc":          "1/1",
					"genre":         "genre",
					"sort_album":    "sort_album",
					"sort_artist":   "sort_artist",
					"sort_name":     "sort_name",
					"title":         "title",
					"track":         "1/15",
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
				FileName:         "filename",
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
				TagList:          map[string]string{},
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

func Test_toStringMap(t *testing.T) {
	type args struct {
		m map[string]any
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "test",
			args: args{
				m: map[string]any{
					"string": "string",
					"int":    int(100),
					"float":  float64(3.14159265358979),
					"bool":   true,
					"slice":  []string{"slice"},
					"map":    map[string]string{"key": "value"},
					"byte":   []byte("byte"),
				},
			},
			want: map[string]string{
				"string": "string",
				"int":    "100",
				"float":  "3.14159265358979",
				"bool":   "true",
				"slice":  "[slice]",
				"map":    "map[key:value]",
				"byte":   "[98 121 116 101]",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toStringMap(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toStringMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
