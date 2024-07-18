// transaction.proto defines service for transaction.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: api/v1/transaction.proto

package apiv1

import (
	reflect "reflect"
	sync "sync"

	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// TransactionErrorCode enumerates transaction error code.
type TransactionErrorCode int32

const (
	// Default enum code according to
	// https://medium.com/@akhaku/protobuf-definition-best-practices-87f281576f31.
	TransactionErrorCode_TRANSACTION_ERROR_CODE_UNSPECIFIED TransactionErrorCode = 0
	// Unexpected behavior occured in system.
	TransactionErrorCode_TRANSACTION_ERROR_CODE_INTERNAL TransactionErrorCode = 1
	// Transaction already exists.
	TransactionErrorCode_TRANSACTION_ERROR_CODE_ALREADY_EXISTS TransactionErrorCode = 2
	// Transaction instance is nil or empty.
	TransactionErrorCode_TRANSACTION_ERROR_CODE_EMPTY_TRANSACTION TransactionErrorCode = 3
	// Sender is invalid.
	TransactionErrorCode_TRANSACTION_ERROR_CODE_INVALID_SENDER TransactionErrorCode = 4
	// Receiver is invalid.
	TransactionErrorCode_TRANSACTION_ERROR_CODE_INVALID_RECEIVER TransactionErrorCode = 5
	// Amount must be numeric and greater than zero.
	TransactionErrorCode_TRANSACTION_ERROR_CODE_INVALID_AMOUNT TransactionErrorCode = 6
	// Idempotency key is missing.
	TransactionErrorCode_TRANSACTION_ERROR_CODE_MISSING_IDEMPOTENCY_KEY TransactionErrorCode = 7
)

// Enum value maps for TransactionErrorCode.
var (
	TransactionErrorCode_name = map[int32]string{
		0: "TRANSACTION_ERROR_CODE_UNSPECIFIED",
		1: "TRANSACTION_ERROR_CODE_INTERNAL",
		2: "TRANSACTION_ERROR_CODE_ALREADY_EXISTS",
		3: "TRANSACTION_ERROR_CODE_EMPTY_TRANSACTION",
		4: "TRANSACTION_ERROR_CODE_INVALID_SENDER",
		5: "TRANSACTION_ERROR_CODE_INVALID_RECEIVER",
		6: "TRANSACTION_ERROR_CODE_INVALID_AMOUNT",
		7: "TRANSACTION_ERROR_CODE_MISSING_IDEMPOTENCY_KEY",
	}
	TransactionErrorCode_value = map[string]int32{
		"TRANSACTION_ERROR_CODE_UNSPECIFIED":             0,
		"TRANSACTION_ERROR_CODE_INTERNAL":                1,
		"TRANSACTION_ERROR_CODE_ALREADY_EXISTS":          2,
		"TRANSACTION_ERROR_CODE_EMPTY_TRANSACTION":       3,
		"TRANSACTION_ERROR_CODE_INVALID_SENDER":          4,
		"TRANSACTION_ERROR_CODE_INVALID_RECEIVER":        5,
		"TRANSACTION_ERROR_CODE_INVALID_AMOUNT":          6,
		"TRANSACTION_ERROR_CODE_MISSING_IDEMPOTENCY_KEY": 7,
	}
)

func (x TransactionErrorCode) Enum() *TransactionErrorCode {
	p := new(TransactionErrorCode)
	*p = x
	return p
}

func (x TransactionErrorCode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TransactionErrorCode) Descriptor() protoreflect.EnumDescriptor {
	return file_api_v1_transaction_proto_enumTypes[0].Descriptor()
}

func (TransactionErrorCode) Type() protoreflect.EnumType {
	return &file_api_v1_transaction_proto_enumTypes[0]
}

func (x TransactionErrorCode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TransactionErrorCode.Descriptor instead.
func (TransactionErrorCode) EnumDescriptor() ([]byte, []int) {
	return file_api_v1_transaction_proto_rawDescGZIP(), []int{0}
}

// CreateTransactionRequest represents request for create transaction.
type CreateTransactionRequest struct {
	state         protoimpl.MessageState
	Transaction   *Transaction `protobuf:"bytes,1,opt,name=transaction,proto3" json:"transaction,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateTransactionRequest) Reset() {
	*x = CreateTransactionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_transaction_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateTransactionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateTransactionRequest) ProtoMessage() {}

func (x *CreateTransactionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_transaction_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateTransactionRequest.ProtoReflect.Descriptor instead.
func (*CreateTransactionRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_transaction_proto_rawDescGZIP(), []int{0}
}

func (x *CreateTransactionRequest) GetTransaction() *Transaction {
	if x != nil {
		return x.Transaction
	}
	return nil
}

// CreateTransactionResponse represents response from create transaction.
type CreateTransactionResponse struct {
	state         protoimpl.MessageState
	Data          *Transaction `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateTransactionResponse) Reset() {
	*x = CreateTransactionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_transaction_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateTransactionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateTransactionResponse) ProtoMessage() {}

func (x *CreateTransactionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_transaction_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateTransactionResponse.ProtoReflect.Descriptor instead.
func (*CreateTransactionResponse) Descriptor() ([]byte, []int) {
	return file_api_v1_transaction_proto_rawDescGZIP(), []int{1}
}

func (x *CreateTransactionResponse) GetData() *Transaction {
	if x != nil {
		return x.Data
	}
	return nil
}

// Transaction represents transaction.
type Transaction struct {
	state         protoimpl.MessageState
	CreatedAt     *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=created_at,proto3" json:"created_at,omitempty"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	SenderId      string                 `protobuf:"bytes,2,opt,name=sender_id,proto3" json:"sender_id,omitempty"`
	ReceiverId    string                 `protobuf:"bytes,3,opt,name=receiver_id,proto3" json:"receiver_id,omitempty"`
	Amount        string                 `protobuf:"bytes,4,opt,name=amount,proto3" json:"amount,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Transaction) Reset() {
	*x = Transaction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_transaction_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Transaction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Transaction) ProtoMessage() {}

func (x *Transaction) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_transaction_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Transaction.ProtoReflect.Descriptor instead.
func (*Transaction) Descriptor() ([]byte, []int) {
	return file_api_v1_transaction_proto_rawDescGZIP(), []int{2}
}

func (x *Transaction) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Transaction) GetSenderId() string {
	if x != nil {
		return x.SenderId
	}
	return ""
}

func (x *Transaction) GetReceiverId() string {
	if x != nil {
		return x.ReceiverId
	}
	return ""
}

func (x *Transaction) GetAmount() string {
	if x != nil {
		return x.Amount
	}
	return ""
}

func (x *Transaction) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

// TransactionError represents message for any error happening in transaction service.
type TransactionError struct {
	state         protoimpl.MessageState
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
	ErrorCode     TransactionErrorCode `protobuf:"varint,1,opt,name=error_code,json=errorCode,proto3,enum=api.v1.TransactionErrorCode" json:"error_code,omitempty"`
}

func (x *TransactionError) Reset() {
	*x = TransactionError{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_transaction_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TransactionError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransactionError) ProtoMessage() {}

func (x *TransactionError) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_transaction_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransactionError.ProtoReflect.Descriptor instead.
func (*TransactionError) Descriptor() ([]byte, []int) {
	return file_api_v1_transaction_proto_rawDescGZIP(), []int{3}
}

func (x *TransactionError) GetErrorCode() TransactionErrorCode {
	if x != nil {
		return x.ErrorCode
	}
	return TransactionErrorCode_TRANSACTION_ERROR_CODE_UNSPECIFIED
}

var File_api_v1_transaction_proto protoreflect.FileDescriptor

var file_api_v1_transaction_proto_rawDesc = []byte{
	0x0a, 0x18, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x61, 0x70, 0x69, 0x2e,
	0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61,
	0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x66, 0x69, 0x65,
	0x6c, 0x64, 0x5f, 0x62, 0x65, 0x68, 0x61, 0x76, 0x69, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f,
	0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x51, 0x0a, 0x18, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x35,
	0x0a, 0x0b, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x72, 0x61,
	0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0b, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x44, 0x0a, 0x19, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54,
	0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x27, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x13, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0xb5, 0x02, 0x0a, 0x0b,
	0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1d, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0d, 0x92, 0x41, 0x07, 0x4a, 0x05, 0x22, 0x31,
	0x32, 0x33, 0x22, 0xe0, 0x41, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x41, 0x0a, 0x09, 0x73, 0x65,
	0x6e, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x23, 0x92,
	0x41, 0x20, 0x32, 0x19, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x27,
	0x73, 0x20, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x27, 0x73, 0x20, 0x69, 0x64, 0x4a, 0x03, 0x22,
	0x31, 0x22, 0x52, 0x09, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x12, 0x47, 0x0a,
	0x0b, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x25, 0x92, 0x41, 0x22, 0x32, 0x1b, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x27, 0x73, 0x20, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x27,
	0x73, 0x20, 0x69, 0x64, 0x4a, 0x03, 0x22, 0x32, 0x22, 0x52, 0x0b, 0x72, 0x65, 0x63, 0x65, 0x69,
	0x76, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x12, 0x3a, 0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x22, 0x92, 0x41, 0x1f, 0x32, 0x14, 0x54, 0x72, 0x61,
	0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x27, 0x73, 0x20, 0x61, 0x6d, 0x6f, 0x75, 0x6e,
	0x74, 0x4a, 0x07, 0x22, 0x31, 0x30, 0x2e, 0x32, 0x33, 0x22, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75,
	0x6e, 0x74, 0x12, 0x3f, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x42, 0x03, 0xe0, 0x41, 0x03, 0x52, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64,
	0x5f, 0x61, 0x74, 0x22, 0x4f, 0x0a, 0x10, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x3b, 0x0a, 0x0a, 0x65, 0x72, 0x72, 0x6f, 0x72,
	0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1c, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x45, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x09, 0x65, 0x72, 0x72, 0x6f, 0x72,
	0x43, 0x6f, 0x64, 0x65, 0x2a, 0xf3, 0x02, 0x0a, 0x14, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x26, 0x0a,
	0x22, 0x54, 0x52, 0x41, 0x4e, 0x53, 0x41, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x45, 0x52, 0x52,
	0x4f, 0x52, 0x5f, 0x43, 0x4f, 0x44, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46,
	0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x23, 0x0a, 0x1f, 0x54, 0x52, 0x41, 0x4e, 0x53, 0x41, 0x43,
	0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x5f, 0x43, 0x4f, 0x44, 0x45, 0x5f,
	0x49, 0x4e, 0x54, 0x45, 0x52, 0x4e, 0x41, 0x4c, 0x10, 0x01, 0x12, 0x29, 0x0a, 0x25, 0x54, 0x52,
	0x41, 0x4e, 0x53, 0x41, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x5f,
	0x43, 0x4f, 0x44, 0x45, 0x5f, 0x41, 0x4c, 0x52, 0x45, 0x41, 0x44, 0x59, 0x5f, 0x45, 0x58, 0x49,
	0x53, 0x54, 0x53, 0x10, 0x02, 0x12, 0x2c, 0x0a, 0x28, 0x54, 0x52, 0x41, 0x4e, 0x53, 0x41, 0x43,
	0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x5f, 0x43, 0x4f, 0x44, 0x45, 0x5f,
	0x45, 0x4d, 0x50, 0x54, 0x59, 0x5f, 0x54, 0x52, 0x41, 0x4e, 0x53, 0x41, 0x43, 0x54, 0x49, 0x4f,
	0x4e, 0x10, 0x03, 0x12, 0x29, 0x0a, 0x25, 0x54, 0x52, 0x41, 0x4e, 0x53, 0x41, 0x43, 0x54, 0x49,
	0x4f, 0x4e, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x5f, 0x43, 0x4f, 0x44, 0x45, 0x5f, 0x49, 0x4e,
	0x56, 0x41, 0x4c, 0x49, 0x44, 0x5f, 0x53, 0x45, 0x4e, 0x44, 0x45, 0x52, 0x10, 0x04, 0x12, 0x2b,
	0x0a, 0x27, 0x54, 0x52, 0x41, 0x4e, 0x53, 0x41, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x45, 0x52,
	0x52, 0x4f, 0x52, 0x5f, 0x43, 0x4f, 0x44, 0x45, 0x5f, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44,
	0x5f, 0x52, 0x45, 0x43, 0x45, 0x49, 0x56, 0x45, 0x52, 0x10, 0x05, 0x12, 0x29, 0x0a, 0x25, 0x54,
	0x52, 0x41, 0x4e, 0x53, 0x41, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52,
	0x5f, 0x43, 0x4f, 0x44, 0x45, 0x5f, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x5f, 0x41, 0x4d,
	0x4f, 0x55, 0x4e, 0x54, 0x10, 0x06, 0x12, 0x32, 0x0a, 0x2e, 0x54, 0x52, 0x41, 0x4e, 0x53, 0x41,
	0x43, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x5f, 0x43, 0x4f, 0x44, 0x45,
	0x5f, 0x4d, 0x49, 0x53, 0x53, 0x49, 0x4e, 0x47, 0x5f, 0x49, 0x44, 0x45, 0x4d, 0x50, 0x4f, 0x54,
	0x45, 0x4e, 0x43, 0x59, 0x5f, 0x4b, 0x45, 0x59, 0x10, 0x07, 0x32, 0xb4, 0x02, 0x0a, 0x19, 0x54,
	0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e,
	0x64, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0xd2, 0x01, 0x0a, 0x11, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x20,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x72,
	0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x21, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x78, 0x92, 0x41, 0x50, 0x0a, 0x0b, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2a, 0x11, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x72, 0x61,
	0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x72, 0x2e, 0x0a, 0x13, 0x0a, 0x0d, 0x41, 0x75,
	0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x28, 0x01, 0x0a,
	0x17, 0x0a, 0x11, 0x58, 0x2d, 0x49, 0x64, 0x65, 0x6d, 0x70, 0x6f, 0x74, 0x65, 0x6e, 0x63, 0x79,
	0x2d, 0x4b, 0x65, 0x79, 0x18, 0x01, 0x28, 0x01, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1f, 0x3a, 0x0b,
	0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x10, 0x2f, 0x76, 0x31,
	0x2f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x1a, 0x42, 0x92,
	0x41, 0x3f, 0x12, 0x3d, 0x54, 0x68, 0x69, 0x73, 0x20, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x20, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x73, 0x20, 0x61, 0x6c, 0x6c, 0x20, 0x75, 0x73,
	0x65, 0x20, 0x63, 0x61, 0x73, 0x65, 0x73, 0x20, 0x74, 0x6f, 0x20, 0x77, 0x6f, 0x72, 0x6b, 0x20,
	0x77, 0x69, 0x74, 0x68, 0x20, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x2e, 0x42, 0x9b, 0x02, 0x92, 0x41, 0xd6, 0x01, 0x12, 0x9c, 0x01, 0x0a, 0x0f, 0x54, 0x72, 0x61,
	0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x20, 0x41, 0x50, 0x49, 0x22, 0x30, 0x0a, 0x0d,
	0x49, 0x6e, 0x64, 0x72, 0x61, 0x20, 0x53, 0x61, 0x70, 0x75, 0x74, 0x72, 0x61, 0x12, 0x1f, 0x68,
	0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x69, 0x6e, 0x64, 0x72, 0x61, 0x73, 0x61, 0x70, 0x75, 0x74, 0x72, 0x61, 0x2a, 0x50,
	0x0a, 0x14, 0x42, 0x53, 0x44, 0x20, 0x33, 0x2d, 0x43, 0x6c, 0x61, 0x75, 0x73, 0x65, 0x20, 0x4c,
	0x69, 0x63, 0x65, 0x6e, 0x73, 0x65, 0x12, 0x38, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x6e, 0x64, 0x72, 0x61,
	0x73, 0x61, 0x70, 0x75, 0x74, 0x72, 0x61, 0x2f, 0x61, 0x72, 0x6a, 0x75, 0x6e, 0x61, 0x2f, 0x62,
	0x6c, 0x6f, 0x62, 0x2f, 0x6d, 0x61, 0x69, 0x6e, 0x2f, 0x4c, 0x49, 0x43, 0x45, 0x4e, 0x53, 0x45,
	0x32, 0x05, 0x31, 0x2e, 0x30, 0x2e, 0x30, 0x1a, 0x0e, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x68, 0x6f,
	0x73, 0x74, 0x3a, 0x38, 0x30, 0x30, 0x30, 0x2a, 0x01, 0x01, 0x32, 0x10, 0x61, 0x70, 0x70, 0x6c,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x6a, 0x73, 0x6f, 0x6e, 0x3a, 0x10, 0x61, 0x70,
	0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x6a, 0x73, 0x6f, 0x6e, 0x5a, 0x3f,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x6e, 0x64, 0x72, 0x61,
	0x73, 0x61, 0x70, 0x75, 0x74, 0x72, 0x61, 0x2f, 0x61, 0x72, 0x6a, 0x75, 0x6e, 0x61, 0x2f, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x3b, 0x61, 0x70, 0x69, 0x76, 0x31, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_v1_transaction_proto_rawDescOnce sync.Once
	file_api_v1_transaction_proto_rawDescData = file_api_v1_transaction_proto_rawDesc
)

func file_api_v1_transaction_proto_rawDescGZIP() []byte {
	file_api_v1_transaction_proto_rawDescOnce.Do(func() {
		file_api_v1_transaction_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_v1_transaction_proto_rawDescData)
	})
	return file_api_v1_transaction_proto_rawDescData
}

var file_api_v1_transaction_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_api_v1_transaction_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_api_v1_transaction_proto_goTypes = []any{
	(TransactionErrorCode)(0),         // 0: api.v1.TransactionErrorCode
	(*CreateTransactionRequest)(nil),  // 1: api.v1.CreateTransactionRequest
	(*CreateTransactionResponse)(nil), // 2: api.v1.CreateTransactionResponse
	(*Transaction)(nil),               // 3: api.v1.Transaction
	(*TransactionError)(nil),          // 4: api.v1.TransactionError
	(*timestamppb.Timestamp)(nil),     // 5: google.protobuf.Timestamp
}
var file_api_v1_transaction_proto_depIdxs = []int32{
	3, // 0: api.v1.CreateTransactionRequest.transaction:type_name -> api.v1.Transaction
	3, // 1: api.v1.CreateTransactionResponse.data:type_name -> api.v1.Transaction
	5, // 2: api.v1.Transaction.created_at:type_name -> google.protobuf.Timestamp
	0, // 3: api.v1.TransactionError.error_code:type_name -> api.v1.TransactionErrorCode
	1, // 4: api.v1.TransactionCommandService.CreateTransaction:input_type -> api.v1.CreateTransactionRequest
	2, // 5: api.v1.TransactionCommandService.CreateTransaction:output_type -> api.v1.CreateTransactionResponse
	5, // [5:6] is the sub-list for method output_type
	4, // [4:5] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_api_v1_transaction_proto_init() }
func file_api_v1_transaction_proto_init() {
	if File_api_v1_transaction_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_v1_transaction_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*CreateTransactionRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_v1_transaction_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*CreateTransactionResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_v1_transaction_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*Transaction); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_v1_transaction_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*TransactionError); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_v1_transaction_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_v1_transaction_proto_goTypes,
		DependencyIndexes: file_api_v1_transaction_proto_depIdxs,
		EnumInfos:         file_api_v1_transaction_proto_enumTypes,
		MessageInfos:      file_api_v1_transaction_proto_msgTypes,
	}.Build()
	File_api_v1_transaction_proto = out.File
	file_api_v1_transaction_proto_rawDesc = nil
	file_api_v1_transaction_proto_goTypes = nil
	file_api_v1_transaction_proto_depIdxs = nil
}
