// Code generated by MockGen. DO NOT EDIT.
// Source: ./service/user/internal/service/user_registrar.go
//
// Generated by this command:
//
//	mockgen -source=./service/user/internal/service/user_registrar.go -destination=./service/user/test/mock//service/user_registrar.go
//

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	uow "github.com/indrasaputra/arjuna/pkg/sdk/uow"
	entity "github.com/indrasaputra/arjuna/service/user/entity"
	gomock "go.uber.org/mock/gomock"
)

// MockRegisterUser is a mock of RegisterUser interface.
type MockRegisterUser struct {
	ctrl     *gomock.Controller
	recorder *MockRegisterUserMockRecorder
}

// MockRegisterUserMockRecorder is the mock recorder for MockRegisterUser.
type MockRegisterUserMockRecorder struct {
	mock *MockRegisterUser
}

// NewMockRegisterUser creates a new mock instance.
func NewMockRegisterUser(ctrl *gomock.Controller) *MockRegisterUser {
	mock := &MockRegisterUser{ctrl: ctrl}
	mock.recorder = &MockRegisterUserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRegisterUser) EXPECT() *MockRegisterUserMockRecorder {
	return m.recorder
}

// Register mocks base method.
func (m *MockRegisterUser) Register(ctx context.Context, user *entity.User, key string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", ctx, user, key)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Register indicates an expected call of Register.
func (mr *MockRegisterUserMockRecorder) Register(ctx, user, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockRegisterUser)(nil).Register), ctx, user, key)
}

// MockRegisterUserRepository is a mock of RegisterUserRepository interface.
type MockRegisterUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRegisterUserRepositoryMockRecorder
}

// MockRegisterUserRepositoryMockRecorder is the mock recorder for MockRegisterUserRepository.
type MockRegisterUserRepositoryMockRecorder struct {
	mock *MockRegisterUserRepository
}

// NewMockRegisterUserRepository creates a new mock instance.
func NewMockRegisterUserRepository(ctrl *gomock.Controller) *MockRegisterUserRepository {
	mock := &MockRegisterUserRepository{ctrl: ctrl}
	mock.recorder = &MockRegisterUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRegisterUserRepository) EXPECT() *MockRegisterUserRepositoryMockRecorder {
	return m.recorder
}

// InsertWithTx mocks base method.
func (m *MockRegisterUserRepository) InsertWithTx(ctx context.Context, tx uow.Tx, user *entity.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertWithTx", ctx, tx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertWithTx indicates an expected call of InsertWithTx.
func (mr *MockRegisterUserRepositoryMockRecorder) InsertWithTx(ctx, tx, user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertWithTx", reflect.TypeOf((*MockRegisterUserRepository)(nil).InsertWithTx), ctx, tx, user)
}

// MockRegisterUserOutboxRepository is a mock of RegisterUserOutboxRepository interface.
type MockRegisterUserOutboxRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRegisterUserOutboxRepositoryMockRecorder
}

// MockRegisterUserOutboxRepositoryMockRecorder is the mock recorder for MockRegisterUserOutboxRepository.
type MockRegisterUserOutboxRepositoryMockRecorder struct {
	mock *MockRegisterUserOutboxRepository
}

// NewMockRegisterUserOutboxRepository creates a new mock instance.
func NewMockRegisterUserOutboxRepository(ctrl *gomock.Controller) *MockRegisterUserOutboxRepository {
	mock := &MockRegisterUserOutboxRepository{ctrl: ctrl}
	mock.recorder = &MockRegisterUserOutboxRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRegisterUserOutboxRepository) EXPECT() *MockRegisterUserOutboxRepositoryMockRecorder {
	return m.recorder
}

// InsertWithTx mocks base method.
func (m *MockRegisterUserOutboxRepository) InsertWithTx(ctx context.Context, tx uow.Tx, payload *entity.UserOutbox) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertWithTx", ctx, tx, payload)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertWithTx indicates an expected call of InsertWithTx.
func (mr *MockRegisterUserOutboxRepositoryMockRecorder) InsertWithTx(ctx, tx, payload any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertWithTx", reflect.TypeOf((*MockRegisterUserOutboxRepository)(nil).InsertWithTx), ctx, tx, payload)
}

// MockIdempotencyKeyRepository is a mock of IdempotencyKeyRepository interface.
type MockIdempotencyKeyRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIdempotencyKeyRepositoryMockRecorder
}

// MockIdempotencyKeyRepositoryMockRecorder is the mock recorder for MockIdempotencyKeyRepository.
type MockIdempotencyKeyRepositoryMockRecorder struct {
	mock *MockIdempotencyKeyRepository
}

// NewMockIdempotencyKeyRepository creates a new mock instance.
func NewMockIdempotencyKeyRepository(ctrl *gomock.Controller) *MockIdempotencyKeyRepository {
	mock := &MockIdempotencyKeyRepository{ctrl: ctrl}
	mock.recorder = &MockIdempotencyKeyRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIdempotencyKeyRepository) EXPECT() *MockIdempotencyKeyRepositoryMockRecorder {
	return m.recorder
}

// Exists mocks base method.
func (m *MockIdempotencyKeyRepository) Exists(ctx context.Context, key string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exists", ctx, key)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exists indicates an expected call of Exists.
func (mr *MockIdempotencyKeyRepositoryMockRecorder) Exists(ctx, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exists", reflect.TypeOf((*MockIdempotencyKeyRepository)(nil).Exists), ctx, key)
}
