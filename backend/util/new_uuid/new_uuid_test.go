package newuuid

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestNewV7(t *testing.T) {
	tests := []struct {
		name      string
		testValue uuid.UUID
		wantErr   bool
	}{
		{
			name:      "normal",
			testValue: uuid.Nil,
			wantErr:   false,
		},
		{
			name:      "testValue",
			testValue: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.testValue != uuid.Nil {
				SetUUIDValue(tt.testValue)
			}

			got, err := NewV7()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewV7() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.DeepEqual(got, uuid.Nil) {
				t.Errorf("NewV7() = %v", got)
			}

			if tt.testValue != uuid.Nil {
				if !reflect.DeepEqual(got, tt.testValue) {
					t.Errorf("NewV7() = %v, want %v", got, tt.testValue)
				}
			}
		})
	}
}
