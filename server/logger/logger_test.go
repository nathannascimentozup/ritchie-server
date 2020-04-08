package logger

import (
	"testing"
)

func TestLoadLogDefinition(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "run",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LoadLogDefinition()
		})
	}
}
