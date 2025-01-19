package handler

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc/metadata"

	"github.com/indrasaputra/arjuna/pkg/sdk/grpc/interceptor"
	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
	"github.com/indrasaputra/arjuna/service/wallet/entity"
	"github.com/indrasaputra/arjuna/service/wallet/internal/service"
)

const (
	headerIdempotencyKey = "x-idempotency-key"
)

// WalletCommand handles HTTP/2 gRPC request for state-changing wallet.
type WalletCommand struct {
	apiv1.UnimplementedWalletCommandServiceServer
	creator  service.CreateWallet
	topup    service.TopupWallet
	transfer service.TransferWallet
}

// NewWalletCommand creates an instance of WalletCommand.
func NewWalletCommand(c service.CreateWallet, t service.TopupWallet, tf service.TransferWallet) *WalletCommand {
	return &WalletCommand{creator: c, topup: t, transfer: tf}
}

// CreateWallet handles HTTP/2 gRPC request similar to POST in HTTP/1.1.
func (wc *WalletCommand) CreateWallet(ctx context.Context, request *apiv1.CreateWalletRequest) (*apiv1.CreateWalletResponse, error) {
	if request == nil || request.GetWallet() == nil {
		slog.ErrorContext(ctx, "[WalletCommand-CreateWallet] empty or nil wallet")
		return nil, entity.ErrEmptyWallet()
	}

	balance, _ := decimal.NewFromString(request.GetWallet().GetBalance())
	req := createWalletFromCreateWalletRequest(request, balance)

	err := wc.creator.Create(ctx, req)
	if err != nil {
		slog.ErrorContext(ctx, "[WalletCommand-CreateWallet] fail register wallet", "error", err)
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

	userID := ctx.Value(interceptor.HeaderKeyUserID).(uuid.UUID)

	if request == nil || request.GetTopup() == nil {
		slog.ErrorContext(ctx, "[WalletCommand-TopupWallet] empty or nil topup")
		return nil, entity.ErrEmptyWallet()
	}

	amount, _ := decimal.NewFromString(request.GetTopup().GetAmount())
	req := createTopupWalletFromTopupWalletRequest(request, userID, amount, key[0])

	err := wc.topup.Topup(ctx, req)
	if err != nil {
		slog.ErrorContext(ctx, "[WalletCommand-TopupWallet] fail topup wallet", "error", err)
		return nil, err
	}
	return &apiv1.TopupWalletResponse{}, nil
}

// TransferBalance handles HTTP/2 gRPC request similar to POST in HTTP/1.1.
func (wc *WalletCommand) TransferBalance(ctx context.Context, request *apiv1.TransferBalanceRequest) (*apiv1.TransferBalanceResponse, error) {
	if request == nil || request.GetTransfer() == nil {
		slog.ErrorContext(ctx, "[WalletCommand-TransferBalance] empty or nil transfer")
		return nil, entity.ErrEmptyWallet()
	}

	amount, _ := decimal.NewFromString(request.GetTransfer().GetAmount())
	req := createTransferWalletFromTransferBalanceRequest(request, amount)

	err := wc.transfer.TransferBalance(ctx, req)
	if err != nil {
		slog.ErrorContext(ctx, "[WalletCommand-TransferBalance] fail transfer wallet", "error", err)
		return nil, err
	}
	return &apiv1.TransferBalanceResponse{}, nil
}

func createWalletFromCreateWalletRequest(request *apiv1.CreateWalletRequest, balance decimal.Decimal) *entity.Wallet {
	return &entity.Wallet{
		UserID:  uuid.MustParse(request.GetWallet().GetUserId()),
		Balance: balance,
	}
}

func createTopupWalletFromTopupWalletRequest(request *apiv1.TopupWalletRequest, userID uuid.UUID, amount decimal.Decimal, key string) *entity.TopupWallet {
	return &entity.TopupWallet{
		WalletID:       uuid.MustParse(request.GetTopup().GetWalletId()),
		UserID:         userID,
		Amount:         amount,
		IdempotencyKey: key,
	}
}

func createTransferWalletFromTransferBalanceRequest(request *apiv1.TransferBalanceRequest, amount decimal.Decimal) *entity.TransferWallet {
	return &entity.TransferWallet{
		SenderID:         uuid.MustParse(request.GetTransfer().GetSenderId()),
		SenderWalletID:   uuid.MustParse(request.GetTransfer().GetSenderWalletId()),
		ReceiverID:       uuid.MustParse(request.GetTransfer().GetReceiverId()),
		ReceiverWalletID: uuid.MustParse(request.GetTransfer().GetReceiverWalletId()),
		Amount:           amount,
	}
}
