package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/indrasaputra/arjuna/service/wallet/entity"
)

// TopupWallet defines interface to topup wallet.
type TopupWallet interface {
	// Topup topups a wallet's balance.
	// It needs idempotency key.
	Topup(ctx context.Context, topup *entity.TopupWallet) (*entity.Wallet, error)
}

// TopupWalletRepository defines the interface to update wallet in repository.
type TopupWalletRepository interface {
	// AddWalletBalance adds certain amount (can be negative) to certain wallet.
	AddWalletBalance(ctx context.Context, id uuid.UUID, amount decimal.Decimal) (*entity.Wallet, error)
}

// WalletTopup is responsible for topup a new wallet.
type WalletTopup struct {
	walletRepo TopupWalletRepository
}

// NewWalletTopup topups an instance of WalletTopup.
func NewWalletTopup(t TopupWalletRepository) *WalletTopup {
	return &WalletTopup{walletRepo: t}
}

// Topup topups wallet's balance.
// It needs idempotency key.
func (wt *WalletTopup) Topup(ctx context.Context, topup *entity.TopupWallet) (*entity.Wallet, error) {
	if topup == nil {
		return nil, entity.ErrEmptyWallet()
	}

	if err := validateTopupWallet(topup); err != nil {
		slog.ErrorContext(ctx, "[WalletTopup-Topup] wallet is invalid", "error", err)
		return nil, err
	}

	wallet, err := wt.walletRepo.AddWalletBalance(ctx, topup.WalletID, topup.Amount)
	if err != nil {
		slog.ErrorContext(ctx, "[WalletTopup-Topup] fail update wallet balance", "error", err)
		return nil, err
	}
	return wallet, nil
}

func validateTopupWallet(topup *entity.TopupWallet) error {
	if topup.WalletID == uuid.Nil {
		return entity.ErrEmptyWallet()
	}
	if topup.UserID == uuid.Nil {
		return entity.ErrInvalidUser()
	}
	if decimal.Zero.Equal(topup.Amount) {
		return entity.ErrInvalidAmount()
	}
	return nil
}
