package handler

import (
	"context"

	"github.com/shopspring/decimal"

	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
	"github.com/indrasaputra/arjuna/service/wallet/entity"
	"github.com/indrasaputra/arjuna/service/wallet/internal/app"
	"github.com/indrasaputra/arjuna/service/wallet/internal/service"
)

// WalletCommand handles HTTP/2 gRPC request for state-changing wallet.
type WalletCommand struct {
	apiv1.UnimplementedWalletCommandServiceServer
	creator service.CreateWallet
}

// NewWalletCommand creates an instance of WalletCommand.
func NewWalletCommand(creator service.CreateWallet) *WalletCommand {
	return &WalletCommand{creator: creator}
}

// CreateWallet handles HTTP/2 gRPC request similar to POST in HTTP/1.1.
func (wc *WalletCommand) CreateWallet(ctx context.Context, request *apiv1.CreateWalletRequest) (*apiv1.CreateWalletResponse, error) {
	if request == nil || request.GetWallet() == nil {
		app.Logger.Errorf(ctx, "[WalletCommand-CreateWallet] empty or nil wallet")
		return nil, entity.ErrEmptyWallet()
	}

	balance, _ := decimal.NewFromString(request.GetWallet().GetBalance())
	req := createWalletFromCreateWalletRequest(request, balance)

	err := wc.creator.Create(ctx, req)
	if err != nil {
		app.Logger.Errorf(ctx, "[WalletCommand-CreateWallet] fail register wallet: %v", err)
		return nil, err
	}
	return &apiv1.CreateWalletResponse{}, nil
}

func createWalletFromCreateWalletRequest(request *apiv1.CreateWalletRequest, balance decimal.Decimal) *entity.Wallet {
	return &entity.Wallet{
		UserID:  request.GetWallet().GetUserId(),
		Balance: balance,
	}
}
