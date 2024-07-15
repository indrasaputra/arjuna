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

// ErrInvalidEmail returns codes.InvalidArgument explained that email is invalid.
func ErrInvalidEmail() error {
	st := status.New(codes.InvalidArgument, "email is invalid")
	te := &apiv1.AuthError{
		ErrorCode: apiv1.AuthErrorCode_AUTH_ERROR_CODE_INVALID_EMAIL,
	}
	res, err := st.WithDetails(te)
	if err != nil {
		return st.Err()
	}
	return res.Err()
}

// ErrInvalidPassword returns codes.InvalidArgument explained that password is invalid.
func ErrInvalidPassword() error {
	st := status.New(codes.InvalidArgument, "password is invalid")
	te := &apiv1.AuthError{
		ErrorCode: apiv1.AuthErrorCode_AUTH_ERROR_CODE_INVALID_EMAIL,
	}
	res, err := st.WithDetails(te)
	if err != nil {
		return st.Err()
	}
	return res.Err()
}

// ErrEmptyAccount returns codes.InvalidArgument explained that the account is empty.
func ErrEmptyAccount() error {
	st := status.New(codes.InvalidArgument, "empty account")
	te := &apiv1.AuthError{
		ErrorCode: apiv1.AuthErrorCode_AUTH_ERROR_CODE_EMPTY_ACCOUNT,
	}
	res, err := st.WithDetails(te)
	if err != nil {
		return st.Err()
	}
	return res.Err()
}

// ErrAlreadyExists returns codes.AlreadyExists explained that the account already exists.
func ErrAlreadyExists() error {
	st := status.New(codes.AlreadyExists, "")
	te := &apiv1.AuthError{
		ErrorCode: apiv1.AuthErrorCode_AUTH_ERROR_CODE_ALREADY_EXISTS,
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
