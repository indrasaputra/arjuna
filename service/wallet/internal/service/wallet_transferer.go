package service

import (
	"context"
	"strings"

	"github.com/shopspring/decimal"

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
	// GetUserWallet gets user's wallet from repository.
	GetUserWallet(ctx context.Context, id string, userID string) (*entity.Wallet, error)
	// AddWalletBalance adds certain amount (can be negative) to certain wallet.
	AddWalletBalance(ctx context.Context, id string, amount decimal.Decimal) error
}

// WalletTransferer is responsible for transfer balance between wallets.
type WalletTransferer struct {
	walletRepo WalletTransfererRepository
}

// NewWalletTransferer creates an instance of WalletTransferer.
func NewWalletTransferer(w WalletTransfererRepository) *WalletTransferer {
	return &WalletTransferer{walletRepo: w}
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
}

// deliberately get wallet from the lower id first to avoid deadlock.
func (wt *WalletTransferer) getSenderAndReceiverWallet(ctx context.Context, transfer *entity.TransferWallet) (*entity.Wallet, *entity.Wallet, error) {
	if transfer.SenderWalletID < transfer.ReceiverWalletID {
		senWallet, err := wt.walletRepo.GetUserWallet(ctx, transfer.SenderWalletID, transfer.SenderID)
		if err != nil {
			app.Logger.Errorf(ctx, "[WalletTransferer-getSenderAndReceiverWallet] sender < receiver; get sender wallet fail: %v", err)
			return nil, nil, err
		}
		recWallet, err := wt.walletRepo.GetUserWallet(ctx, transfer.ReceiverWalletID, transfer.ReceiverID)
		if err != nil {
			app.Logger.Errorf(ctx, "[WalletTransferer-getSenderAndReceiverWallet] sender < receiver; get receiver wallet fail: %v", err)
			return nil, nil, err
		}
		return senWallet, recWallet, nil
	}

	recWallet, err := wt.walletRepo.GetUserWallet(ctx, transfer.ReceiverWalletID, transfer.ReceiverID)
	if err != nil {
		app.Logger.Errorf(ctx, "[WalletTransferer-getSenderAndReceiverWallet] sender >= receiver; get receiver wallet fail: %v", err)
		return nil, nil, err
	}
	senWallet, err := wt.walletRepo.GetUserWallet(ctx, transfer.SenderWalletID, transfer.SenderID)
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
