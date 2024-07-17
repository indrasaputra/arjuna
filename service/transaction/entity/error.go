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
	te := &apiv1.TransactionError{
		ErrorCode: apiv1.TransactionErrorCode_TRANSACTION_ERROR_CODE_INTERNAL,
	}
	res, err := st.WithDetails(te)
	if err != nil {
		return st.Err()
	}
	return res.Err()
}

// ErrAlreadyExists returns codes.AlreadyExists explained that the transaction already exists.
func ErrAlreadyExists() error {
	st := status.New(codes.AlreadyExists, "")
	te := &apiv1.TransactionError{
		ErrorCode: apiv1.TransactionErrorCode_TRANSACTION_ERROR_CODE_ALREADY_EXISTS,
	}
	res, err := st.WithDetails(te)
	if err != nil {
		return st.Err()
	}
	return res.Err()
}

// ErrEmptyTransaction returns codes.InvalidArgument explained that the instance is empty or nil.
func ErrEmptyTransaction() error {
	st := status.New(codes.InvalidArgument, "")
	br := createBadRequest(&errdetails.BadRequest_FieldViolation{
		Field:       "transaction instance",
		Description: "empty or nil",
	})

	te := &apiv1.TransactionError{
		ErrorCode: apiv1.TransactionErrorCode_TRANSACTION_ERROR_CODE_EMPTY_TRANSACTION,
	}
	res, err := st.WithDetails(br, te)
	if err != nil {
		return st.Err()
	}
	return res.Err()
}

// ErrInvalidSender returns codes.InvalidArgument explained that the sender is invalid.
func ErrInvalidSender() error {
	st := status.New(codes.InvalidArgument, "")
	br := createBadRequest(&errdetails.BadRequest_FieldViolation{
		Field:       "sender_id",
		Description: "empty",
	})

	te := &apiv1.TransactionError{
		ErrorCode: apiv1.TransactionErrorCode_TRANSACTION_ERROR_CODE_INVALID_SENDER,
	}
	res, err := st.WithDetails(br, te)
	if err != nil {
		return st.Err()
	}
	return res.Err()
}

// ErrInvalidReceiver returns codes.InvalidArgument explained that the receiver is invalid.
func ErrInvalidReceiver() error {
	st := status.New(codes.InvalidArgument, "")
	br := createBadRequest(&errdetails.BadRequest_FieldViolation{
		Field:       "receiver_id",
		Description: "empty",
	})

	te := &apiv1.TransactionError{
		ErrorCode: apiv1.TransactionErrorCode_TRANSACTION_ERROR_CODE_INVALID_RECEIVER,
	}
	res, err := st.WithDetails(br, te)
	if err != nil {
		return st.Err()
	}
	return res.Err()
}

// ErrInvalidAmount returns codes.InvalidArgument explained that the amount is invalid.
func ErrInvalidAmount() error {
	st := status.New(codes.InvalidArgument, "")
	br := createBadRequest(&errdetails.BadRequest_FieldViolation{
		Field:       "amount",
		Description: "less than or equal to zero",
	})

	te := &apiv1.TransactionError{
		ErrorCode: apiv1.TransactionErrorCode_TRANSACTION_ERROR_CODE_INVALID_AMOUNT,
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
	te := &apiv1.TransactionError{
		ErrorCode: apiv1.TransactionErrorCode_TRANSACTION_ERROR_CODE_MISSING_IDEMPOTENCY_KEY,
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
