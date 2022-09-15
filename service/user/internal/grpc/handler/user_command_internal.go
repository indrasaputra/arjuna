package handler

import (
	"context"

	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/service"
)

// UserCommandInternal handles HTTP/2 gRPC request for state-changing user.
type UserCommandInternal struct {
	apiv1.UnimplementedUserCommandInternalServiceServer
	deleter service.DeleteUser
}

// NewUserCommandInternal creates an instance of UserCommandInternal.
func NewUserCommandInternal(deleter service.DeleteUser) *UserCommandInternal {
	return &UserCommandInternal{deleter: deleter}
}

// DeleteUser handles HTTP/2 gRPC request similar to POST in HTTP/1.1.
func (uci *UserCommandInternal) DeleteUser(ctx context.Context, request *apiv1.DeleteUserRequest) (*apiv1.DeleteUserResponse, error) {
	if request == nil || request.GetId() == "" {
		return nil, entity.ErrEmptyUser()
	}

	err := uci.deleter.HardDelete(ctx, request.GetId())
	if err != nil {
		return nil, err
	}
	return &apiv1.DeleteUserResponse{}, nil
}
