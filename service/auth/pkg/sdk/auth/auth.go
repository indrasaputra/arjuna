package auth

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
	"github.com/indrasaputra/arjuna/service/auth/entity"
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

// Client is responsible to connect to auth use cases.
type Client struct {
	handler apiv1.AuthServiceClient
	config  *Config
}

// NewClient creates an instance of Client.
func NewClient(cfg *Config) (*Client, error) {
	conn, err := grpc.NewClient(cfg.Host, cfg.Options...)
	if err != nil {
		return nil, status.New(codes.Unavailable, "").Err()
	}

	return &Client{
		handler: apiv1.NewAuthServiceClient(conn),
		config:  cfg,
	}, nil
}

// Register registers an account.
func (c *Client) Register(ctx context.Context, account *entity.Account) error {
	req := &apiv1.RegisterAccountRequest{Account: &apiv1.Account{
		UserId:   account.UserID.String(),
		Email:    account.Email,
		Password: account.Password,
	}}

	token := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.config.Username, c.config.Password)))
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(headerAuthorization, fmt.Sprintf("basic %s", token)))

	_, err := c.handler.RegisterAccount(ctx, req)
	return err
}

// ParseToken parses the token.
func ParseToken(tokenString string, secret []byte) (*entity.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &entity.Claims{}, func(token *jwt.Token) (any, error) {
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
