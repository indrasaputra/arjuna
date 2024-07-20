package service

import (
	"context"
	"strings"

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
	// GetUserWalletWithTx gets user's wallet from repository using transaction.
	GetUserWalletWithTx(ctx context.Context, tx uow.Tx, id string, userID string) (*entity.Wallet, error)
	// AddWalletBalanceWithTx adds certain amount (can be negative) to certain wallet using transaction.
	AddWalletBalanceWithTx(ctx context.Context, tx uow.Tx, id string, amount decimal.Decimal) error
}

// WalletTransferer is responsible for transfer balance between wallets.
type WalletTransferer struct {
	walletRepo WalletTransfererRepository
	uow        uow.UnitOfWork
}

// NewWalletTransferer creates an instance of WalletTransferer.
func NewWalletTransferer(w WalletTransfererRepository, u uow.UnitOfWork) *WalletTransferer {
	return &WalletTransferer{walletRepo: w, uow: u}
}

// TransferBalance transfers certain amount of balance from sender to receiver.
// Sender's balance must be sufficient to make a transfer.
func (wt *WalletTransferer) TransferBalance(ctx context.Context, transfer *entity.TransferWallet) error {
	if transfer == nil {
		return entity.ErrInvalidTransfer()
	}
	sanitizeTransferWallet(transfer)
	if err := validateTransferWalletRequest(transfer); err != nil {
		return err
	}
	if err := wt.processTransferBalance(ctx, transfer); err != nil {
		return err
	}
	return nil
}

func (wt *WalletTransferer) processTransferBalance(ctx context.Context, transfer *entity.TransferWallet) error {
	tx, err := wt.uow.Begin(ctx)
	if err != nil {
		app.Logger.Errorf(ctx, "[WalletTransferer-processTransferBalance] fail begin transaction: %v", err)
		return entity.ErrInternal("something went wrong")
	}

	senWallet, recWallet, err := wt.getSenderAndReceiverWallet(ctx, tx, transfer)
	if err != nil {
		_ = wt.uow.Finish(ctx, tx, err)
		return err
	}

	if senWallet == nil || recWallet == nil {
		_ = wt.uow.Finish(ctx, tx, entity.ErrInvalidUser())
		return entity.ErrInvalidUser()
	}
	if senWallet.Balance.LessThan(transfer.Amount) {
		_ = wt.uow.Finish(ctx, tx, entity.ErrInsufficientBalance())
		return entity.ErrInsufficientBalance()
	}

	if err := wt.updateUserBalances(ctx, tx, transfer); err != nil {
		_ = wt.uow.Finish(ctx, tx, err)
		return err
	}
	return wt.uow.Finish(ctx, tx, nil)
}

// deliberately get wallet from the lower id first to avoid deadlock.
func (wt *WalletTransferer) getSenderAndReceiverWallet(ctx context.Context, tx uow.Tx, transfer *entity.TransferWallet) (*entity.Wallet, *entity.Wallet, error) {
	if transfer.SenderWalletID < transfer.ReceiverWalletID {
		senWallet, err := wt.walletRepo.GetUserWalletWithTx(ctx, tx, transfer.SenderWalletID, transfer.SenderID)
		if err != nil {
			app.Logger.Errorf(ctx, "[WalletTransferer-getSenderAndReceiverWallet] sender < receiver; get sender wallet fail: %v", err)
			return nil, nil, err
		}
		recWallet, err := wt.walletRepo.GetUserWalletWithTx(ctx, tx, transfer.ReceiverWalletID, transfer.ReceiverID)
		if err != nil {
			app.Logger.Errorf(ctx, "[WalletTransferer-getSenderAndReceiverWallet] sender < receiver; get receiver wallet fail: %v", err)
			return nil, nil, err
		}
		return senWallet, recWallet, nil
	}

	recWallet, err := wt.walletRepo.GetUserWalletWithTx(ctx, tx, transfer.ReceiverWalletID, transfer.ReceiverID)
	if err != nil {
		app.Logger.Errorf(ctx, "[WalletTransferer-getSenderAndReceiverWallet] sender >= receiver; get receiver wallet fail: %v", err)
		return nil, nil, err
	}
	senWallet, err := wt.walletRepo.GetUserWalletWithTx(ctx, tx, transfer.SenderWalletID, transfer.SenderID)
	if err != nil {
		app.Logger.Errorf(ctx, "[WalletTransferer-getSenderAndReceiverWallet] sender >= receiver; get sender wallet fail: %v", err)
		return nil, nil, err
	}
	return senWallet, recWallet, nil
}

func (wt *WalletTransferer) updateUserBalances(ctx context.Context, tx uow.Tx, transfer *entity.TransferWallet) error {
	if err := wt.walletRepo.AddWalletBalanceWithTx(ctx, tx, transfer.SenderWalletID, transfer.Amount.Neg()); err != nil {
		app.Logger.Errorf(ctx, "[WalletTransferer-updateUserBalances] subtract sender balance fail: %v", err)
		return err
	}
	if err := wt.walletRepo.AddWalletBalanceWithTx(ctx, tx, transfer.ReceiverWalletID, transfer.Amount); err != nil {
		app.Logger.Errorf(ctx, "[WalletTransferer-updateUserBalances] add receiver balance fail: %v", err)
		return err
	}
	return nil
}

func sanitizeTransferWallet(transfer *entity.TransferWallet) {
	transfer.SenderID = strings.TrimSpace(transfer.SenderID)
	transfer.SenderWalletID = strings.TrimSpace(transfer.SenderWalletID)
	transfer.ReceiverID = strings.TrimSpace(transfer.ReceiverID)
	transfer.ReceiverWalletID = strings.TrimSpace(transfer.ReceiverWalletID)
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
