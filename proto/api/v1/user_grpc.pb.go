// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package apiv1

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// UserCommandServiceClient is the client API for UserCommandService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserCommandServiceClient interface {
	// Register a new user.
	//
	// This endpoint registers a new user.
	RegisterUser(ctx context.Context, in *RegisterUserRequest, opts ...grpc.CallOption) (*RegisterUserResponse, error)
}

type userCommandServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserCommandServiceClient(cc grpc.ClientConnInterface) UserCommandServiceClient {
	return &userCommandServiceClient{cc}
}

func (c *userCommandServiceClient) RegisterUser(ctx context.Context, in *RegisterUserRequest, opts ...grpc.CallOption) (*RegisterUserResponse, error) {
	out := new(RegisterUserResponse)
	err := c.cc.Invoke(ctx, "/api.v1.UserCommandService/RegisterUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserCommandServiceServer is the server API for UserCommandService service.
// All implementations must embed UnimplementedUserCommandServiceServer
// for forward compatibility
type UserCommandServiceServer interface {
	// Register a new user.
	//
	// This endpoint registers a new user.
	RegisterUser(context.Context, *RegisterUserRequest) (*RegisterUserResponse, error)
	mustEmbedUnimplementedUserCommandServiceServer()
}

// UnimplementedUserCommandServiceServer must be embedded to have forward compatible implementations.
type UnimplementedUserCommandServiceServer struct {
}

func (UnimplementedUserCommandServiceServer) RegisterUser(context.Context, *RegisterUserRequest) (*RegisterUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterUser not implemented")
}
func (UnimplementedUserCommandServiceServer) mustEmbedUnimplementedUserCommandServiceServer() {}

// UnsafeUserCommandServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserCommandServiceServer will
// result in compilation errors.
type UnsafeUserCommandServiceServer interface {
	mustEmbedUnimplementedUserCommandServiceServer()
}

func RegisterUserCommandServiceServer(s grpc.ServiceRegistrar, srv UserCommandServiceServer) {
	s.RegisterService(&UserCommandService_ServiceDesc, srv)
}

func _UserCommandService_RegisterUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserCommandServiceServer).RegisterUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.v1.UserCommandService/RegisterUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserCommandServiceServer).RegisterUser(ctx, req.(*RegisterUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserCommandService_ServiceDesc is the grpc.ServiceDesc for UserCommandService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserCommandService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.v1.UserCommandService",
	HandlerType: (*UserCommandServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterUser",
			Handler:    _UserCommandService_RegisterUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/v1/user.proto",
}

// UserCommandInternalServiceClient is the client API for UserCommandInternalService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserCommandInternalServiceClient interface {
	// Delete a user.
	//
	// This endpoint deletes a new user.
	// It is expected to be hidden or internal use only.
	DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*DeleteUserResponse, error)
}

type userCommandInternalServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserCommandInternalServiceClient(cc grpc.ClientConnInterface) UserCommandInternalServiceClient {
	return &userCommandInternalServiceClient{cc}
}

