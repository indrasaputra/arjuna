package interceptor

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	grpcauth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdkauth "github.com/indrasaputra/arjuna/service/auth/pkg/sdk/auth"
)

const (
	bearer = "bearer"
	// HeaderKeyUserID contains user's ID.
	HeaderKeyUserID = HeaderKey("X-User-ID")
	// HeaderKeyEmail contains user's email.
	HeaderKeyEmail = HeaderKey("X-User-Email")
)

// HeaderKey represents a string for request header key.
type HeaderKey string

// AuthBearer intercepts the request and check for bearer authorization.
// If success, it will inject the claims to context.
func AuthBearer(secret []byte) func(context.Context) (context.Context, error) {
	return func(ctx context.Context) (context.Context, error) {
		token, err := grpcauth.AuthFromMD(ctx, bearer)
		if err != nil {
			return ctx, status.Error(codes.Unauthenticated, "unauthenticated")
		}

		claims, err := sdkauth.ParseToken(token, secret)
		if err != nil {
			return ctx, status.Error(codes.Unauthenticated, "unauthenticated")
		}

		ctx = context.WithValue(ctx, HeaderKeyUserID, claims.UserID)
		ctx = context.WithValue(ctx, HeaderKeyEmail, claims.Email)
		return ctx, nil
	}
}

// SkipMethod skips the interceptor for the given methods.
func SkipMethod(methods ...string) func(context.Context, interceptors.CallMeta) bool {
	return func(_ context.Context, c interceptors.CallMeta) bool {
		for _, m := range methods {
			if c.FullMethod() == m {
				return false
			}
		}
		return true
	}
}
