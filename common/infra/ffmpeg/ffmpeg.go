package ffmpeg

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/walnuts1018/mucaron/common/config"
)

type FFMPEG struct {
	BaseURL          *url.URL
	FPS              int
	Preset           Preset
	VideoCodec       string
	VideoQualityKeys []VideoQualityKey
}

func NewFFMPEG(cfg config.Config) (*FFMPEG, error) {
	url, err := url.Parse(cfg.MinIOPublicBaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %v", err)
	}
	url = url.JoinPath(cfg.MinIOBucket)

	return &FFMPEG{
		BaseURL:    url,
		FPS:        30,
		Preset:     Veryslow,
		VideoCodec: "libx264",
		VideoQualityKeys: []VideoQualityKey{
			VideoQualityKey360P,
			VideoQualityKey720P,
			VideoQualityKey1080P,
		},
	}, nil
}

func (f *FFMPEG) CreateArgs(id uuid.UUID, inputFileName string, audioOnly bool) ([]string, error) {
	args := []string{
		"-i", inputFileName,
		"-y",
		"-hide_banner",
		"-preset", string(f.Preset),
		"-keyint_min", "100",
		"-g", "100",
		"-sc_threshold", "0",
		"-r", fmt.Sprintf("%d", f.FPS),
		"-c:v", f.VideoCodec,
		"-pix_fmt", "yuv420p",
	}

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
		"-hls_base_url", fmt.Sprintf("%s/", f.BaseURL.String()),
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

func (f *FFMPEG) Encode(id uuid.UUID, movie io.Reader, audioOnly bool) (string, error) {
	inputFile, err := os.CreateTemp("", fmt.Sprintf("mucaron-input-%s", id))
	if err != nil {
		return "", err
	}
	defer inputFile.Close()
	_, err = io.Copy(inputFile, movie)
	if err != nil {
		return "", err
	}

	command := "ffmpeg"

	args, err := f.CreateArgs(id, inputFile.Name(), audioOnly)
	if err != nil {
		return "", err
	}
	slog.Debug("args created", slog.Any(
		"args",
		args,
	))

	cmd := exec.Command(command, args...)
	cmd.Dir = filepath.Dir(inputFile.Name())

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		slog.Error("ffmpeg error",
			slog.String("stdout", stdout.String()),
			slog.String("stderr", stderr.String()),
		)
		return "", fmt.Errorf("failed to run ffmpeg: %w", err)
	}

	return path.Join(cmd.Dir, id.String()), nil
}
