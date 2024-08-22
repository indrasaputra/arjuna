// Code generated by MockGen. DO NOT EDIT.
// Source: ./service/transaction/internal/service/transaction_creator.go
//
// Generated by this command:
//
//	mockgen -source=./service/transaction/internal/service/transaction_creator.go -destination=./service/transaction/test/mock//service/transaction_creator.go
//

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"

	entity "github.com/indrasaputra/arjuna/service/transaction/entity"
)

// MockCreateTransaction is a mock of CreateTransaction interface.
type MockCreateTransaction struct {
	ctrl     *gomock.Controller
	recorder *MockCreateTransactionMockRecorder
}

// MockCreateTransactionMockRecorder is the mock recorder for MockCreateTransaction.
type MockCreateTransactionMockRecorder struct {
	mock *MockCreateTransaction
}

// NewMockCreateTransaction creates a new mock instance.
func NewMockCreateTransaction(ctrl *gomock.Controller) *MockCreateTransaction {
	mock := &MockCreateTransaction{ctrl: ctrl}
	mock.recorder = &MockCreateTransactionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCreateTransaction) EXPECT() *MockCreateTransactionMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockCreateTransaction) Create(ctx context.Context, transaction *entity.Transaction, key string) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, transaction, key)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockCreateTransactionMockRecorder) Create(ctx, transaction, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockCreateTransaction)(nil).Create), ctx, transaction, key)
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

// MockCreateTransactionRepository is a mock of CreateTransactionRepository interface.
type MockCreateTransactionRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCreateTransactionRepositoryMockRecorder
}

// MockCreateTransactionRepositoryMockRecorder is the mock recorder for MockCreateTransactionRepository.
type MockCreateTransactionRepositoryMockRecorder struct {
	mock *MockCreateTransactionRepository
}

// NewMockCreateTransactionRepository creates a new mock instance.
func NewMockCreateTransactionRepository(ctrl *gomock.Controller) *MockCreateTransactionRepository {
	mock := &MockCreateTransactionRepository{ctrl: ctrl}
	mock.recorder = &MockCreateTransactionRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCreateTransactionRepository) EXPECT() *MockCreateTransactionRepositoryMockRecorder {
	return m.recorder
}

// Insert mocks base method.
func (m *MockCreateTransactionRepository) Insert(ctx context.Context, transaction *entity.Transaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, transaction)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert.
func (mr *MockCreateTransactionRepositoryMockRecorder) Insert(ctx, transaction any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockCreateTransactionRepository)(nil).Insert), ctx, transaction)
}
