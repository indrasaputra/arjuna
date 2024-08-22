package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/indrasaputra/arjuna/service/wallet/entity"
	"github.com/indrasaputra/arjuna/service/wallet/internal/app"
)

// CreateWallet defines interface to create wallet.
type CreateWallet interface {
	// Create creates a new wallet.
	Create(ctx context.Context, wallet *entity.Wallet) error
}

// CreateWalletRepository defines the interface to insert wallet to repository.
type CreateWalletRepository interface {
	// Insert inserts a wallet.
	Insert(ctx context.Context, wallet *entity.Wallet) error
}

// WalletCreator is responsible for creating a new wallet.
type WalletCreator struct {
	walletRepo CreateWalletRepository
}

// NewWalletCreator creates an instance of WalletCreator.
func NewWalletCreator(t CreateWalletRepository) *WalletCreator {
	return &WalletCreator{walletRepo: t}
}

// Create creates a new wallet.
// It needs idempotency key.
func (wc *WalletCreator) Create(ctx context.Context, wallet *entity.Wallet) error {
	if err := validateWallet(wallet); err != nil {
		app.Logger.Errorf(ctx, "[WalletCreator-Create] wallet is invalid: %v", err)
		return err
	}

	setWalletID(wallet)
	setWalletAuditableProperties(wallet)

	err := wc.walletRepo.Insert(ctx, wallet)
	if err != nil {
		app.Logger.Errorf(ctx, "[WalletCreator-Create] fail save to repository: %v", err)
		return err
	}
	return nil
}

func validateWallet(wallet *entity.Wallet) error {
	if wallet == nil {
		return entity.ErrEmptyWallet()
	}
	if wallet.UserID == uuid.Nil {
		return entity.ErrInvalidUser()
	}
	return nil
}

func setWalletID(wallet *entity.Wallet) {
	wallet.ID = generateUniqueID()
}

func generateUniqueID() uuid.UUID {
	return uuid.Must(uuid.NewV7())
}

func setWalletAuditableProperties(wallet *entity.Wallet) {
	wallet.CreatedAt = time.Now().UTC()
	wallet.UpdatedAt = time.Now().UTC()
	wallet.CreatedBy = wallet.UserID.String()
	wallet.UpdatedBy = wallet.UserID.String()
}
