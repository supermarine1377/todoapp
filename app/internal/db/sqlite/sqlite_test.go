// package sqlite は、データベースとしてSQLiteを使う場合の機能を提供する
package sqlite_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/supermarine1377/todoapp/app/internal/db/sqlite"
)

type mockConfig struct {
	dsn string
}

func (mc mockConfig) DSN() string {
	return mc.dsn
}

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		setupFile func(tempDir string) (path string, err error)
		wantErr   bool
		err       error
	}{
		{
			name: "Valid file with read/write permissions",
			setupFile: func(tempDir string) (path string, err error) {
				f, err := os.CreateTemp(tempDir, "")
				if err != nil {
					return "", err
				}
				return f.Name(), nil
			},
			wantErr: false,
		},
		{
			name: "File does not exist",
			setupFile: func(_ string) (path string, err error) {
				dummy := "dummy"
				_ = os.Remove(dummy)
				return dummy, nil
			},
			wantErr: true,
			err:     sqlite.ErrNoSuchFileOrDir,
		},
		{
			name: "File without write permission",
			setupFile: func(tempDir string) (path string, err error) {
				f, err := os.CreateTemp(tempDir, "")
				if err != nil {
					return "", err
				}
				_ = f.Chmod(0400)
				return f.Name(), nil
			},
			wantErr: true,
			err:     sqlite.ErrFileNotAccessble,
		},
	}
	for _, tt := range tests {
		temp, err := os.MkdirTemp(".", "temp_")
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			_ = os.RemoveAll(temp)
		})

		t.Run(tt.name, func(t *testing.T) {
			path, err := tt.setupFile(temp)
			if err != nil {
				t.Fatal(err)
			}
			mc := mockConfig{dsn: path}
			_, err = sqlite.New(mc)
			if !tt.wantErr {
				require.NoError(t, err)
			}
			if tt.wantErr {
				require.Error(t, err)
				if tt.err != nil {
					require.ErrorIs(t, err, tt.err)
				}
			}
		})
	}
}
