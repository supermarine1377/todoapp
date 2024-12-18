// package loader は、サーバーの設定を読み込む機能を提供する
package loader_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supermarine1377/todoapp/app/internal/config/loader"
)

type envVar struct {
	Port string
	DSN  string
}

func init() {
	os.Clearenv()
}

func (ev *envVar) Set(t *testing.T) {
	if ev.Port != "" {
		t.Setenv("PORT", ev.Port)
	}
	if ev.DSN != "" {
		t.Setenv("DATABASE_DSN", ev.DSN)
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		envVar  envVar
		want    *loader.Config
		wantErr bool
	}{
		{
			name: "When PORT environment variable is set",
			envVar: envVar{
				Port: "8080",
				DSN:  "path",
			},
			want: &loader.Config{
				Port: 8080,
				DB: loader.DB{
					UseSQLite: true,
					DSN:       "path",
				},
			},
			wantErr: false,
		},
		{
			name: "When PORT environment variable is not set",
			want: &loader.Config{
				Port: 8080,
				DB: loader.DB{
					UseSQLite: true,
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid environment variable",
			envVar: envVar{
				Port: "a",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.envVar.Set(t)

			got, err := loader.Parse()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
