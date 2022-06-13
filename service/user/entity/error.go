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
	te := &apiv1.UserError{
		ErrorCode: apiv1.UserErrorCode_USER_ERROR_CODE_INTERNAL,
	}
	res, err := st.WithDetails(te)
	if err != nil {
		return st.Err()
	}
	return res.Err()
}

// ErrEmptyUser returns codes.InvalidArgument explained that the instance is empty or nil.
func ErrEmptyUser() error {
	st := status.New(codes.InvalidArgument, "")
	br := createBadRequest(&errdetails.BadRequest_FieldViolation{
		Field:       "user instance",
		Description: "empty or nil",
	})

	te := &apiv1.UserError{
		ErrorCode: apiv1.UserErrorCode_USER_ERROR_CODE_EMPTY_USER,
	}
	res, err := st.WithDetails(br, te)
	if err != nil {
		return st.Err()
	}
	return res.Err()
}

// ErrAlreadyExists returns codes.AlreadyExists explained that the user already exists.
func ErrAlreadyExists() error {
	st := status.New(codes.AlreadyExists, "")
	te := &apiv1.UserError{
		ErrorCode: apiv1.UserErrorCode_USER_ERROR_CODE_ALREADY_EXISTS,
	}
	res, err := st.WithDetails(te)
	if err != nil {
		return st.Err()
	}
	return res.Err()
}

// ErrInvalidName returns codes.InvalidArgument explained that the user's name is invalid.
func ErrInvalidName() error {
	st := status.New(codes.InvalidArgument, "")
	br := createBadRequest(&errdetails.BadRequest_FieldViolation{
		Field:       "name",
		Description: "contain character outside of alphabet",
	})

	te := &apiv1.UserError{
		ErrorCode: apiv1.UserErrorCode_USER_ERROR_CODE_INVALID_NAME,
	}
	res, err := st.WithDetails(br, te)
	if err != nil {
		return st.Err()
	}
	return res.Err()
}

// ErrInvalidEmail returns codes.InvalidArgument explained that the user's email is invalid.
func ErrInvalidEmail() error {
	st := status.New(codes.InvalidArgument, "")
	br := createBadRequest(&errdetails.BadRequest_FieldViolation{
		Field: "email",
	})

	te := &apiv1.UserError{
		ErrorCode: apiv1.UserErrorCode_USER_ERROR_CODE_INVALID_EMAIL,
	}
	res, err := st.WithDetails(br, te)
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
