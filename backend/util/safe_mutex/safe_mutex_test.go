package safemutex

import (
	"fmt"
	"testing"
)

func TestMutex(t *testing.T) {
	tests := []struct {
		name string
		m    *Mutex
	}{
		{
			name: "multiple unlock",
			m:    NewMutex(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Lock()
			fmt.Printf("locked\n")

			tt.m.Unlock()
			fmt.Printf("unlocked\n")
			tt.m.Unlock()
			fmt.Printf("unlocked\n")

			tt.m.Lock()
			fmt.Printf("locked\n")
			tt.m.Unlock()
			fmt.Printf("unlocked\n")
		})
	}
}
