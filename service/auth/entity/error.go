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
	te := &apiv1.AuthError{
		ErrorCode: apiv1.AuthErrorCode_AUTH_ERROR_CODE_INTERNAL,
	}
	res, err := st.WithDetails(te)
	if err != nil {
		return st.Err()
	}
	return res.Err()
}

// ErrEmptyField returns codes.InvalidArgument explained that the field is empty.
func ErrEmptyField(field string) error {
	st := status.New(codes.InvalidArgument, "")
	br := createBadRequest(&errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: "empty or nil",
	})

	te := &apiv1.AuthError{
		ErrorCode: apiv1.AuthErrorCode_AUTH_ERROR_CODE_EMPTY_FIELD,
	}
	res, err := st.WithDetails(br, te)
	if err != nil {
		return st.Err()
	}
	return res.Err()
}

// ErrUnauthorized returns codes.Unauthenticated explained that credential might be wrong.
func ErrUnauthorized() error {
	st := status.New(codes.Unauthenticated, "")
	te := &apiv1.AuthError{
		ErrorCode: apiv1.AuthErrorCode_AUTH_ERROR_CODE_UNAUTHORIZED,
	}
	res, err := st.WithDetails(te)
	if err != nil {
		return st.Err()
	}
	return res.Err()
}

// ErrInvalidArgument returns codes.InvalidArgument explained that argument(s) is invalid.
func ErrInvalidArgument(message string) error {
	st := status.New(codes.InvalidArgument, message)
	te := &apiv1.AuthError{
		ErrorCode: apiv1.AuthErrorCode_AUTH_ERROR_CODE_INVALID_ARGUMENT,
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
