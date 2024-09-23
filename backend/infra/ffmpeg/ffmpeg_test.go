package ffmpeg

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

const testFilesDir = "./test_files"

var testfiles = map[string]map[string]string{
	"video": {
		"avi":  "file_example_AVI_480_750kB.avi",
		"mov":  "file_example_MOV_480_700kB.mov",
		"mp4":  "file_example_MP4_480_1_5MG.mp4",
		"webm": "file_example_WEBM_480_900KB.webm",
	},
	"audio": {
		"mp3": "file_example_MP3_700KB.mp3",
		"wav": "file_example_WAV_1MG.wav",
	},
}

func init() {
	slog.SetDefault(
		slog.New(
			slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelDebug},
			),
		),
	)
}

func TestFFMPEG_CreateArgs(t *testing.T) {
	f := &FFMPEG{
		BaseURL:    &url.URL{Scheme: "https", Host: "localhost", Path: "mucaron-test"},
		FPS:        30,
		Preset:     Veryslow,
		VideoCodec: "libx264",
		VideoQualityKeys: []VideoQualityKey{
			VideoQualityKey360P,
			VideoQualityKey720P,
			VideoQualityKey1080P,
		}}

	type args struct {
		id            uuid.UUID
		inputFileName string
		audioOnly     bool
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "audioOnly",
			args: args{
				id:            uuid.MustParse("0191456e-7c30-76db-a572-b646efcf2e6f"),
				inputFileName: "input.mp4",
				audioOnly:     true,
			},
			want: []string{
				"-i", "input.mp4",
				"-y",
				"-hide_banner",
				"-preset", "veryslow",
				"-keyint_min", "100",
				"-g", "100",
				"-sc_threshold", "0",
				"-r", "30",
				"-c:v", "libx264",
				"-pix_fmt", "yuv420p",
				"-map", "a:0",
				"-c:a", "copy",
				"-f", "hls",
				"-hls_time", "4",
				"-hls_playlist_type", "vod",
				"-hls_flags", "independent_segments",
				"-hls_base_url", "https://localhost/mucaron-test/",
				"-master_pl_name", "primary.m3u8",
				"-hls_segment_filename", "0191456e-7c30-76db-a572-b646efcf2e6f/stream_%v/s%06d.ts",
				"-strftime_mkdir", "1",
				"-var_stream_map", "a:0",
				"0191456e-7c30-76db-a572-b646efcf2e6f/stream_%v.m3u8",
			},
		},
		{
			name: "with video",
			args: args{
				id:            uuid.MustParse("0191456e-7c30-76db-a572-b646efcf2e6f"),
				inputFileName: "input.mp4",
				audioOnly:     false,
			},
			want: []string{
				"-i", "input.mp4",
				"-y",
				"-hide_banner",
				"-preset", "veryslow",
				"-keyint_min", "100",
				"-g", "100",
				"-sc_threshold", "0",
				"-r", "30",
				"-c:v", "libx264",
				"-pix_fmt", "yuv420p",

				// 360p
				"-map", "v:0",
				"-vf:0", "scale=-1:360",
				"-b:v:0", "365k",
				"-maxrate:0", "390k",
				"-bufsize:0", "640k",

				// 720p
				"-map", "v:0",
				"-vf:1", "scale=-1:720",
				"-b:v:1", "4.5M",
				"-maxrate:1", "4.8M",
				"-bufsize:1", "8M",

				// 1080p
				"-map", "v:0",
				"-vf:2", "scale=-1:1080",
				"-b:v:2", "7.8M",
				"-maxrate:2", "8.3M",
				"-bufsize:2", "14M",

				"-map", "a:0",
				"-map", "a:0",
				"-map", "a:0",
				"-map", "a:0",

				"-c:a", "copy",
				"-f", "hls",
				"-hls_time", "4",
				"-hls_playlist_type", "vod",
				"-hls_flags", "independent_segments",
				"-hls_base_url", "https://localhost/mucaron-test/",
				"-master_pl_name", "primary.m3u8",
				"-hls_segment_filename", "0191456e-7c30-76db-a572-b646efcf2e6f/stream_%v/s%06d.ts",
				"-strftime_mkdir", "1",
				"-var_stream_map", "v:0,a:0 v:1,a:1 v:2,a:2 a:3",
				"0191456e-7c30-76db-a572-b646efcf2e6f/stream_%v.m3u8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := f.createArgs(tt.args.id, tt.args.inputFileName, tt.args.audioOnly)
			if (err != nil) != tt.wantErr {
				t.Errorf("FFMPEG.CreateArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FFMPEG.CreateArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFFMPEG_Encode(t *testing.T) {
	f := &FFMPEG{
		BaseURL:    &url.URL{Scheme: "https", Host: "localhost", Path: "mucaron-test"},
		FPS:        30,
		Preset:     Ultrafast,
		VideoCodec: "libx264",
		VideoQualityKeys: []VideoQualityKey{
			VideoQualityKey360P,
		}}

	type args struct {
		id        uuid.UUID
		path      string
		audioOnly bool
	}
	type test struct {
		name    string
		args    args
		wantErr bool
	}

	tests := make([]test, 0)

	// video
	for k, v := range testfiles["video"] {
		id, err := uuid.NewV7()
		if err != nil {
			t.Errorf("failed to gen id: %s", err)
		}
		tests = append(tests, test{
			name: k,
			args: args{
				id:        id,
				path:      path.Join(testFilesDir, v),
				audioOnly: false,
			},
			wantErr: false,
		})
	}

	// audio
	for k, v := range testfiles["audio"] {
		id, err := uuid.NewV7()
		if err != nil {
			t.Errorf("failed to gen id: %s", err)
		}
		tests = append(tests, test{
			name: k,
			args: args{
				id:        id,
				path:      path.Join(testFilesDir, v),
				audioOnly: true,
			},
			wantErr: false,
		})
	}

	workdir, err := os.Getwd()
	if err != nil {
		t.Errorf("failed to get workdir: %v", err)
		return
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hlsDir, err := f.Encode(tt.args.id, filepath.Join(workdir, tt.args.path), tt.args.audioOnly)
			if (err != nil) != tt.wantErr {
				t.Errorf("FFMPEG.Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(hlsDir)
		})
	}
}
