package wallet

import (
	"context"
	"encoding/base64"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
	"github.com/indrasaputra/arjuna/service/wallet/entity"
)

const (
	headerAuthorization = "authorization"
)

// Config defines configuration to work with Client.
type Config struct {
	Host     string
	Username string
	Password string
	Options  []grpc.DialOption
}

// Client is responsible to connect to wallet use cases.
type Client struct {
	handler apiv1.WalletCommandServiceClient
	config  *Config
}

// NewClient creates an instance of Client.
func NewClient(cfg *Config) (*Client, error) {
	conn, err := grpc.NewClient(cfg.Host, cfg.Options...)
	if err != nil {
		return nil, status.New(codes.Unavailable, "").Err()
	}

	return &Client{
		handler: apiv1.NewWalletCommandServiceClient(conn),
		config:  cfg,
	}, nil
}

// CreateWallet creates a wallet.
func (c *Client) CreateWallet(ctx context.Context, wallet *entity.Wallet) error {
	req := &apiv1.CreateWalletRequest{Wallet: &apiv1.Wallet{
		UserId:  wallet.UserID.String(),
		Balance: wallet.Balance.String(),
	}}

	token := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.config.Username, c.config.Password)))
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(headerAuthorization, fmt.Sprintf("basic %s", token)))

	_, err := c.handler.CreateWallet(ctx, req)
	return err
}
