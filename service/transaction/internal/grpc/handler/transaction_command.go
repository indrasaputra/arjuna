package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc/metadata"

	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
	"github.com/indrasaputra/arjuna/service/transaction/entity"
	"github.com/indrasaputra/arjuna/service/transaction/internal/app"
	"github.com/indrasaputra/arjuna/service/transaction/internal/service"
)

const (
	headerIdempotencyKey = "x-idempotency-key"
)

// TransactionCommand handles HTTP/2 gRPC request for state-changing transaction.
type TransactionCommand struct {
	apiv1.UnimplementedTransactionCommandServiceServer
	creator service.CreateTransaction
	deleter service.DeleteTransaction
}

// NewTransactionCommand creates an instance of TransactionCommand.
func NewTransactionCommand(c service.CreateTransaction, d service.DeleteTransaction) *TransactionCommand {
	return &TransactionCommand{creator: c, deleter: d}
}

// CreateTransaction handles HTTP/2 gRPC request similar to POST in HTTP/1.1.
func (tc *TransactionCommand) CreateTransaction(ctx context.Context, request *apiv1.CreateTransactionRequest) (*apiv1.CreateTransactionResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, entity.ErrInternal("metadata not found from incoming context")
	}
	key := md[headerIdempotencyKey]
	if len(key) == 0 {
		return nil, entity.ErrMissingIdempotencyKey()
	}

	if request == nil || request.GetTransaction() == nil {
		app.Logger.Errorf(ctx, "[TransactionCommand-CreateTransaction] empty or nil transaction")
		return nil, entity.ErrEmptyTransaction()
	}

	amount, _ := decimal.NewFromString(request.GetTransaction().GetAmount())
	id, err := tc.creator.Create(ctx, createTransactionFromCreateTransactionRequest(request, amount), key[0])
	if err != nil {
		app.Logger.Errorf(ctx, "[TransactionCommand-CreateTransaction] fail register transaction: %v", err)
		return nil, err
	}
	return &apiv1.CreateTransactionResponse{Data: &apiv1.Transaction{Id: id.String()}}, nil
}

// DeleteAllTransactions handles HTTP/2 gRPC request similar to POST in HTTP/1.1.
func (tc *TransactionCommand) DeleteAllTransactions(ctx context.Context, _ *apiv1.DeleteAllTransactionsRequest) (*apiv1.DeleteAllTransactionsResponse, error) {
	err := tc.deleter.DeleteAllTransactions(ctx)
	if err != nil {
		app.Logger.Errorf(ctx, "[TransactionCommand-DeleteAllTransactions] delete all transactions fail: %v", err)
		return nil, err
	}
	return &apiv1.DeleteAllTransactionsResponse{}, nil
}

func createTransactionFromCreateTransactionRequest(request *apiv1.CreateTransactionRequest, amount decimal.Decimal) *entity.Transaction {
	return &entity.Transaction{
		SenderID:   uuid.MustParse(request.GetTransaction().GetSenderId()),
		ReceiverID: uuid.MustParse(request.GetTransaction().GetReceiverId()),
		Amount:     amount,
	}
}
