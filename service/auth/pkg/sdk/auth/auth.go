package auth

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
	"github.com/indrasaputra/arjuna/service/auth/entity"
)

// DialConfig defines configuration to work with Client.
type DialConfig struct {
	// Host defines server host.
	Host string
	// Options defines list of dial option used to make a connection to server.
	Options []grpc.DialOption
}

// Client is responsible to connect to auth use cases.
type Client struct {
	handler apiv1.AuthServiceClient
}

// NewClient creates an instance of Client.
func NewClient(cfg *DialConfig) (*Client, error) {
	conn, err := grpc.NewClient(cfg.Host, cfg.Options...)
	if err != nil {
		return nil, status.New(codes.Unavailable, "").Err()
	}

	return &Client{
		handler: apiv1.NewAuthServiceClient(conn),
	}, nil
}

// Register registers an account.
func (c *Client) Register(ctx context.Context, account *entity.Account) error {
	req := &apiv1.RegisterAccountRequest{Account: &apiv1.Account{
		UserId:   account.UserID,
		Email:    account.Email,
		Password: account.Password,
	}}

	_, err := c.handler.RegisterAccount(ctx, req)
	return err
}

// ParseToken parses the token.
func ParseToken(tokenString string, secret []byte) (*entity.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &entity.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, entity.ErrInvalidArgument("unexpected signing method")
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*entity.Claims); ok {
		return claims, nil
	}
	return nil, entity.ErrInternal("unknown claims type")
}
