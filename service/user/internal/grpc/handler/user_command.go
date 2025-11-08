package handler

import (
	"context"
	"log/slog"

	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/service"
)

// UserCommand handles HTTP/2 gRPC request for state-changing user.
type UserCommand struct {
	apiv1.UnimplementedUserCommandServiceServer
	registrar service.RegisterUser
}

// NewUserCommand creates an instance of UserCommand.
func NewUserCommand(registrar service.RegisterUser) *UserCommand {
	return &UserCommand{registrar: registrar}
}

// RegisterUser handles HTTP/2 gRPC request similar to POST in HTTP/1.1.
func (uc *UserCommand) RegisterUser(ctx context.Context, request *apiv1.RegisterUserRequest) (*apiv1.RegisterUserResponse, error) {
	if request == nil || request.GetUser() == nil {
		slog.ErrorContext(ctx, "[UserCommand-RegisterUser] empty or nil user")
		return nil, entity.ErrEmptyUser()
	}

	id, err := uc.registrar.Register(ctx, createUserFromRegisterUserRequest(request))
	if err != nil {
		slog.ErrorContext(ctx, "[UserCommand-RegisterUser] fail register user", "error", err)
		return nil, err
	}
	return &apiv1.RegisterUserResponse{Data: &apiv1.User{Id: id.String()}}, nil
}

func createUserFromRegisterUserRequest(request *apiv1.RegisterUserRequest) *entity.User {
	return &entity.User{
		Name:     request.GetUser().GetName(),
		Email:    request.GetUser().GetEmail(),
		Password: request.GetUser().GetPassword(),
	}
}
