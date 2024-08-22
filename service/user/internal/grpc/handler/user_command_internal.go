package handler

import (
	"context"

	"github.com/google/uuid"

	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/app"
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
		app.Logger.Errorf(ctx, "[UserCommand-DeleteUser] empty or nil user")
		return nil, entity.ErrEmptyUser()
	}

	err := uci.deleter.HardDelete(ctx, uuid.MustParse(request.GetId()))
	if err != nil {
		app.Logger.Errorf(ctx, "[UserCommand-DeleteUser] fail delete user: %v", err)
		return nil, err
	}
	return &apiv1.DeleteUserResponse{}, nil
}
