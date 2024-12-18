// Code generated by MockGen. DO NOT EDIT.
// Source: task.go
//
// Generated by this command:
//
//	mockgen -source=task.go -destination=./mock/mock.go -package=mock
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	task "github.com/supermarine1377/todoapp/app/internal/model/entity/task"
	gomock "go.uber.org/mock/gomock"
)

// MockTaskRepository is a mock of TaskRepository interface.
type MockTaskRepository struct {
	ctrl     *gomock.Controller
	recorder *MockTaskRepositoryMockRecorder
	isgomock struct{}
}

// MockTaskRepositoryMockRecorder is the mock recorder for MockTaskRepository.
type MockTaskRepositoryMockRecorder struct {
	mock *MockTaskRepository
}

// NewMockTaskRepository creates a new mock instance.
func NewMockTaskRepository(ctrl *gomock.Controller) *MockTaskRepository {
	mock := &MockTaskRepository{ctrl: ctrl}
	mock.recorder = &MockTaskRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTaskRepository) EXPECT() *MockTaskRepositoryMockRecorder {
	return m.recorder
}

// CreateCtx mocks base method.
func (m *MockTaskRepository) CreateCtx(ctx context.Context, task *task.Task) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCtx", ctx, task)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateCtx indicates an expected call of CreateCtx.
func (mr *MockTaskRepositoryMockRecorder) CreateCtx(ctx, task any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCtx", reflect.TypeOf((*MockTaskRepository)(nil).CreateCtx), ctx, task)
}

// GetCtx mocks base method.
func (m *MockTaskRepository) GetCtx(ctx context.Context, id int) (*task.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCtx", ctx, id)
	ret0, _ := ret[0].(*task.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCtx indicates an expected call of GetCtx.
func (mr *MockTaskRepositoryMockRecorder) GetCtx(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCtx", reflect.TypeOf((*MockTaskRepository)(nil).GetCtx), ctx, id)
}

// ListCtx mocks base method.
func (m *MockTaskRepository) ListCtx(ctx context.Context, offset, limit int) (*task.Tasks, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListCtx", ctx, offset, limit)
	ret0, _ := ret[0].(*task.Tasks)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListCtx indicates an expected call of ListCtx.
func (mr *MockTaskRepositoryMockRecorder) ListCtx(ctx, offset, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListCtx", reflect.TypeOf((*MockTaskRepository)(nil).ListCtx), ctx, offset, limit)
}
