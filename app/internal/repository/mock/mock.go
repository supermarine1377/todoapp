// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go
//
// Generated by this command:
//
//	mockgen -source=repository.go -destination=./mock/mock.go -package=mock
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockDB is a mock of DB interface.
type MockDB struct {
	ctrl     *gomock.Controller
	recorder *MockDBMockRecorder
	isgomock struct{}
}

// MockDBMockRecorder is the mock recorder for MockDB.
type MockDBMockRecorder struct {
	mock *MockDB
}

// NewMockDB creates a new mock instance.
func NewMockDB(ctrl *gomock.Controller) *MockDB {
	mock := &MockDB{ctrl: ctrl}
	mock.recorder = &MockDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDB) EXPECT() *MockDBMockRecorder {
	return m.recorder
}

// InsertCtx mocks base method.
func (m *MockDB) InsertCtx(ctx context.Context, p any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertCtx", ctx, p)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertCtx indicates an expected call of InsertCtx.
func (mr *MockDBMockRecorder) InsertCtx(ctx, p any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertCtx", reflect.TypeOf((*MockDB)(nil).InsertCtx), ctx, p)
}

// SelectCtx mocks base method.
func (m *MockDB) SelectCtx(ctx context.Context, p any, columns []string, offset, limit int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectCtx", ctx, p, columns, offset, limit)
	ret0, _ := ret[0].(error)
	return ret0
}

// SelectCtx indicates an expected call of SelectCtx.
func (mr *MockDBMockRecorder) SelectCtx(ctx, p, columns, offset, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectCtx", reflect.TypeOf((*MockDB)(nil).SelectCtx), ctx, p, columns, offset, limit)
}

// SelectWithIDCtx mocks base method.
func (m *MockDB) SelectWithIDCtx(ctx context.Context, p any, columns []string, id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectWithIDCtx", ctx, p, columns, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// SelectWithIDCtx indicates an expected call of SelectWithIDCtx.
func (mr *MockDBMockRecorder) SelectWithIDCtx(ctx, p, columns, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectWithIDCtx", reflect.TypeOf((*MockDB)(nil).SelectWithIDCtx), ctx, p, columns, id)
}
