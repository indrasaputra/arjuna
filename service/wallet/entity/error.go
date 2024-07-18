package entity

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
)

// ErrInternal returns codes.Internal explained that unexpected behavior occurred in system.
func ErrInternal(message string) error {
	st := status.New(codes.Internal, message)
	te := &apiv1.WalletError{
		ErrorCode: apiv1.WalletErrorCode_WALLET_ERROR_CODE_INTERNAL,
	}
	res, err := st.WithDetails(te)
	if err != nil {
		return st.Err()
	}
	return res.Err()
}

// ErrAlreadyExists returns codes.AlreadyExists explained that the wallet already exists.
func ErrAlreadyExists() error {
	st := status.New(codes.AlreadyExists, "")
	te := &apiv1.WalletError{
		ErrorCode: apiv1.WalletErrorCode_WALLET_ERROR_CODE_ALREADY_EXISTS,
	}
	res, err := st.WithDetails(te)
	if err != nil {
		return st.Err()
	}
	return res.Err()
}

// ErrEmptyWallet returns codes.InvalidArgument explained that the instance is empty or nil.
func ErrEmptyWallet() error {
	st := status.New(codes.InvalidArgument, "")
	br := createBadRequest(&errdetails.BadRequest_FieldViolation{
		Field:       "wallet instance",
		Description: "empty or nil",
	})

	te := &apiv1.WalletError{
		ErrorCode: apiv1.WalletErrorCode_WALLET_ERROR_CODE_EMPTY_WALLET,
	}
	res, err := st.WithDetails(br, te)
	if err != nil {
		return st.Err()
	}
	return res.Err()
}

// ErrInvalidBalance returns codes.InvalidArgument explained that the balance invalid.
func ErrInvalidBalance() error {
	st := status.New(codes.InvalidArgument, "")
	br := createBadRequest(&errdetails.BadRequest_FieldViolation{
		Field:       "balance",
		Description: "must be numeric and greater than or equal to zero",
	})

	te := &apiv1.WalletError{
		ErrorCode: apiv1.WalletErrorCode_WALLET_ERROR_CODE_INVALID_BALANCE,
	}
	res, err := st.WithDetails(br, te)
	if err != nil {
		return st.Err()
	}
	return res.Err()
}

// ErrInvalidUser returns codes.InvalidArgument explained that the user is invalid.
func ErrInvalidUser() error {
	st := status.New(codes.InvalidArgument, "")
	br := createBadRequest(&errdetails.BadRequest_FieldViolation{
		Field:       "user_id",
		Description: "empty",
	})

	te := &apiv1.WalletError{
		ErrorCode: apiv1.WalletErrorCode_WALLET_ERROR_CODE_INVALID_USER,
	}
	res, err := st.WithDetails(br, te)
	if err != nil {
		return st.Err()
	}
	return res.Err()
}

// ErrMissingIdempotencyKey returns codes.InvalidArgument explained that the idempotency key is missing.
func ErrMissingIdempotencyKey() error {
	st := status.New(codes.InvalidArgument, "missing idempotency key")
	te := &apiv1.WalletError{
		ErrorCode: apiv1.WalletErrorCode_WALLET_ERROR_CODE_MISSING_IDEMPOTENCY_KEY,
	}
	res, err := st.WithDetails(te)
	if err != nil {
		return st.Err()
	}
	return res.Err()
}

func createBadRequest(details ...*errdetails.BadRequest_FieldViolation) *errdetails.BadRequest {
	return &errdetails.BadRequest{
		FieldViolations: details,
	}
}
