package vault

import (
	"github.com/hashicorp/vault/api"
	"testing"
)

func TestConfig_Start(t *testing.T) {
	type fields struct {
		vaultConfig *api.Config
	}
	tests := []struct {
		name   string
		fields fields
		want   error
	}{
		{
			name:   "create vault client",
			fields: fields{vaultConfig: api.DefaultConfig()},
			want:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewConfig()

			got, err := vc.Start()

			if err != tt.want {
				t.Errorf(" got %v want %v", got, tt.want)
			}
		})
	}
}
