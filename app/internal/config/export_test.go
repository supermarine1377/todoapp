package config

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supermarine1377/todoapp/app/internal/config/loader"
)

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
