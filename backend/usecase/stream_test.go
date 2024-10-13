package usecase

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestUsecase_GetStreamM3U8(t *testing.T) {
	type args struct {
		ctx      context.Context
		musicID  uuid.UUID
		streamID string
	}
	tests := []struct {
		name    string
		u       *Usecase
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.u.GetStreamM3U8(tt.args.ctx, tt.args.musicID, tt.args.streamID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.GetStreamM3U8() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Usecase.GetStreamM3U8() = %v, want %v", got, tt.want)
			}
		})
	}
}
