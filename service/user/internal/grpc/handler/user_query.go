package handler

import (
	"context"
	"log/slog"

	"google.golang.org/protobuf/types/known/timestamppb"

	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/service"
)

// UserQuery handles HTTP/2 gRPC request for retrieving user.
type UserQuery struct {
	apiv1.UnimplementedUserQueryServiceServer
	getter service.GetUser
}

// NewUserQuery creates an instance of UserQuery.
func NewUserQuery(getter service.GetUser) *UserQuery {
	return &UserQuery{getter: getter}
}

// GetAllUsers handles HTTP/2 gRPC request similar to GET in HTTP/1.1.
func (uc *UserQuery) GetAllUsers(ctx context.Context, request *apiv1.GetAllUsersRequest) (*apiv1.GetAllUsersResponse, error) {
	if request == nil {
		return nil, entity.ErrEmptyUser()
	}

	users, err := uc.getter.GetAll(ctx, uint(request.GetLimit()))
	if err != nil {
		slog.ErrorContext(ctx, "[UserQuery-GetAllUsers] fail get all users", "error", err)
		return nil, err
	}
	return createGetAllUsersResponse(users), nil
}

func createGetAllUsersResponse(users []*entity.User) *apiv1.GetAllUsersResponse {
	resp := &apiv1.GetAllUsersResponse{}
	for _, user := range users {
		resp.Data = append(resp.Data, createProtoUser(user))
	}
	return resp
}

func createProtoUser(user *entity.User) *apiv1.User {
	return &apiv1.User{
		Id:        user.ID.String(),
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}
