package wallet

import (
	"context"
	"log/slog"

	"github.com/gogo/status"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc/codes"

	"github.com/indrasaputra/arjuna/service/user/entity"
	enwallet "github.com/indrasaputra/arjuna/service/wallet/entity"
	sdkwallet "github.com/indrasaputra/arjuna/service/wallet/pkg/sdk/wallet"
)

// Wallet is responsible to connect to wallet service.
type Wallet struct {
	client *sdkwallet.Client
}

// NewWallet creates an instance of Wallet.
func NewWallet(c *sdkwallet.Client) *Wallet {
	return &Wallet{client: c}
}

// CreateWallet creates an account.
func (a *Wallet) CreateWallet(ctx context.Context, user *entity.User) error {
	req := &enwallet.Wallet{UserID: user.ID, Balance: decimal.Zero}
	err := a.client.CreateWallet(ctx, req)
	if err != nil {
		slog.ErrorContext(ctx, "[Wallet-CreateWallet] fail call register", "error", err)
	}
	if status.Code(err) == codes.AlreadyExists {
		return entity.ErrAlreadyExists()
	}
	return err
}
