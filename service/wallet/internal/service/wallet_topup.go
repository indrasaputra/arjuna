package service

import (
	"context"
	"strings"

	"github.com/shopspring/decimal"

	"github.com/indrasaputra/arjuna/service/wallet/entity"
	"github.com/indrasaputra/arjuna/service/wallet/internal/app"
)

// TopupWallet defines interface to topup wallet.
type TopupWallet interface {
	// Topup topups a wallet's balance.
	// It needs idempotency key.
	Topup(ctx context.Context, topup *entity.TopupWallet) error
}

// TopupWalletRepository defines the interface to update wallet in repository.
type TopupWalletRepository interface {
	// AddWalletBalance updates the balance.
	AddWalletBalance(ctx context.Context, topup *entity.TopupWallet) error
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
		app.Logger.Errorf(ctx, "[UserRegistrar-Register] fail check idempotency key: %s - %v", topup.IdempotencyKey, err)
		return err
	}

	sanitizeTopupWallet(topup)
	if err := validateTopupWallet(topup); err != nil {
		app.Logger.Errorf(ctx, "[WalletTopup-Topup] wallet is invalid: %v", err)
		return err
	}

	err := wt.walletRepo.AddWalletBalance(ctx, topup)
	if err != nil {
		app.Logger.Errorf(ctx, "[WalletTopup-Topup] fail update to repository: %v", err)
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

func sanitizeTopupWallet(topup *entity.TopupWallet) {
	topup.WalletID = strings.TrimSpace(topup.WalletID)
	topup.UserID = strings.TrimSpace(topup.UserID)
}

func validateTopupWallet(topup *entity.TopupWallet) error {
	if topup == nil {
		return entity.ErrEmptyWallet()
	}
	if topup.WalletID == "" {
		return entity.ErrEmptyWallet()
	}
	if topup.UserID == "" {
		return entity.ErrInvalidUser()
	}
	if decimal.Zero.Equal(topup.Amount) {
		return entity.ErrInvalidAmount()
	}
	return nil
}
