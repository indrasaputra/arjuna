// wallet.proto defines service for wallet.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: api/v1/wallet.proto

package apiv1

import (
	reflect "reflect"
	sync "sync"

	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// WalletErrorCode enumerates wallet error code.
type WalletErrorCode int32

const (
	// Default enum code according to
	// https://medium.com/@akhaku/protobuf-definition-best-practices-87f281576f31.
	WalletErrorCode_WALLET_ERROR_CODE_UNSPECIFIED WalletErrorCode = 0
	// Unexpected behavior occured in system.
	WalletErrorCode_WALLET_ERROR_CODE_INTERNAL WalletErrorCode = 1
	// Wallet already exists.
	WalletErrorCode_WALLET_ERROR_CODE_ALREADY_EXISTS WalletErrorCode = 2
	// Wallet instance is nil or empty.
	WalletErrorCode_WALLET_ERROR_CODE_EMPTY_WALLET WalletErrorCode = 3
	// Balance must be numeric and greater than or equal to zero.
	WalletErrorCode_WALLET_ERROR_CODE_INVALID_BALANCE WalletErrorCode = 6
	// Idempotency key is missing.
	WalletErrorCode_WALLET_ERROR_CODE_MISSING_IDEMPOTENCY_KEY WalletErrorCode = 7
	// User is invalid.
	WalletErrorCode_WALLET_ERROR_CODE_INVALID_USER WalletErrorCode = 8
	// Balance must be numeric and greater than zero.
	WalletErrorCode_WALLET_ERROR_CODE_INVALID_AMOUNT WalletErrorCode = 9
)

// Enum value maps for WalletErrorCode.
var (
	WalletErrorCode_name = map[int32]string{
		0: "WALLET_ERROR_CODE_UNSPECIFIED",
		1: "WALLET_ERROR_CODE_INTERNAL",
		2: "WALLET_ERROR_CODE_ALREADY_EXISTS",
		3: "WALLET_ERROR_CODE_EMPTY_WALLET",
		6: "WALLET_ERROR_CODE_INVALID_BALANCE",
		7: "WALLET_ERROR_CODE_MISSING_IDEMPOTENCY_KEY",
		8: "WALLET_ERROR_CODE_INVALID_USER",
		9: "WALLET_ERROR_CODE_INVALID_AMOUNT",
	}
	WalletErrorCode_value = map[string]int32{
		"WALLET_ERROR_CODE_UNSPECIFIED":             0,
		"WALLET_ERROR_CODE_INTERNAL":                1,
		"WALLET_ERROR_CODE_ALREADY_EXISTS":          2,
		"WALLET_ERROR_CODE_EMPTY_WALLET":            3,
		"WALLET_ERROR_CODE_INVALID_BALANCE":         6,
		"WALLET_ERROR_CODE_MISSING_IDEMPOTENCY_KEY": 7,
		"WALLET_ERROR_CODE_INVALID_USER":            8,
		"WALLET_ERROR_CODE_INVALID_AMOUNT":          9,
	}
)

func (x WalletErrorCode) Enum() *WalletErrorCode {
	p := new(WalletErrorCode)
	*p = x
	return p
}

func (x WalletErrorCode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (WalletErrorCode) Descriptor() protoreflect.EnumDescriptor {
	return file_api_v1_wallet_proto_enumTypes[0].Descriptor()
}

func (WalletErrorCode) Type() protoreflect.EnumType {
	return &file_api_v1_wallet_proto_enumTypes[0]
}

func (x WalletErrorCode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use WalletErrorCode.Descriptor instead.
func (WalletErrorCode) EnumDescriptor() ([]byte, []int) {
	return file_api_v1_wallet_proto_rawDescGZIP(), []int{0}
}

// CreateWalletRequest represents request for create wallet.
type CreateWalletRequest struct {
	state         protoimpl.MessageState
	Wallet        *Wallet `protobuf:"bytes,1,opt,name=wallet,proto3" json:"wallet,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateWalletRequest) Reset() {
	*x = CreateWalletRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_wallet_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateWalletRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateWalletRequest) ProtoMessage() {}

func (x *CreateWalletRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_wallet_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateWalletRequest.ProtoReflect.Descriptor instead.
func (*CreateWalletRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_wallet_proto_rawDescGZIP(), []int{0}
}

func (x *CreateWalletRequest) GetWallet() *Wallet {
	if x != nil {
		return x.Wallet
	}
	return nil
}

// CreateWalletResponse represents response from create wallet.
type CreateWalletResponse struct {
	state         protoimpl.MessageState
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateWalletResponse) Reset() {
	*x = CreateWalletResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_wallet_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateWalletResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateWalletResponse) ProtoMessage() {}

func (x *CreateWalletResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_wallet_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateWalletResponse.ProtoReflect.Descriptor instead.
func (*CreateWalletResponse) Descriptor() ([]byte, []int) {
	return file_api_v1_wallet_proto_rawDescGZIP(), []int{1}
}

// TopupWalletRequest represents request for topup wallet.
type TopupWalletRequest struct {
	state         protoimpl.MessageState
	Topup         *Topup `protobuf:"bytes,1,opt,name=topup,proto3" json:"topup,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TopupWalletRequest) Reset() {
	*x = TopupWalletRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_wallet_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TopupWalletRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TopupWalletRequest) ProtoMessage() {}

func (x *TopupWalletRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_wallet_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TopupWalletRequest.ProtoReflect.Descriptor instead.
func (*TopupWalletRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_wallet_proto_rawDescGZIP(), []int{2}
}

func (x *TopupWalletRequest) GetTopup() *Topup {
	if x != nil {
		return x.Topup
	}
	return nil
}

// TopupWalletResponse represents response from topup wallet.
type TopupWalletResponse struct {
	state         protoimpl.MessageState
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TopupWalletResponse) Reset() {
	*x = TopupWalletResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_wallet_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TopupWalletResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TopupWalletResponse) ProtoMessage() {}

func (x *TopupWalletResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_wallet_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TopupWalletResponse.ProtoReflect.Descriptor instead.
func (*TopupWalletResponse) Descriptor() ([]byte, []int) {
	return file_api_v1_wallet_proto_rawDescGZIP(), []int{3}
}

// Wallet represents wallet.
type Wallet struct {
	state         protoimpl.MessageState
	Id            string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId        string `protobuf:"bytes,2,opt,name=user_id,proto3" json:"user_id,omitempty"`
	Balance       string `protobuf:"bytes,3,opt,name=balance,proto3" json:"balance,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Wallet) Reset() {
	*x = Wallet{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_wallet_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Wallet) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Wallet) ProtoMessage() {}

func (x *Wallet) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_wallet_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Wallet.ProtoReflect.Descriptor instead.
func (*Wallet) Descriptor() ([]byte, []int) {
	return file_api_v1_wallet_proto_rawDescGZIP(), []int{4}
}

func (x *Wallet) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Wallet) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *Wallet) GetBalance() string {
	if x != nil {
		return x.Balance
	}
	return ""
}

// Topup represents topup.
type Topup struct {
	state         protoimpl.MessageState
	WalletId      string `protobuf:"bytes,1,opt,name=wallet_id,proto3" json:"wallet_id,omitempty"`
	Amount        string `protobuf:"bytes,2,opt,name=amount,proto3" json:"amount,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Topup) Reset() {
	*x = Topup{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_wallet_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Topup) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Topup) ProtoMessage() {}

func (x *Topup) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_wallet_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Topup.ProtoReflect.Descriptor instead.
func (*Topup) Descriptor() ([]byte, []int) {
	return file_api_v1_wallet_proto_rawDescGZIP(), []int{5}
}

func (x *Topup) GetWalletId() string {
	if x != nil {
		return x.WalletId
	}
	return ""
}

func (x *Topup) GetAmount() string {
	if x != nil {
		return x.Amount
	}
	return ""
}

// WalletError represents message for any error happening in wallet service.
type WalletError struct {
	state         protoimpl.MessageState
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
	ErrorCode     WalletErrorCode `protobuf:"varint,1,opt,name=error_code,json=errorCode,proto3,enum=api.v1.WalletErrorCode" json:"error_code,omitempty"`
}

func (x *WalletError) Reset() {
	*x = WalletError{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_wallet_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WalletError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WalletError) ProtoMessage() {}

func (x *WalletError) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_wallet_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WalletError.ProtoReflect.Descriptor instead.
func (*WalletError) Descriptor() ([]byte, []int) {
	return file_api_v1_wallet_proto_rawDescGZIP(), []int{6}
}

func (x *WalletError) GetErrorCode() WalletErrorCode {
	if x != nil {
		return x.ErrorCode
	}
	return WalletErrorCode_WALLET_ERROR_CODE_UNSPECIFIED
}

var File_api_v1_wallet_proto protoreflect.FileDescriptor

var file_api_v1_wallet_proto_rawDesc = []byte{
	0x0a, 0x13, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x62, 0x65,
	0x68, 0x61, 0x76, 0x69, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69,
	0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3d, 0x0a, 0x13,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x26, 0x0a, 0x06, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x57, 0x61, 0x6c,
	0x6c, 0x65, 0x74, 0x52, 0x06, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x22, 0x16, 0x0a, 0x14, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x39, 0x0a, 0x12, 0x54, 0x6f, 0x70, 0x75, 0x70, 0x57, 0x61, 0x6c, 0x6c,
	0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x23, 0x0a, 0x05, 0x74, 0x6f, 0x70,
	0x75, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76,
	0x31, 0x2e, 0x54, 0x6f, 0x70, 0x75, 0x70, 0x52, 0x05, 0x74, 0x6f, 0x70, 0x75, 0x70, 0x22, 0x15,
	0x0a, 0x13, 0x54, 0x6f, 0x70, 0x75, 0x70, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x99, 0x01, 0x0a, 0x06, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74,
	0x12, 0x1d, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0d, 0x92, 0x41,
	0x07, 0x4a, 0x05, 0x22, 0x31, 0x32, 0x33, 0x22, 0xe0, 0x41, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x36, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x1c, 0x92, 0x41, 0x19, 0x32, 0x12, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x27, 0x73, 0x20,
	0x75, 0x73, 0x65, 0x72, 0x27, 0x73, 0x20, 0x69, 0x64, 0x4a, 0x03, 0x22, 0x31, 0x22, 0x52, 0x07,
	0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x12, 0x38, 0x0a, 0x07, 0x62, 0x61, 0x6c, 0x61, 0x6e,
	0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x1e, 0x92, 0x41, 0x1b, 0x32, 0x10, 0x57,
	0x61, 0x6c, 0x6c, 0x65, 0x74, 0x27, 0x73, 0x20, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x4a,
	0x07, 0x22, 0x31, 0x30, 0x2e, 0x32, 0x33, 0x22, 0x52, 0x07, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63,
	0x65, 0x22, 0x70, 0x0a, 0x05, 0x54, 0x6f, 0x70, 0x75, 0x70, 0x12, 0x33, 0x0a, 0x09, 0x77, 0x61,
	0x6c, 0x6c, 0x65, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x15, 0x92,
	0x41, 0x12, 0x32, 0x0b, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x27, 0x73, 0x20, 0x69, 0x64, 0x4a,
	0x03, 0x22, 0x31, 0x22, 0x52, 0x09, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x5f, 0x69, 0x64, 0x12,
	0x32, 0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x1a, 0x92, 0x41, 0x17, 0x32, 0x0c, 0x54, 0x6f, 0x70, 0x75, 0x70, 0x20, 0x61, 0x6d, 0x6f, 0x75,
	0x6e, 0x74, 0x4a, 0x07, 0x22, 0x31, 0x30, 0x2e, 0x32, 0x33, 0x22, 0x52, 0x06, 0x61, 0x6d, 0x6f,
	0x75, 0x6e, 0x74, 0x22, 0x45, 0x0a, 0x0b, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x45, 0x72, 0x72,
	0x6f, 0x72, 0x12, 0x36, 0x0a, 0x0a, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f, 0x63, 0x6f, 0x64, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x17, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e,
	0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x52,
	0x09, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x2a, 0xbe, 0x02, 0x0a, 0x0f, 0x57,
	0x61, 0x6c, 0x6c, 0x65, 0x74, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x21,
	0x0a, 0x1d, 0x57, 0x41, 0x4c, 0x4c, 0x45, 0x54, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x5f, 0x43,
	0x4f, 0x44, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10,
	0x00, 0x12, 0x1e, 0x0a, 0x1a, 0x57, 0x41, 0x4c, 0x4c, 0x45, 0x54, 0x5f, 0x45, 0x52, 0x52, 0x4f,
	0x52, 0x5f, 0x43, 0x4f, 0x44, 0x45, 0x5f, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x4e, 0x41, 0x4c, 0x10,
	0x01, 0x12, 0x24, 0x0a, 0x20, 0x57, 0x41, 0x4c, 0x4c, 0x45, 0x54, 0x5f, 0x45, 0x52, 0x52, 0x4f,
	0x52, 0x5f, 0x43, 0x4f, 0x44, 0x45, 0x5f, 0x41, 0x4c, 0x52, 0x45, 0x41, 0x44, 0x59, 0x5f, 0x45,
	0x58, 0x49, 0x53, 0x54, 0x53, 0x10, 0x02, 0x12, 0x22, 0x0a, 0x1e, 0x57, 0x41, 0x4c, 0x4c, 0x45,
	0x54, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x5f, 0x43, 0x4f, 0x44, 0x45, 0x5f, 0x45, 0x4d, 0x50,
	0x54, 0x59, 0x5f, 0x57, 0x41, 0x4c, 0x4c, 0x45, 0x54, 0x10, 0x03, 0x12, 0x25, 0x0a, 0x21, 0x57,
	0x41, 0x4c, 0x4c, 0x45, 0x54, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x5f, 0x43, 0x4f, 0x44, 0x45,
	0x5f, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x5f, 0x42, 0x41, 0x4c, 0x41, 0x4e, 0x43, 0x45,
	0x10, 0x06, 0x12, 0x2d, 0x0a, 0x29, 0x57, 0x41, 0x4c, 0x4c, 0x45, 0x54, 0x5f, 0x45, 0x52, 0x52,
	0x4f, 0x52, 0x5f, 0x43, 0x4f, 0x44, 0x45, 0x5f, 0x4d, 0x49, 0x53, 0x53, 0x49, 0x4e, 0x47, 0x5f,
	0x49, 0x44, 0x45, 0x4d, 0x50, 0x4f, 0x54, 0x45, 0x4e, 0x43, 0x59, 0x5f, 0x4b, 0x45, 0x59, 0x10,
	0x07, 0x12, 0x22, 0x0a, 0x1e, 0x57, 0x41, 0x4c, 0x4c, 0x45, 0x54, 0x5f, 0x45, 0x52, 0x52, 0x4f,
	0x52, 0x5f, 0x43, 0x4f, 0x44, 0x45, 0x5f, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x5f, 0x55,
	0x53, 0x45, 0x52, 0x10, 0x08, 0x12, 0x24, 0x0a, 0x20, 0x57, 0x41, 0x4c, 0x4c, 0x45, 0x54, 0x5f,
	0x45, 0x52, 0x52, 0x4f, 0x52, 0x5f, 0x43, 0x4f, 0x44, 0x45, 0x5f, 0x49, 0x4e, 0x56, 0x41, 0x4c,
	0x49, 0x44, 0x5f, 0x41, 0x4d, 0x4f, 0x55, 0x4e, 0x54, 0x10, 0x09, 0x32, 0xd6, 0x02, 0x0a, 0x14,
	0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x4b, 0x0a, 0x0c, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x57, 0x61,
	0x6c, 0x6c, 0x65, 0x74, 0x12, 0x1b, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1c, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x12, 0xb1, 0x01, 0x0a, 0x0b, 0x54, 0x6f, 0x70, 0x75, 0x70, 0x57, 0x61, 0x6c, 0x6c, 0x65,
	0x74, 0x12, 0x1a, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x6f, 0x70, 0x75, 0x70,
	0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x6f, 0x70, 0x75, 0x70, 0x57, 0x61, 0x6c, 0x6c,
	0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x69, 0x92, 0x41, 0x45, 0x0a,
	0x06, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x2a, 0x0b, 0x54, 0x6f, 0x70, 0x75, 0x70, 0x57, 0x61,
	0x6c, 0x6c, 0x65, 0x74, 0x72, 0x2e, 0x0a, 0x13, 0x0a, 0x0d, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72,
	0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x28, 0x01, 0x0a, 0x17, 0x0a, 0x11, 0x58,
	0x2d, 0x49, 0x64, 0x65, 0x6d, 0x70, 0x6f, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x2d, 0x4b, 0x65, 0x79,
	0x18, 0x01, 0x28, 0x01, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1b, 0x3a, 0x05, 0x74, 0x6f, 0x70, 0x75,
	0x70, 0x1a, 0x12, 0x2f, 0x76, 0x31, 0x2f, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x73, 0x2f, 0x74,
	0x6f, 0x70, 0x75, 0x70, 0x73, 0x1a, 0x3d, 0x92, 0x41, 0x3a, 0x12, 0x38, 0x54, 0x68, 0x69, 0x73,
	0x20, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x20, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65,
	0x73, 0x20, 0x61, 0x6c, 0x6c, 0x20, 0x75, 0x73, 0x65, 0x20, 0x63, 0x61, 0x73, 0x65, 0x73, 0x20,
	0x74, 0x6f, 0x20, 0x77, 0x6f, 0x72, 0x6b, 0x20, 0x77, 0x69, 0x74, 0x68, 0x20, 0x77, 0x61, 0x6c,
	0x6c, 0x65, 0x74, 0x2e, 0x42, 0x91, 0x02, 0x92, 0x41, 0xd1, 0x01, 0x12, 0x97, 0x01, 0x0a, 0x0a,
	0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x20, 0x41, 0x50, 0x49, 0x22, 0x30, 0x0a, 0x0d, 0x49, 0x6e,
	0x64, 0x72, 0x61, 0x20, 0x53, 0x61, 0x70, 0x75, 0x74, 0x72, 0x61, 0x12, 0x1f, 0x68, 0x74, 0x74,
	0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x69, 0x6e, 0x64, 0x72, 0x61, 0x73, 0x61, 0x70, 0x75, 0x74, 0x72, 0x61, 0x2a, 0x50, 0x0a, 0x14,
	0x42, 0x53, 0x44, 0x20, 0x33, 0x2d, 0x43, 0x6c, 0x61, 0x75, 0x73, 0x65, 0x20, 0x4c, 0x69, 0x63,
	0x65, 0x6e, 0x73, 0x65, 0x12, 0x38, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x6e, 0x64, 0x72, 0x61, 0x73, 0x61,
	0x70, 0x75, 0x74, 0x72, 0x61, 0x2f, 0x61, 0x72, 0x6a, 0x75, 0x6e, 0x61, 0x2f, 0x62, 0x6c, 0x6f,
	0x62, 0x2f, 0x6d, 0x61, 0x69, 0x6e, 0x2f, 0x4c, 0x49, 0x43, 0x45, 0x4e, 0x53, 0x45, 0x32, 0x05,
	0x31, 0x2e, 0x30, 0x2e, 0x30, 0x1a, 0x0e, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x68, 0x6f, 0x73, 0x74,
	0x3a, 0x38, 0x30, 0x30, 0x30, 0x2a, 0x01, 0x01, 0x32, 0x10, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x6a, 0x73, 0x6f, 0x6e, 0x3a, 0x10, 0x61, 0x70, 0x70, 0x6c,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x6a, 0x73, 0x6f, 0x6e, 0x5a, 0x3a, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x6e, 0x64, 0x72, 0x61, 0x73, 0x61,
	0x70, 0x75, 0x74, 0x72, 0x61, 0x2f, 0x61, 0x72, 0x6a, 0x75, 0x6e, 0x61, 0x2f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2f, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x76, 0x31, 0x3b, 0x61, 0x70, 0x69, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_v1_wallet_proto_rawDescOnce sync.Once
	file_api_v1_wallet_proto_rawDescData = file_api_v1_wallet_proto_rawDesc
)

func file_api_v1_wallet_proto_rawDescGZIP() []byte {
	file_api_v1_wallet_proto_rawDescOnce.Do(func() {
		file_api_v1_wallet_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_v1_wallet_proto_rawDescData)
	})
	return file_api_v1_wallet_proto_rawDescData
}

var file_api_v1_wallet_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_api_v1_wallet_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_api_v1_wallet_proto_goTypes = []any{
	(WalletErrorCode)(0),         // 0: api.v1.WalletErrorCode
	(*CreateWalletRequest)(nil),  // 1: api.v1.CreateWalletRequest
	(*CreateWalletResponse)(nil), // 2: api.v1.CreateWalletResponse
	(*TopupWalletRequest)(nil),   // 3: api.v1.TopupWalletRequest
	(*TopupWalletResponse)(nil),  // 4: api.v1.TopupWalletResponse
	(*Wallet)(nil),               // 5: api.v1.Wallet
	(*Topup)(nil),                // 6: api.v1.Topup
	(*WalletError)(nil),          // 7: api.v1.WalletError
}
var file_api_v1_wallet_proto_depIdxs = []int32{
	5, // 0: api.v1.CreateWalletRequest.wallet:type_name -> api.v1.Wallet
	6, // 1: api.v1.TopupWalletRequest.topup:type_name -> api.v1.Topup
	0, // 2: api.v1.WalletError.error_code:type_name -> api.v1.WalletErrorCode
	1, // 3: api.v1.WalletCommandService.CreateWallet:input_type -> api.v1.CreateWalletRequest
	3, // 4: api.v1.WalletCommandService.TopupWallet:input_type -> api.v1.TopupWalletRequest
	2, // 5: api.v1.WalletCommandService.CreateWallet:output_type -> api.v1.CreateWalletResponse
	4, // 6: api.v1.WalletCommandService.TopupWallet:output_type -> api.v1.TopupWalletResponse
	5, // [5:7] is the sub-list for method output_type
	3, // [3:5] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_api_v1_wallet_proto_init() }
func file_api_v1_wallet_proto_init() {
	if File_api_v1_wallet_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_v1_wallet_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*CreateWalletRequest); i {
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
		file_api_v1_wallet_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*CreateWalletResponse); i {
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
		file_api_v1_wallet_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*TopupWalletRequest); i {
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
		file_api_v1_wallet_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*TopupWalletResponse); i {
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
		file_api_v1_wallet_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*Wallet); i {
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
		file_api_v1_wallet_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*Topup); i {
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
		file_api_v1_wallet_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*WalletError); i {
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
			RawDescriptor: file_api_v1_wallet_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_v1_wallet_proto_goTypes,
		DependencyIndexes: file_api_v1_wallet_proto_depIdxs,
		EnumInfos:         file_api_v1_wallet_proto_enumTypes,
		MessageInfos:      file_api_v1_wallet_proto_msgTypes,
	}.Build()
	File_api_v1_wallet_proto = out.File
	file_api_v1_wallet_proto_rawDesc = nil
	file_api_v1_wallet_proto_goTypes = nil
	file_api_v1_wallet_proto_depIdxs = nil
}
