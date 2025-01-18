// Code generated by MockGen. DO NOT EDIT.
// Source: ./service/wallet/internal/service/wallet_transferer.go
//
// Generated by this command:
//
//	mockgen -source=./service/wallet/internal/service/wallet_transferer.go -destination=./service/wallet/test/mock//service/wallet_transferer.go
//

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	uuid "github.com/google/uuid"
	decimal "github.com/shopspring/decimal"
	gomock "go.uber.org/mock/gomock"

	entity "github.com/indrasaputra/arjuna/service/wallet/entity"
)

// MockTransferWallet is a mock of TransferWallet interface.
type MockTransferWallet struct {
	ctrl     *gomock.Controller
	recorder *MockTransferWalletMockRecorder
}

// MockTransferWalletMockRecorder is the mock recorder for MockTransferWallet.
type MockTransferWalletMockRecorder struct {
	mock *MockTransferWallet
}

// NewMockTransferWallet creates a new mock instance.
func NewMockTransferWallet(ctrl *gomock.Controller) *MockTransferWallet {
	mock := &MockTransferWallet{ctrl: ctrl}
	mock.recorder = &MockTransferWalletMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransferWallet) EXPECT() *MockTransferWalletMockRecorder {
	return m.recorder
}

// TransferBalance mocks base method.
func (m *MockTransferWallet) TransferBalance(ctx context.Context, transfer *entity.TransferWallet) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TransferBalance", ctx, transfer)
	ret0, _ := ret[0].(error)
	return ret0
}

// TransferBalance indicates an expected call of TransferBalance.
func (mr *MockTransferWalletMockRecorder) TransferBalance(ctx, transfer any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TransferBalance", reflect.TypeOf((*MockTransferWallet)(nil).TransferBalance), ctx, transfer)
}

// MockWalletTransfererRepository is a mock of WalletTransfererRepository interface.
type MockWalletTransfererRepository struct {
	ctrl     *gomock.Controller
	recorder *MockWalletTransfererRepositoryMockRecorder
}

// MockWalletTransfererRepositoryMockRecorder is the mock recorder for MockWalletTransfererRepository.
type MockWalletTransfererRepositoryMockRecorder struct {
	mock *MockWalletTransfererRepository
}

// NewMockWalletTransfererRepository creates a new mock instance.
func NewMockWalletTransfererRepository(ctrl *gomock.Controller) *MockWalletTransfererRepository {
	mock := &MockWalletTransfererRepository{ctrl: ctrl}
	mock.recorder = &MockWalletTransfererRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWalletTransfererRepository) EXPECT() *MockWalletTransfererRepositoryMockRecorder {
	return m.recorder
}

// AddWalletBalance mocks base method.
func (m *MockWalletTransfererRepository) AddWalletBalance(ctx context.Context, id uuid.UUID, amount decimal.Decimal) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddWalletBalance", ctx, id, amount)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddWalletBalance indicates an expected call of AddWalletBalance.
func (mr *MockWalletTransfererRepositoryMockRecorder) AddWalletBalance(ctx, id, amount any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddWalletBalance", reflect.TypeOf((*MockWalletTransfererRepository)(nil).AddWalletBalance), ctx, id, amount)
}

// GetUserWalletForUpdate mocks base method.
func (m *MockWalletTransfererRepository) GetUserWalletForUpdate(ctx context.Context, id, userID uuid.UUID) (*entity.Wallet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserWalletForUpdate", ctx, id, userID)
	ret0, _ := ret[0].(*entity.Wallet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserWalletForUpdate indicates an expected call of GetUserWalletForUpdate.
func (mr *MockWalletTransfererRepositoryMockRecorder) GetUserWalletForUpdate(ctx, id, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserWalletForUpdate", reflect.TypeOf((*MockWalletTransfererRepository)(nil).GetUserWalletForUpdate), ctx, id, userID)
}
