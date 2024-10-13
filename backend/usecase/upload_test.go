package usecase

import "testing"

func Test_replaceM3U8URL(t *testing.T) {
	type args struct {
		content        string
		serverEndpoint string
		musicID        string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				content: `#EXTM3U
#EXT-X-VERSION:6
#EXT-X-STREAM-INF:BANDWIDTH=680836,RESOLUTION=1080x1080,CODECS="avc1.640032,mp4a.40.2"
stream_0.m3u8

#EXT-X-STREAM-INF:BANDWIDTH=5229336,RESOLUTION=1080x1080,CODECS="avc1.640032,mp4a.40.2"
stream_1.m3u8

#EXT-X-STREAM-INF:BANDWIDTH=8859336,RESOLUTION=1080x1080,CODECS="avc1.640032,mp4a.40.2"
stream_2.m3u8

#EXT-X-STREAM-INF:BANDWIDTH=279336,CODECS="mp4a.40.2"
stream_3.m3u8
`,
				serverEndpoint: "http://localhost:8080",
				musicID:        "music_id",
			},
			want: `#EXTM3U
#EXT-X-VERSION:6
#EXT-X-STREAM-INF:BANDWIDTH=680836,RESOLUTION=1080x1080,CODECS="avc1.640032,mp4a.40.2"
http://localhost:8080/api/v1/music/music_id/stream/stream_0.m3u8

#EXT-X-STREAM-INF:BANDWIDTH=5229336,RESOLUTION=1080x1080,CODECS="avc1.640032,mp4a.40.2"
http://localhost:8080/api/v1/music/music_id/stream/stream_1.m3u8

#EXT-X-STREAM-INF:BANDWIDTH=8859336,RESOLUTION=1080x1080,CODECS="avc1.640032,mp4a.40.2"
http://localhost:8080/api/v1/music/music_id/stream/stream_2.m3u8

#EXT-X-STREAM-INF:BANDWIDTH=279336,CODECS="mp4a.40.2"
http://localhost:8080/api/v1/music/music_id/stream/stream_3.m3u8
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := replaceM3U8URL(tt.args.content, tt.args.serverEndpoint, tt.args.musicID)
			if (err != nil) != tt.wantErr {
				t.Errorf("replaceM3U8URL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("replaceM3U8URL() = %v, want %v", got, tt.want)
			}
		})
	}
}
