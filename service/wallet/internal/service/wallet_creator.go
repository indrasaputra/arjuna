package service

import (
	"context"
	"strings"
	"time"

	"github.com/segmentio/ksuid"

	"github.com/indrasaputra/arjuna/service/wallet/entity"
	"github.com/indrasaputra/arjuna/service/wallet/internal/app"
)

// CreateWallet defines interface to create wallet.
type CreateWallet interface {
	// Create creates a new wallet.
	// It needs idempotency key.
	Create(ctx context.Context, wallet *entity.Wallet, key string) error
}

// IdempotencyKeyRepository defines  interface for idempotency check flow and repository.
type IdempotencyKeyRepository interface {
	// Exists check if given key exists in repository.
	Exists(ctx context.Context, key string) (bool, error)
}

// CreateWalletRepository defines the interface to insert wallet to repository.
type CreateWalletRepository interface {
	// Insert inserts a wallet.
	Insert(ctx context.Context, wallet *entity.Wallet) error
}

// WalletCreator is responsible for creating a new wallet.
type WalletCreator struct {
	walletRepo CreateWalletRepository
	keyRepo    IdempotencyKeyRepository
}

// NewWalletCreator creates an instance of WalletCreator.
func NewWalletCreator(t CreateWalletRepository, k IdempotencyKeyRepository) *WalletCreator {
	return &WalletCreator{walletRepo: t, keyRepo: k}
}

// Create creates a new wallet.
// It needs idempotency key.
func (wc *WalletCreator) Create(ctx context.Context, wallet *entity.Wallet, key string) error {
	if err := wc.validateIdempotencyKey(ctx, key); err != nil {
		app.Logger.Errorf(ctx, "[WalletCreator-Create] fail check idempotency key: %s - %v", key, err)
		return err
	}

	sanitizeWallet(wallet)
	if err := validateWallet(wallet); err != nil {
		app.Logger.Errorf(ctx, "[WalletCreator-Create] wallet is invalid: %v", err)
		return err
	}

	if err := setWalletID(ctx, wallet); err != nil {
		app.Logger.Errorf(ctx, "[WalletCreator-Create] fail set wallet id: %v", err)
		return err
	}
	setWalletAuditableProperties(wallet)

	err := wc.walletRepo.Insert(ctx, wallet)
	if err != nil {
		app.Logger.Errorf(ctx, "[WalletCreator-Create] fail save to repository: %v", err)
		return err
	}
	return nil
}

func (wc *WalletCreator) validateIdempotencyKey(ctx context.Context, key string) error {
	res, err := wc.keyRepo.Exists(ctx, key)
	if err != nil {
		return err
	}
	if res {
		return entity.ErrAlreadyExists()
	}
	return nil
}

func sanitizeWallet(wallet *entity.Wallet) {
	if wallet == nil {
		return
	}
	wallet.UserID = strings.TrimSpace(wallet.UserID)
}

func validateWallet(wallet *entity.Wallet) error {
	if wallet == nil {
		return entity.ErrEmptyWallet()
	}
	if wallet.UserID == "" {
		return entity.ErrInvalidUser()
	}
	return nil
}

func setWalletID(ctx context.Context, wallet *entity.Wallet) error {
	id, err := generateUniqueID(ctx)
	if err != nil {
		app.Logger.Errorf(ctx, "[setWalletID] fail generate unique id: %v", err)
		return entity.ErrInternal("fail to create wallet's ID")
	}
	wallet.ID = id
	return nil
}

func generateUniqueID(ctx context.Context) (string, error) {
	id, err := ksuid.NewRandom()
	if err != nil {
		app.Logger.Errorf(ctx, "[generateUniqueID] fail generate ksuid: %v", err)
		return "", entity.ErrInternal("fail to generate unique ID")
	}
	return id.String(), err
}

func setWalletAuditableProperties(wallet *entity.Wallet) {
	wallet.CreatedAt = time.Now().UTC()
	wallet.UpdatedAt = time.Now().UTC()
	wallet.CreatedBy = wallet.UserID
	wallet.UpdatedBy = wallet.UserID
}
