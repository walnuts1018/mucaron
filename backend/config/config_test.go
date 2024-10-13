package config

import (
	"log/slog"
	"os"
	"reflect"
	"testing"

	"dario.cat/mergo"
	_ "github.com/joho/godotenv/autoload"
)

var requiredEnvs = map[string]string{
	"MINIO_ACCESS_KEY": "test",
	"MINIO_SECRET_KEY": "test",
	"REDIS_PASSWORD":   "test",
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		envs    map[string]string //env
		want    Config
		wantErr bool
	}{
		{
			name: "check custom type default",
			envs: map[string]string{},
			want: Config{
				ServerPort: "8080",
				LogLevel:   slog.LevelInfo,
			},
			wantErr: false,
		},
		{
			name: "normal",
			envs: map[string]string{
				"SERVER_PORT":     "9000",
				"SERVER_ENDPOINT": "https://example.com",
			},
			want: Config{
				ServerPort:     "9000",
				ServerEndpoint: "https://example.com",
			},
			wantErr: false,
		},
		{
			name: "check custom type",
			envs: map[string]string{
				"LOG_LEVEL":      "debug",
				"SESSION_SECRET": "testtesttesttesttesttesttesttest",
			},
			want: Config{
				LogLevel:      slog.LevelDebug,
				SessionSecret: "testtesttesttesttesttesttesttest",
			},
			wantErr: false,
		},
		{
			name: "check PSQL",
			envs: map[string]string{
				"PSQL_HOST":     "host",
				"PSQL_PORT":     "15432",
				"PSQL_DATABASE": "db",
				"PSQL_USER":     "user",
				"PSQL_PASSWORD": "password",
				"PSQL_SSL_MODE": "sslmode",
				"PSQL_TIMEZONE": "timezone",
			},
			want: Config{
				PSQLDSN: "host=host port=15432 user=user password=password dbname=db sslmode=sslmode TimeZone=timezone",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()

			var envs = requiredEnvs
			for k, v := range tt.envs {
				envs[k] = v
			}

			for k, v := range envs {
				if err := os.Setenv(k, v); err != nil {
					t.Errorf("failed to set env: %v", err)
					return
				}
				defer os.Unsetenv(k)
			}

			got, err := Load()
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			ok, err := equal(got, tt.want)
			if err != nil {
				t.Errorf("failed to check config: %v", err)
				return
			}
			if !ok {
				t.Errorf("Load() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func equal(got, want Config) (bool, error) {
	merged := want
	if err := mergo.Merge(&merged, got); err != nil {
		return false, err
	}

	return reflect.DeepEqual(merged, got), nil
}

func Test_equal(t *testing.T) {
	type args struct {
		got  Config
		want Config
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				got: Config{
					ServerPort:     "8080",
					ServerEndpoint: "http://localhost:8080",
					LogLevel:       slog.LevelDebug,
				},
				want: Config{
					ServerPort: "8080",
					LogLevel:   slog.LevelDebug,
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "not equal",
			args: args{
				got: Config{
					ServerPort:     "8080",
					ServerEndpoint: "http://localhost:8080",
					LogLevel:       slog.LevelInfo,
				},
				want: Config{
					ServerPort: "9090",
					LogLevel:   slog.LevelDebug,
				},
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := equal(tt.args.got, tt.args.want)
			if (err != nil) != tt.wantErr {
				t.Errorf("equal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseSessionSecret(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "normal 32byte",
			args: args{
				v: "testtesttesttesttesttesttesttest",
			},
			want:    "testtesttesttesttesttesttesttest",
			wantErr: false,
		},
		{
			name: "random",
			args: args{
				v: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSessionSecret(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseSessionSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != "" {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("parseSessionSecret() = %v, want %v", got, tt.want)
				}
			} else {
				t.Logf("parseSessionSecret() = %v", got)
			}
		})
	}
}
