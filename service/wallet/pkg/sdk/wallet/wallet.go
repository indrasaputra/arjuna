package wallet

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
	"github.com/indrasaputra/arjuna/service/wallet/entity"
)

// DialConfig defines configuration to work with Client.
type DialConfig struct {
	// Host defines server host.
	Host string
	// Options defines list of dial option used to make a connection to server.
	Options []grpc.DialOption
}

// Client is responsible to connect to wallet use cases.
type Client struct {
	handler apiv1.WalletCommandServiceClient
}

// NewClient creates an instance of Client.
func NewClient(cfg *DialConfig) (*Client, error) {
	conn, err := grpc.NewClient(cfg.Host, cfg.Options...)
	if err != nil {
		return nil, status.New(codes.Unavailable, "").Err()
	}

	return &Client{
		handler: apiv1.NewWalletCommandServiceClient(conn),
	}, nil
}

// CreateWallet creates a wallet.
func (c *Client) CreateWallet(ctx context.Context, wallet *entity.Wallet) error {
	req := &apiv1.CreateWalletRequest{Wallet: &apiv1.Wallet{
		UserId:  wallet.UserID,
		Balance: wallet.Balance.String(),
	}}

	_, err := c.handler.CreateWallet(ctx, req)
	return err
}
