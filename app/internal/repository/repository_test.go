package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/supermarine1377/todoapp/app/internal/model/entity/task"
	"github.com/supermarine1377/todoapp/app/internal/repository"
	"github.com/supermarine1377/todoapp/app/internal/repository/mock"
	"go.uber.org/mock/gomock"
)

func TestTaskRepository_ListCtx(t *testing.T) {
	type args struct {
		offset int
		limit  int
	}
	tests := []struct {
		name        string
		args        args
		prepareMock func(t *testing.T) *mock.MockDB
		want        *task.Tasks
		wantErr     bool
	}{
		{
			name: "No error",
			args: args{
				offset: 0,
				limit:  10,
			},
			prepareMock: func(t *testing.T) *mock.MockDB {
				ctrl := gomock.NewController(t)
				mockDB := mock.NewMockDB(ctrl)
				mockDB.EXPECT().SelectCtx(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					0,
					10,
				).Return(nil)
				return mockDB
			},
			wantErr: false,
		},
		{
			name: "Error",
			args: args{
				offset: 0,
				limit:  10,
			},
			prepareMock: func(t *testing.T) *mock.MockDB {
				ctrl := gomock.NewController(t)
				mockDB := mock.NewMockDB(ctrl)
				mockDB.EXPECT().SelectCtx(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					0,
					10,
				).Return(errors.New("Any error"))
				return mockDB
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := tt.prepareMock(t)
			tr := repository.NewTaskRepository(mockDB)
			_, err := tr.ListCtx(context.Background(), tt.args.offset, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("TaskRepository.ListCtx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestTaskRepository_CreateCtx(t *testing.T) {
	type args struct {
		offset int
		limit  int
	}
	tests := []struct {
		name        string
		args        args
		prepareMock func(t *testing.T) *mock.MockDB
		want        *task.Tasks
		wantErr     bool
	}{
		{
			name: "No error",
			args: args{
				offset: 0,
				limit:  10,
			},
			prepareMock: func(t *testing.T) *mock.MockDB {
				ctrl := gomock.NewController(t)
				mockDB := mock.NewMockDB(ctrl)
				mockDB.EXPECT().InsertCtx(
					gomock.Any(),
					gomock.Any(),
				).Return(nil)
				return mockDB
			},
			wantErr: false,
		},
		{
			name: "Error",
			args: args{
				offset: 0,
				limit:  10,
			},
			prepareMock: func(t *testing.T) *mock.MockDB {
				ctrl := gomock.NewController(t)
				mockDB := mock.NewMockDB(ctrl)
				mockDB.EXPECT().InsertCtx(
					gomock.Any(),
					gomock.Any(),
				).Return(errors.New("Any error"))
				return mockDB
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := tt.prepareMock(t)
			tr := repository.NewTaskRepository(mockDB)
			if err := tr.CreateCtx(context.Background(), &task.Task{}); (err != nil) != tt.wantErr {
				t.Errorf("TaskRepository.CreateCtx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTaskRepository_GetCtx(t *testing.T) {
	tests := []struct {
		name        string
		prepareMock func(t *testing.T) *mock.MockDB
		wantErr     bool
	}{
		{
			name: "No error",
			prepareMock: func(t *testing.T) *mock.MockDB {
				ctrl := gomock.NewController(t)
				mockDB := mock.NewMockDB(ctrl)
				mockDB.EXPECT().SelectWithIDCtx(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(nil)
				return mockDB
			},
			wantErr: false,
		},
		{
			name: "Got error",
			prepareMock: func(t *testing.T) *mock.MockDB {
				ctrl := gomock.NewController(t)
				mockDB := mock.NewMockDB(ctrl)
				mockDB.EXPECT().SelectWithIDCtx(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(errors.New("error"))
				return mockDB
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := tt.prepareMock(t)
			tr := repository.NewTaskRepository(mockDB)
			_, err := tr.GetCtx(context.Background(), 1)
			if !tt.wantErr {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
