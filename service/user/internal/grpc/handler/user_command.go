package handler

import (
	"context"

	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/service"
)

// UserCommand handles HTTP/2 gRPC request for state-changing user .
type UserCommand struct {
	apiv1.UnimplementedUserCommandServiceServer
	registrator service.RegisterUser
}

// NewUserCommand creates an instance of UserCommand.
func NewUserCommand(registrator service.RegisterUser) *UserCommand {
	return &UserCommand{registrator: registrator}
}

// RegisterUser handles HTTP/2 gRPC request similar to POST in HTTP/1.1.
func (uc *UserCommand) RegisterUser(ctx context.Context, request *apiv1.RegisterUserRequest) (*apiv1.RegisterUserResponse, error) {
	if request == nil || request.GetUser() == nil {
		return nil, entity.ErrEmptyUser()
	}

	_, err := uc.registrator.Register(ctx, createUserFromRegisterUserRequest(request))
	if err != nil {
		return nil, err
	}
	return &apiv1.RegisterUserResponse{}, nil
}

func createUserFromRegisterUserRequest(request *apiv1.RegisterUserRequest) *entity.User {
	return &entity.User{
		Name:     request.GetUser().GetName(),
		Email:    request.GetUser().GetEmail(),
		Password: request.GetUser().GetPassword(),
	}
}
