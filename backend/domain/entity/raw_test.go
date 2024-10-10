package entity

import (
	"strings"
	"testing"

	sliceequal "github.com/walnuts1018/mucaron/backend/util/slice_equal"
)

func TestNewRawMusicMetadataTags(t *testing.T) {
	type args struct {
		m map[string]any
	}
	tests := []struct {
		name string
		args args
		want []RawMusicMetadataTag
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
			want: []RawMusicMetadataTag{
				{
					Key:   "string",
					Value: "string",
				},
				{
					Key:   "int",
					Value: "100",
				},
				{
					Key:   "float",
					Value: "3.14159265358979",
				},
				{
					Key:   "bool",
					Value: "true",
				},
				{
					Key:   "slice",
					Value: "[slice]",
				},
				{
					Key:   "map",
					Value: "map[key:value]",
				},
				{
					Key:   "byte",
					Value: "[98 121 116 101]",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRawMusicMetadataTags(tt.args.m); !sliceequal.Equal(got, tt.want, func(a, b RawMusicMetadataTag) int {
				return strings.Compare(a.Key, b.Key)
			}) {
				t.Errorf("NewRawMusicMetadataTags() = %v, want %v", got, tt.want)
			}
		})
	}
}
