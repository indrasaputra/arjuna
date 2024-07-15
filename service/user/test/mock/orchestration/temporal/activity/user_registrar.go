// Code generated by MockGen. DO NOT EDIT.
// Source: ./service/user/internal/orchestration/temporal/activity/user_registrar.go
//
// Generated by this command:
//
//	mockgen -source=./service/user/internal/orchestration/temporal/activity/user_registrar.go -destination=./service/user/test/mock//orchestration/temporal/activity/user_registrar.go
//

// Package mock_activity is a generated GoMock package.
package mock_activity

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"

	entity "github.com/indrasaputra/arjuna/service/user/entity"
)

// MockRegisterUserConnection is a mock of RegisterUserConnection interface.
type MockRegisterUserConnection struct {
	ctrl     *gomock.Controller
	recorder *MockRegisterUserConnectionMockRecorder
}

// MockRegisterUserConnectionMockRecorder is the mock recorder for MockRegisterUserConnection.
type MockRegisterUserConnectionMockRecorder struct {
	mock *MockRegisterUserConnection
}

// NewMockRegisterUserConnection creates a new mock instance.
func NewMockRegisterUserConnection(ctrl *gomock.Controller) *MockRegisterUserConnection {
	mock := &MockRegisterUserConnection{ctrl: ctrl}
	mock.recorder = &MockRegisterUserConnectionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRegisterUserConnection) EXPECT() *MockRegisterUserConnectionMockRecorder {
	return m.recorder
}

// CreateAccount mocks base method.
func (m *MockRegisterUserConnection) CreateAccount(ctx context.Context, user *entity.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccount", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateAccount indicates an expected call of CreateAccount.
func (mr *MockRegisterUserConnectionMockRecorder) CreateAccount(ctx, user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccount", reflect.TypeOf((*MockRegisterUserConnection)(nil).CreateAccount), ctx, user)
}

// MockRegisterUserDatabase is a mock of RegisterUserDatabase interface.
type MockRegisterUserDatabase struct {
	ctrl     *gomock.Controller
	recorder *MockRegisterUserDatabaseMockRecorder
}

// MockRegisterUserDatabaseMockRecorder is the mock recorder for MockRegisterUserDatabase.
type MockRegisterUserDatabaseMockRecorder struct {
	mock *MockRegisterUserDatabase
}

// NewMockRegisterUserDatabase creates a new mock instance.
func NewMockRegisterUserDatabase(ctrl *gomock.Controller) *MockRegisterUserDatabase {
	mock := &MockRegisterUserDatabase{ctrl: ctrl}
	mock.recorder = &MockRegisterUserDatabaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRegisterUserDatabase) EXPECT() *MockRegisterUserDatabaseMockRecorder {
	return m.recorder
}

// HardDelete mocks base method.
func (m *MockRegisterUserDatabase) HardDelete(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HardDelete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// HardDelete indicates an expected call of HardDelete.
func (mr *MockRegisterUserDatabaseMockRecorder) HardDelete(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HardDelete", reflect.TypeOf((*MockRegisterUserDatabase)(nil).HardDelete), ctx, id)
}
