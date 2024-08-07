// Code generated by MockGen. DO NOT EDIT.
// Source: ./service/user/internal/service/user_deleter.go
//
// Generated by this command:
//
//	mockgen -source=./service/user/internal/service/user_deleter.go -destination=./service/user/test/mock//service/user_deleter.go
//

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"

	uow "github.com/indrasaputra/arjuna/pkg/sdk/uow"
	entity "github.com/indrasaputra/arjuna/service/user/entity"
)

// MockDeleteUser is a mock of DeleteUser interface.
type MockDeleteUser struct {
	ctrl     *gomock.Controller
	recorder *MockDeleteUserMockRecorder
}

// MockDeleteUserMockRecorder is the mock recorder for MockDeleteUser.
type MockDeleteUserMockRecorder struct {
	mock *MockDeleteUser
}

// NewMockDeleteUser creates a new mock instance.
func NewMockDeleteUser(ctrl *gomock.Controller) *MockDeleteUser {
	mock := &MockDeleteUser{ctrl: ctrl}
	mock.recorder = &MockDeleteUserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDeleteUser) EXPECT() *MockDeleteUserMockRecorder {
	return m.recorder
}

// HardDelete mocks base method.
func (m *MockDeleteUser) HardDelete(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HardDelete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// HardDelete indicates an expected call of HardDelete.
func (mr *MockDeleteUserMockRecorder) HardDelete(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HardDelete", reflect.TypeOf((*MockDeleteUser)(nil).HardDelete), ctx, id)
}

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

// GetByID mocks base method.
func (m *MockDeleteUserRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockDeleteUserRepositoryMockRecorder) GetByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockDeleteUserRepository)(nil).GetByID), ctx, id)
}

// HardDeleteWithTx mocks base method.
func (m *MockDeleteUserRepository) HardDeleteWithTx(ctx context.Context, tx uow.Tx, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HardDeleteWithTx", ctx, tx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// HardDeleteWithTx indicates an expected call of HardDeleteWithTx.
func (mr *MockDeleteUserRepositoryMockRecorder) HardDeleteWithTx(ctx, tx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HardDeleteWithTx", reflect.TypeOf((*MockDeleteUserRepository)(nil).HardDeleteWithTx), ctx, tx, id)
}
