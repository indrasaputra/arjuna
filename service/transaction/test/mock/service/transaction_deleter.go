// Code generated by MockGen. DO NOT EDIT.
// Source: ./service/transaction/internal/service/transaction_deleter.go
//
// Generated by this command:
//
//	mockgen -source=./service/transaction/internal/service/transaction_deleter.go -destination=./service/transaction/test/mock//service/transaction_deleter.go
//

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockDeleteTransaction is a mock of DeleteTransaction interface.
type MockDeleteTransaction struct {
	ctrl     *gomock.Controller
	recorder *MockDeleteTransactionMockRecorder
}

// MockDeleteTransactionMockRecorder is the mock recorder for MockDeleteTransaction.
type MockDeleteTransactionMockRecorder struct {
	mock *MockDeleteTransaction
}

// NewMockDeleteTransaction creates a new mock instance.
func NewMockDeleteTransaction(ctrl *gomock.Controller) *MockDeleteTransaction {
	mock := &MockDeleteTransaction{ctrl: ctrl}
	mock.recorder = &MockDeleteTransactionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDeleteTransaction) EXPECT() *MockDeleteTransactionMockRecorder {
	return m.recorder
}

// DeleteAllTransactions mocks base method.
func (m *MockDeleteTransaction) DeleteAllTransactions(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAllTransactions", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAllTransactions indicates an expected call of DeleteAllTransactions.
func (mr *MockDeleteTransactionMockRecorder) DeleteAllTransactions(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAllTransactions", reflect.TypeOf((*MockDeleteTransaction)(nil).DeleteAllTransactions), ctx)
}

// MockDeleteTransactionRepository is a mock of DeleteTransactionRepository interface.
type MockDeleteTransactionRepository struct {
	ctrl     *gomock.Controller
	recorder *MockDeleteTransactionRepositoryMockRecorder
}

// MockDeleteTransactionRepositoryMockRecorder is the mock recorder for MockDeleteTransactionRepository.
type MockDeleteTransactionRepositoryMockRecorder struct {
	mock *MockDeleteTransactionRepository
}

// NewMockDeleteTransactionRepository creates a new mock instance.
func NewMockDeleteTransactionRepository(ctrl *gomock.Controller) *MockDeleteTransactionRepository {
	mock := &MockDeleteTransactionRepository{ctrl: ctrl}
	mock.recorder = &MockDeleteTransactionRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDeleteTransactionRepository) EXPECT() *MockDeleteTransactionRepositoryMockRecorder {
	return m.recorder
}

// DeleteAll mocks base method.
func (m *MockDeleteTransactionRepository) DeleteAll(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAll", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAll indicates an expected call of DeleteAll.
func (mr *MockDeleteTransactionRepositoryMockRecorder) DeleteAll(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAll", reflect.TypeOf((*MockDeleteTransactionRepository)(nil).DeleteAll), ctx)
}
