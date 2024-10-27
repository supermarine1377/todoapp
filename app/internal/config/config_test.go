package config

import (
	"testing"

	config_loader "github.com/supermarine1377/todoapp/app/internal/config/loader"
)

func TestConfig_Port(t *testing.T) {
	type fields struct {
		config config_loader.Config
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "1st",
			fields: fields{
				config: config_loader.Config{
					Port: 8080,
				},
			},
			want: 8080,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				config: tt.fields.config,
			}
			if got := c.Port(); got != tt.want {
				t.Errorf("Config.Port() = %v, want %v", got, tt.want)
			}
		})
	}
}
