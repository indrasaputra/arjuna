// transaction.proto defines service for transaction.

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             (unknown)
// source: api/v1/transaction.proto

package apiv1

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	TransactionCommandService_CreateTransaction_FullMethodName = "/api.v1.TransactionCommandService/CreateTransaction"
)

// TransactionCommandServiceClient is the client API for TransactionCommandService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// TransactionCommandService provides all use cases to work with transaction.
type TransactionCommandServiceClient interface {
	// CreateTransaction.
	//
	// This endpoint creates a transaction.
	CreateTransaction(ctx context.Context, in *CreateTransactionRequest, opts ...grpc.CallOption) (*CreateTransactionResponse, error)
}

type transactionCommandServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTransactionCommandServiceClient(cc grpc.ClientConnInterface) TransactionCommandServiceClient {
	return &transactionCommandServiceClient{cc}
}

func (c *transactionCommandServiceClient) CreateTransaction(ctx context.Context, in *CreateTransactionRequest, opts ...grpc.CallOption) (*CreateTransactionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateTransactionResponse)
	err := c.cc.Invoke(ctx, TransactionCommandService_CreateTransaction_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TransactionCommandServiceServer is the server API for TransactionCommandService service.
// All implementations must embed UnimplementedTransactionCommandServiceServer
// for forward compatibility
//
// TransactionCommandService provides all use cases to work with transaction.
type TransactionCommandServiceServer interface {
	// CreateTransaction.
	//
	// This endpoint creates a transaction.
	CreateTransaction(context.Context, *CreateTransactionRequest) (*CreateTransactionResponse, error)
	mustEmbedUnimplementedTransactionCommandServiceServer()
}

// UnimplementedTransactionCommandServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTransactionCommandServiceServer struct {
}

func (UnimplementedTransactionCommandServiceServer) CreateTransaction(context.Context, *CreateTransactionRequest) (*CreateTransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTransaction not implemented")
}
func (UnimplementedTransactionCommandServiceServer) mustEmbedUnimplementedTransactionCommandServiceServer() {
}

// UnsafeTransactionCommandServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TransactionCommandServiceServer will
// result in compilation errors.
type UnsafeTransactionCommandServiceServer interface {
	mustEmbedUnimplementedTransactionCommandServiceServer()
}

func RegisterTransactionCommandServiceServer(s grpc.ServiceRegistrar, srv TransactionCommandServiceServer) {
	s.RegisterService(&TransactionCommandService_ServiceDesc, srv)
}

func _TransactionCommandService_CreateTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateTransactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionCommandServiceServer).CreateTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TransactionCommandService_CreateTransaction_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionCommandServiceServer).CreateTransaction(ctx, req.(*CreateTransactionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TransactionCommandService_ServiceDesc is the grpc.ServiceDesc for TransactionCommandService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TransactionCommandService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.v1.TransactionCommandService",
	HandlerType: (*TransactionCommandServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateTransaction",
			Handler:    _TransactionCommandService_CreateTransaction_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/v1/transaction.proto",
}
