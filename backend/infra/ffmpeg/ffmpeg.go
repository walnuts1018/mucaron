package ffmpeg

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/walnuts1018/mucaron/backend/config"
	"github.com/walnuts1018/mucaron/backend/util/fileutil"
)

var baseURL = url.URL{Scheme: "https", Host: "to-be-replaced.example.com"}
var OutDirPrefix = "mucaron-outdir"

type FFMPEG struct {
	baseURL          *url.URL
	FPS              int
	Preset           Preset
	VideoCodec       string
	VideoQualityKeys []VideoQualityKey
	logFileDir       string
}

func NewFFMPEG(cfg config.Config) (*FFMPEG, error) {
	return &FFMPEG{
		baseURL:    &baseURL,
		FPS:        30,
		Preset:     Medium,
		VideoCodec: "libx264",
		VideoQualityKeys: []VideoQualityKey{
			VideoQualityKey360P,
			VideoQualityKey720P,
			VideoQualityKey1080P,
		},
		logFileDir: filepath.Join(cfg.LogDir, "ffmpeg"),
	}, nil
}

func (f *FFMPEG) GetOutDirPrefix() string {
	return OutDirPrefix
}

func (f *FFMPEG) createArgs(id string, inputFileName string, audioOnly bool) ([]string, error) {
	args := []string{
		"-i", inputFileName,
		"-y",
		"-hide_banner",
		"-progress",
		"-preset", string(f.Preset),
		"-keyint_min", "100",
		"-g", "100",
		"-sc_threshold", "0",
		"-r", fmt.Sprintf("%d", f.FPS),
		"-c:v", f.VideoCodec,
		"-pix_fmt", "yuv420p",
	}

// ffmpeg -i test_files/file_example_MP4_480_1_5MG.mp4 -y -hide_banner -progress - -preset medium -keyint_min 100 -g 100 -sc_threshold 0 -r 30 -c:v libx264 -c:a copy -pix_fmt yuv420p \
// -map "v:0?" -filter:v:0 scale=-2:360 -b:v:0 365k -maxrate:0 390k -bufsize:0 640k \
// -map "v:0?" -filter:v:1 scale=-2:720 -b:v:1 4.5M -maxrate:1 4.8M -bufsize:1 8M \
// -map "v:0?" -filter:v:2 scale=-2:1080 -b:v:2 7.8M -maxrate:2 8.3M -bufsize:2 14M \
// -map 0:a \
// -init_seg_name init\$RepresentationID\$.\$ext\$ -media_seg_name chunk\$RepresentationID\$-\$Number%05d\$.\$ext\$ \
// -use_template 1 -use_timeline 1  \
// -seg_duration 4 -adaptation_sets "id=0,streams=a id=1,streams=v" \
// -f dash Dash/dash.mpd


	if !audioOnly {
		for i, quality := range f.VideoQualityKeys {
			args = append(args,
				"-map", "v:0",
				fmt.Sprintf("-vf:%d", i), fmt.Sprintf("scale=%s", videoQualities[quality].Scale),
				fmt.Sprintf("-b:v:%d", i), videoQualities[quality].Bitrate,
				fmt.Sprintf("-maxrate:%d", i), videoQualities[quality].MaxBitrate,
				fmt.Sprintf("-bufsize:%d", i), videoQualities[quality].Bufsize,
			)
		}
	}

	args = append(args, "-map", "a:0")
	if !audioOnly {
		for i := 0; i < len(f.VideoQualityKeys); i++ {
			args = append(args, "-map", "a:0")
		}
	}

	args = append(args,
		"-c:a", "copy",
		"-f", "hls",
		"-hls_time", "4",
		"-hls_playlist_type", "vod",
		"-hls_flags", "independent_segments",
		"-hls_base_url", fmt.Sprintf("%s/", f.baseURL.String()),
		"-master_pl_name", "primary.m3u8",
		"-hls_segment_filename", fmt.Sprintf("%s/stream_%%v/s%%06d.ts", id),
		"-strftime_mkdir", "1",
	)

	var var_stream_maps []string
	if audioOnly {
		var_stream_maps = append(var_stream_maps, "a:0")
	} else {
		for i := 0; i < len(f.VideoQualityKeys); i++ {
			var_stream_maps = append(var_stream_maps, fmt.Sprintf("v:%d,a:%d", i, i))
		}
		var_stream_maps = append(var_stream_maps, fmt.Sprintf("a:%d", len(f.VideoQualityKeys)))
	}

	args = append(args, "-var_stream_map", strings.Join(var_stream_maps, " "))
	args = append(args, fmt.Sprintf("%s/stream_%%v.m3u8", id))
	return args, nil
}

/*
動画ファイルを受け取って、HLSとしてエンコードします
戻り値はエンコード先ディレクトリの絶対パスです。
*/
func (f *FFMPEG) Encode(id string, path string, audioOnly bool) (string, error) {
	workdir, err := os.MkdirTemp("", OutDirPrefix)
	if err != nil {
		return "", err
	}

	args, err := f.createArgs(id, path, audioOnly)
	if err != nil {
		return "", err
	}
	slog.Debug("ffmpeg args", slog.Any("args", args))

	cmd := exec.Command("ffmpeg", args...)
	cmd.Dir = workdir

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	logfile, err := fileutil.CreateFileRecursive(filepath.Join(f.logFileDir, id+".log"))
	if err != nil {
		return "", fmt.Errorf("failed to create log file: %w", err)
	}

	cmd.Stdout = io.MultiWriter(logfile, &stdout)
	cmd.Stderr = io.MultiWriter(logfile, &stderr)

	if err := cmd.Run(); err != nil {
		slog.Error("ffmpeg error",
			slog.String("stdout", stdout.String()),
			slog.String("stderr", stderr.String()),
		)
		return "", fmt.Errorf("failed to run ffmpeg: %w", err)
	}

	return filepath.Join(workdir, id), nil
}
