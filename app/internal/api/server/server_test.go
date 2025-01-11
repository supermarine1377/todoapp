// package server はサーバーを提供する
package server

import (
	"testing"
)

func TestWithPort(t *testing.T) {
	type args struct {
		port int
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "Invalid port: negative number",
			args:    args{port: -1},
			wantErr: ErrInvalidPort,
		},
		{
			name:    "Invalid port: zero",
			args:    args{port: 0},
			wantErr: ErrInvalidPort,
		},
		{
			name:    "Valid port",
			args:    args{port: 8080},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithPort(tt.args.port)
			opts := &options{}
			err := opt(opts)

			if err != tt.wantErr {
				t.Errorf("WithPort() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil && *opts.port != tt.args.port {
				t.Errorf("WithPort() port = %v, want %v", *opts.port, tt.args.port)
			}
		})
	}
}

func TestWithDSN(t *testing.T) {
	type args struct {
		dsn string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "Empty DSN",
			args:    args{dsn: ""},
			wantErr: ErrInvalidDSN,
		},
		{
			name:    "Valid PostgreSQL DSN",
			args:    args{dsn: "postgres://user:pass@localhost:5432/dbname?sslmode=disable"},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithDSN(tt.args.dsn)
			opts := &options{}
			err := opt(opts)

			if err != tt.wantErr {
				t.Errorf("WithDSN() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil && *opts.dsn != tt.args.dsn {
				t.Errorf("WithDSN() dsn = %v, want %v", opts.dsn, tt.args.dsn)
			}
		})
	}
}
