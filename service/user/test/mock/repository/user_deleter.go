// Code generated by MockGen. DO NOT EDIT.
// Source: ./service/user/internal/repository/user_deleter.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	entity "github.com/indrasaputra/arjuna/service/user/entity"
)

// MockDeleteUserRepository is a mock of DeleteUserRepository interface.
type MockDeleteUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockDeleteUserRepositoryMockRecorder
}

// MockDeleteUserRepositoryMockRecorder is the mock recorder for MockDeleteUserRepository.
type MockDeleteUserRepositoryMockRecorder struct {
	mock *MockDeleteUserRepository
}

// NewMockDeleteUserRepository creates a new mock instance.
func NewMockDeleteUserRepository(ctrl *gomock.Controller) *MockDeleteUserRepository {
	mock := &MockDeleteUserRepository{ctrl: ctrl}
	mock.recorder = &MockDeleteUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDeleteUserRepository) EXPECT() *MockDeleteUserRepositoryMockRecorder {
	return m.recorder
}

// HardDelete mocks base method.
func (m *MockDeleteUserRepository) HardDelete(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HardDelete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// HardDelete indicates an expected call of HardDelete.
func (mr *MockDeleteUserRepositoryMockRecorder) HardDelete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HardDelete", reflect.TypeOf((*MockDeleteUserRepository)(nil).HardDelete), ctx, id)
}

// MockDeleteUserPostgres is a mock of DeleteUserPostgres interface.
type MockDeleteUserPostgres struct {
	ctrl     *gomock.Controller
	recorder *MockDeleteUserPostgresMockRecorder
}

// MockDeleteUserPostgresMockRecorder is the mock recorder for MockDeleteUserPostgres.
type MockDeleteUserPostgresMockRecorder struct {
	mock *MockDeleteUserPostgres
}

// NewMockDeleteUserPostgres creates a new mock instance.
func NewMockDeleteUserPostgres(ctrl *gomock.Controller) *MockDeleteUserPostgres {
	mock := &MockDeleteUserPostgres{ctrl: ctrl}
	mock.recorder = &MockDeleteUserPostgresMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDeleteUserPostgres) EXPECT() *MockDeleteUserPostgresMockRecorder {
	return m.recorder
}

// GetByID mocks base method.
func (m *MockDeleteUserPostgres) GetByID(ctx context.Context, id string) (*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockDeleteUserPostgresMockRecorder) GetByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockDeleteUserPostgres)(nil).GetByID), ctx, id)
}

// HardDelete mocks base method.
func (m *MockDeleteUserPostgres) HardDelete(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HardDelete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// HardDelete indicates an expected call of HardDelete.
func (mr *MockDeleteUserPostgresMockRecorder) HardDelete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HardDelete", reflect.TypeOf((*MockDeleteUserPostgres)(nil).HardDelete), ctx, id)
}
