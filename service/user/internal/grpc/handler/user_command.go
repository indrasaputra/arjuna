package handler

import (
	"context"

	"google.golang.org/grpc/metadata"

	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/app"
	"github.com/indrasaputra/arjuna/service/user/internal/service"
)

const (
	headerIdempotencyKey = "x-idempotency-key"
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
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, entity.ErrInternal("metadata not found from incoming context")
	}
	key := md[headerIdempotencyKey]
	if len(key) == 0 {
		return nil, entity.ErrMissingIdempotencyKey()
	}

	if request == nil || request.GetUser() == nil {
		app.Logger.Errorf(ctx, "[UserCommand-RegisterUser] empty or nil user")
		return nil, entity.ErrEmptyUser()
	}

	id, err := uc.registrar.Register(ctx, createUserFromRegisterUserRequest(request), key[0])
	if err != nil {
		app.Logger.Errorf(ctx, "[UserCommand-RegisterUser] fail register user: %v", err)
		return nil, err
	}
	return &apiv1.RegisterUserResponse{Data: &apiv1.User{Id: id}}, nil
}

func createUserFromRegisterUserRequest(request *apiv1.RegisterUserRequest) *entity.User {
	return &entity.User{
		Name:     request.GetUser().GetName(),
		Email:    request.GetUser().GetEmail(),
		Password: request.GetUser().GetPassword(),
	}
}
