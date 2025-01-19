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
	Topup(ctx context.Context, topup *entity.TopupWallet) error
}

// TopupWalletRepository defines the interface to update wallet in repository.
type TopupWalletRepository interface {
	// AddWalletBalance adds certain amount (can be negative) to certain wallet.
	AddWalletBalance(ctx context.Context, id uuid.UUID, amount decimal.Decimal) error
}

// IdempotencyKeyRepository defines  interface for idempotency check flow and repository.
type IdempotencyKeyRepository interface {
	// Exists check if given key exists in repository.
	Exists(ctx context.Context, key string) (bool, error)
}

// WalletTopup is responsible for topup a new wallet.
type WalletTopup struct {
	walletRepo TopupWalletRepository
	keyRepo    IdempotencyKeyRepository
}

// NewWalletTopup topups an instance of WalletTopup.
func NewWalletTopup(t TopupWalletRepository, k IdempotencyKeyRepository) *WalletTopup {
	return &WalletTopup{walletRepo: t, keyRepo: k}
}

// Topup topups wallet's balance.
// It needs idempotency key.
func (wt *WalletTopup) Topup(ctx context.Context, topup *entity.TopupWallet) error {
	if topup == nil {
		return entity.ErrEmptyWallet()
	}

	if err := wt.validateIdempotencyKey(ctx, topup.IdempotencyKey); err != nil {
		slog.ErrorContext(ctx, "[WalletTopup-Topup] fail check idempotency key", "idempotency_key", topup.IdempotencyKey, "error", err)
		return err
	}

	if err := validateTopupWallet(topup); err != nil {
		slog.ErrorContext(ctx, "[WalletTopup-Topup] wallet is invalid", "error", err)
		return err
	}

	err := wt.walletRepo.AddWalletBalance(ctx, topup.WalletID, topup.Amount)
	if err != nil {
		slog.ErrorContext(ctx, "[WalletTopup-Topup] fail update wallet balance", "error", err)
		return err
	}
	return nil
}

func (wt *WalletTopup) validateIdempotencyKey(ctx context.Context, key string) error {
	res, err := wt.keyRepo.Exists(ctx, key)
	if err != nil {
		return err
	}
	if res {
		return entity.ErrAlreadyExists()
	}
	return nil
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
