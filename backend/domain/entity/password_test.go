package entity

import (
	"reflect"
	"testing"
)

func TestNewLoginInfo(t *testing.T) {
	type args struct {
		rawPassword RawPassword
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				rawPassword: "password",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLoginInfo(tt.args.rawPassword)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLoginInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got.IsCorrectPassword(tt.args.rawPassword), true) {
				t.Errorf("NewLoginInfo() = %v", got)
			}
		})
	}
}
