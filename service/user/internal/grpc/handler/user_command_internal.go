package handler

// import (
// 	"context"

// 	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
// 	"github.com/indrasaputra/arjuna/service/user/entity"
// 	"github.com/indrasaputra/arjuna/service/user/internal/service"
// )

// // UserCommandInternal handles HTTP/2 gRPC request for state-changing user.
// // It should be only for internal use.
// type UserCommandInternal struct {
// 	apiv1.UnimplementedUserCommandInternalServiceServer
// 	deleter service.DeleteUser
// }

// // NewUserCommandInternal creates an instance of UserCommandInternal.
// func NewUserCommandInternal(deleter service.DeleteUser) *UserCommandInternal {
// 	return &UserCommandInternal{deleter: deleter}
// }

// // DeleteUser handles HTTP/2 gRPC request similar to POST in HTTP/1.1.
// func (uc *UserCommandInternal) DeleteUser(ctx context.Context, request *apiv1.DeleteUserRequest) (*apiv1.DeleteUserResponse, error) {
// 	if request == nil || request.GetUser() == nil {
// 		return nil, entity.ErrEmptyUser()
// 	}

// 	id, err := uc.deleter.Delete(ctx, createUserFromRegisterUserRequest(request))
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &apiv1.RegisterUserResponse{Data: &apiv1.User{Id: id}}, nil
// }

// func createUserFromRegisterUserRequest(request *apiv1.RegisterUserRequest) *entity.User {
// 	return &entity.User{
// 		Name:     request.GetUser().GetName(),
// 		Email:    request.GetUser().GetEmail(),
// 		Password: request.GetUser().GetPassword(),
// 	}
// }
