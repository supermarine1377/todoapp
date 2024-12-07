package config

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supermarine1377/todoapp/app/internal/config/loader"
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

func TestNew(t *testing.T) {
	type args struct {
		clf ConfigLoaderFunc
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "When config.ConfigLoaderFunc returns an error",
			args: args{
				func() (*loader.Config, error) {
					return nil, errors.New("Some error")
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "When config.ConfigLoaderFunc returns no error",
			args: args{
				func() (*loader.Config, error) {
					return &loader.Config{
						Port: 8080,
					}, nil
				},
			},
			want: &Config{
				config: loader.Config{
					Port: 8080,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.clf)
			assert.Equal(t, err != nil, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
