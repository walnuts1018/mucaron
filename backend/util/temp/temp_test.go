package temp

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"
)

func TestCreateTempFile(t *testing.T) {
	type args struct {
		r        io.Reader
		filename string
	}
	tests := []struct {
		name        string
		args        args
		wantPath    string
		wantContent []byte
		wantErr     bool
	}{
		{
			name: "normal",
			args: args{
				r:        strings.NewReader("test"),
				filename: "TestCreateTempFile",
			},
			wantContent: []byte("test"),
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateTempFile(tt.args.r, tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTempFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			defer func() {
				got.Close()
				os.Remove(got.UseFile().Name())
			}()

			gotContent, err := io.ReadAll(got.UseFile())
			if err != nil {
				t.Errorf("failed to read all: %v", err)
				return
			}
			if !reflect.DeepEqual(gotContent, tt.wantContent) {
				t.Errorf("CreateTempFile() = %v, want %v", string(gotContent), string(tt.wantContent))
				return
			}
		})
	}
}

func TestTempFile_Close(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		tmpfile, err := CreateTempFile(strings.NewReader("test"), "TestTempFile_Close")
		if err != nil {
			t.Errorf("failed to create temp file: %v", err)
			return
		}
		filename := tmpfile.UseFile().Name()

		wg := sync.WaitGroup{}

		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(tempFile *TempFile) {
				file := tempFile.UseFile()
				fmt.Printf("using TempFile, id: %v", i)
				defer func() {
					fmt.Printf("closing TempFile, id: %v", i)
					if err := tempFile.Close(); err != nil {
						t.Errorf("failed to close: %v", err)
					}
					wg.Done()
				}()

				content, err := os.ReadFile(file.Name())
				if err != nil {
					t.Errorf("failed to read file: %v", err)
					return
				}

				if string(content) != "test" {
					t.Errorf("unexpected content: %v", string(content))
					return
				}
			}(tmpfile)
		}

		wg.Wait()
		tmpfile.Close()

		closed, count := tmpfile.checkClosed()
		if !closed {
			t.Errorf("Not closed, count: %v", count)
			return
		}

		if _, err := os.Stat(filename); err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("ok, file is removed")
			} else {
				t.Errorf("failed to stat: %v", err)
			}
		}
	})
}
