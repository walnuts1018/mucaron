package gormmodel

import "testing"

func TestUUIDModel_CreateID(t *testing.T) {
	tests := []struct {
		name    string
		u       *UUIDModel
		wantErr bool
	}{
		{
			name: "normal",
			u:    &UUIDModel{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.u.CreateID(); (err != nil) != tt.wantErr {
				t.Errorf("UUIDModel.CreateID() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.u.ID.String() == "" {
				t.Errorf("UUIDModel.CreateID() = %v", tt.u.ID)
			}
		})
	}
}
