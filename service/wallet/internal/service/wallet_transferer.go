package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/indrasaputra/arjuna/pkg/sdk/uow"
	"github.com/indrasaputra/arjuna/service/wallet/entity"
	"github.com/indrasaputra/arjuna/service/wallet/internal/app"
)

// TransferWallet defines interface to transfer wallet.
type TransferWallet interface {
	// TransferBalance transfers a wallet's balance.
	TransferBalance(ctx context.Context, transfer *entity.TransferWallet) error
}

// WalletTransfererRepository defines the interface to get wallet in repository.
type WalletTransfererRepository interface {
	// GetUserWalletForUpdate gets user's wallet from repository for update.
	GetUserWalletForUpdate(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entity.Wallet, error)
	// AddWalletBalance adds certain amount (can be negative) to certain wallet.
	AddWalletBalance(ctx context.Context, id uuid.UUID, amount decimal.Decimal) error
}

// WalletTransferer is responsible for transfer balance between wallets.
type WalletTransferer struct {
	walletRepo WalletTransfererRepository
	txManager  uow.TxManager
}

// NewWalletTransferer creates an instance of WalletTransferer.
func NewWalletTransferer(w WalletTransfererRepository, m uow.TxManager) *WalletTransferer {
	return &WalletTransferer{walletRepo: w, txManager: m}
}

// TransferBalance transfers certain amount of balance from sender to receiver.
// Sender's balance must be sufficient to make a transfer.
func (wt *WalletTransferer) TransferBalance(ctx context.Context, transfer *entity.TransferWallet) error {
	if transfer == nil {
		return entity.ErrInvalidTransfer()
	}
	if err := validateTransferWalletRequest(transfer); err != nil {
		return err
	}
	if err := wt.processTransferBalance(ctx, transfer); err != nil {
		return err
	}
	return nil
}

func (wt *WalletTransferer) processTransferBalance(ctx context.Context, transfer *entity.TransferWallet) error {
	err := wt.txManager.Do(ctx, func(ctx context.Context) error {
		senWallet, recWallet, err := wt.getSenderAndReceiverWallet(ctx, transfer)
		if err != nil {
			return err
		}

		if senWallet == nil || recWallet == nil {
			return entity.ErrInvalidUser()
		}
		if senWallet.Balance.LessThan(transfer.Amount) {
			return entity.ErrInsufficientBalance()
		}

		if err := wt.updateUserBalances(ctx, transfer); err != nil {
			return err
		}
		return nil
	})
	return err
}

// Prevent deadlocks by consistently ordering wallet lock acquisition:
// We always lock wallets in order of ascending wallet ID, regardless of whether a wallet
// is the sender or receiver. This prevents deadlock scenarios where:
// 1. Transaction A transfers from wallet 2 -> wallet 1
// 2. Transaction B transfers from wallet 1 -> wallet 2
// Without ordering, A could lock wallet 2 while B locks wallet 1, leading to deadlock.
// By always locking the lower ID first, both transactions will attempt to lock wallet 1
// before wallet 2, creating a consistent ordering that prevents circular wait conditions.
func (wt *WalletTransferer) getSenderAndReceiverWallet(ctx context.Context, transfer *entity.TransferWallet) (*entity.Wallet, *entity.Wallet, error) {
	if transfer.SenderWalletID.String() < transfer.ReceiverWalletID.String() {
		senWallet, err := wt.walletRepo.GetUserWalletForUpdate(ctx, transfer.SenderWalletID, transfer.SenderID)
		if err != nil {
			app.Logger.Errorf(ctx, "[WalletTransferer-getSenderAndReceiverWallet] sender < receiver; get sender wallet fail: %v", err)
			return nil, nil, err
		}
		recWallet, err := wt.walletRepo.GetUserWalletForUpdate(ctx, transfer.ReceiverWalletID, transfer.ReceiverID)
		if err != nil {
			app.Logger.Errorf(ctx, "[WalletTransferer-getSenderAndReceiverWallet] sender < receiver; get receiver wallet fail: %v", err)
			return nil, nil, err
		}
		return senWallet, recWallet, nil
	}

	recWallet, err := wt.walletRepo.GetUserWalletForUpdate(ctx, transfer.ReceiverWalletID, transfer.ReceiverID)
	if err != nil {
		app.Logger.Errorf(ctx, "[WalletTransferer-getSenderAndReceiverWallet] sender >= receiver; get receiver wallet fail: %v", err)
		return nil, nil, err
	}
	senWallet, err := wt.walletRepo.GetUserWalletForUpdate(ctx, transfer.SenderWalletID, transfer.SenderID)
	if err != nil {
		app.Logger.Errorf(ctx, "[WalletTransferer-getSenderAndReceiverWallet] sender >= receiver; get sender wallet fail: %v", err)
		return nil, nil, err
	}
	return senWallet, recWallet, nil
}

func (wt *WalletTransferer) updateUserBalances(ctx context.Context, transfer *entity.TransferWallet) error {
	if err := wt.walletRepo.AddWalletBalance(ctx, transfer.SenderWalletID, transfer.Amount.Neg()); err != nil {
		app.Logger.Errorf(ctx, "[WalletTransferer-updateUserBalances] subtract sender balance fail: %v", err)
		return err
	}
	if err := wt.walletRepo.AddWalletBalance(ctx, transfer.ReceiverWalletID, transfer.Amount); err != nil {
		app.Logger.Errorf(ctx, "[WalletTransferer-updateUserBalances] add receiver balance fail: %v", err)
		return err
	}
	return nil
}

func validateTransferWalletRequest(transfer *entity.TransferWallet) error {
	if transfer.SenderID == transfer.ReceiverID {
		return entity.ErrSameAccount()
	}
	if transfer.Amount.Equal(decimal.Zero) {
		return entity.ErrInvalidAmount()
	}
	return nil
}