func (c *userCommandInternalServiceClient) DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*DeleteUserResponse, error) {
	out := new(DeleteUserResponse)
	err := c.cc.Invoke(ctx, "/api.v1.UserCommandInternalService/DeleteUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserCommandInternalServiceServer is the server API for UserCommandInternalService service.
// All implementations must embed UnimplementedUserCommandInternalServiceServer
// for forward compatibility
type UserCommandInternalServiceServer interface {
	// Delete a user.
	//
	// This endpoint deletes a new user.
	// It is expected to be hidden or internal use only.
	DeleteUser(context.Context, *DeleteUserRequest) (*DeleteUserResponse, error)
	mustEmbedUnimplementedUserCommandInternalServiceServer()
}

// UnimplementedUserCommandInternalServiceServer must be embedded to have forward compatible implementations.
type UnimplementedUserCommandInternalServiceServer struct {
}

func (UnimplementedUserCommandInternalServiceServer) DeleteUser(context.Context, *DeleteUserRequest) (*DeleteUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUser not implemented")
}
func (UnimplementedUserCommandInternalServiceServer) mustEmbedUnimplementedUserCommandInternalServiceServer() {
}

// UnsafeUserCommandInternalServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserCommandInternalServiceServer will
// result in compilation errors.
type UnsafeUserCommandInternalServiceServer interface {
	mustEmbedUnimplementedUserCommandInternalServiceServer()
}

func RegisterUserCommandInternalServiceServer(s grpc.ServiceRegistrar, srv UserCommandInternalServiceServer) {
	s.RegisterService(&UserCommandInternalService_ServiceDesc, srv)
}

func _UserCommandInternalService_DeleteUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserCommandInternalServiceServer).DeleteUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.v1.UserCommandInternalService/DeleteUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserCommandInternalServiceServer).DeleteUser(ctx, req.(*DeleteUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserCommandInternalService_ServiceDesc is the grpc.ServiceDesc for UserCommandInternalService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserCommandInternalService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.v1.UserCommandInternalService",
	HandlerType: (*UserCommandInternalServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DeleteUser",
			Handler:    _UserCommandInternalService_DeleteUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/v1/user.proto",
}

// UserQueryServiceClient is the client API for UserQueryService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserQueryServiceClient interface {
	// Get all users.
	//
	// This endpoint gets all available users in the system.
	// Currently, it only retrieves 10 users at most.
	GetAllUsers(ctx context.Context, in *GetAllUsersRequest, opts ...grpc.CallOption) (*GetAllUsersResponse, error)
}

type userQueryServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserQueryServiceClient(cc grpc.ClientConnInterface) UserQueryServiceClient {
	return &userQueryServiceClient{cc}
}

func (c *userQueryServiceClient) GetAllUsers(ctx context.Context, in *GetAllUsersRequest, opts ...grpc.CallOption) (*GetAllUsersResponse, error) {
	out := new(GetAllUsersResponse)
	err := c.cc.Invoke(ctx, "/api.v1.UserQueryService/GetAllUsers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserQueryServiceServer is the server API for UserQueryService service.
// All implementations must embed UnimplementedUserQueryServiceServer
// for forward compatibility
type UserQueryServiceServer interface {
	// Get all users.
	//
	// This endpoint gets all available users in the system.
	// Currently, it only retrieves 10 users at most.
	GetAllUsers(context.Context, *GetAllUsersRequest) (*GetAllUsersResponse, error)
	mustEmbedUnimplementedUserQueryServiceServer()
}

// UnimplementedUserQueryServiceServer must be embedded to have forward compatible implementations.
type UnimplementedUserQueryServiceServer struct {
}

func (UnimplementedUserQueryServiceServer) GetAllUsers(context.Context, *GetAllUsersRequest) (*GetAllUsersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllUsers not implemented")
}
func (UnimplementedUserQueryServiceServer) mustEmbedUnimplementedUserQueryServiceServer() {}

// UnsafeUserQueryServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserQueryServiceServer will
// result in compilation errors.
type UnsafeUserQueryServiceServer interface {
	mustEmbedUnimplementedUserQueryServiceServer()
}

func RegisterUserQueryServiceServer(s grpc.ServiceRegistrar, srv UserQueryServiceServer) {
	s.RegisterService(&UserQueryService_ServiceDesc, srv)
}

func _UserQueryService_GetAllUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllUsersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserQueryServiceServer).GetAllUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.v1.UserQueryService/GetAllUsers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserQueryServiceServer).GetAllUsers(ctx, req.(*GetAllUsersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserQueryService_ServiceDesc is the grpc.ServiceDesc for UserQueryService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserQueryService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.v1.UserQueryService",
	HandlerType: (*UserQueryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAllUsers",
			Handler:    _UserQueryService_GetAllUsers_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/v1/user.proto",
}
