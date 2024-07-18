package handler

import (
	"context"

	"github.com/shopspring/decimal"
	"google.golang.org/grpc/metadata"

	"github.com/indrasaputra/arjuna/pkg/sdk/grpc/interceptor"
	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
	"github.com/indrasaputra/arjuna/service/wallet/entity"
	"github.com/indrasaputra/arjuna/service/wallet/internal/app"
	"github.com/indrasaputra/arjuna/service/wallet/internal/service"
)

const (
	headerIdempotencyKey = "x-idempotency-key"
)

// WalletCommand handles HTTP/2 gRPC request for state-changing wallet.
type WalletCommand struct {
	apiv1.UnimplementedWalletCommandServiceServer
	creator service.CreateWallet
	topup   service.TopupWallet
}

// NewWalletCommand creates an instance of WalletCommand.
func NewWalletCommand(c service.CreateWallet, t service.TopupWallet) *WalletCommand {
	return &WalletCommand{creator: c, topup: t}
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

// TopupWallet handles HTTP/2 gRPC request similar to POST in HTTP/1.1.
func (wc *WalletCommand) TopupWallet(ctx context.Context, request *apiv1.TopupWalletRequest) (*apiv1.TopupWalletResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, entity.ErrInternal("metadata not found from incoming context")
	}
	key := md[headerIdempotencyKey]
	if len(key) == 0 {
		return nil, entity.ErrMissingIdempotencyKey()
	}

	userID := ctx.Value(interceptor.HeaderKeyUserID).(interceptor.HeaderKey)

	if request == nil || request.GetTopup() == nil {
		app.Logger.Errorf(ctx, "[WalletCommand-TopupWallet] empty or nil topup")
		return nil, entity.ErrEmptyWallet()
	}

	amount, _ := decimal.NewFromString(request.GetTopup().GetAmount())
	req := createTopupWalletFromTopupWalletRequest(request, string(userID), amount)

	err := wc.topup.Topup(ctx, req)
	if err != nil {
		app.Logger.Errorf(ctx, "[WalletCommand-TopupWallet] fail topup wallet: %v", err)
		return nil, err
	}
	return &apiv1.TopupWalletResponse{}, nil
}

func createWalletFromCreateWalletRequest(request *apiv1.CreateWalletRequest, balance decimal.Decimal) *entity.Wallet {
	return &entity.Wallet{
		UserID:  request.GetWallet().GetUserId(),
		Balance: balance,
	}
}

func createTopupWalletFromTopupWalletRequest(request *apiv1.TopupWalletRequest, userID string, amount decimal.Decimal) *entity.TopupWallet {
	return &entity.TopupWallet{
		WalletID: request.GetTopup().GetWalletId(),
		UserID:   userID,
		Amount:   amount,
	}
}
