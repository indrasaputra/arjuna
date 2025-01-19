// auth.proto defines service for auth.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: api/v1/auth.proto

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

// AuthErrorCode enumerates auth error code.
type AuthErrorCode int32

const (
	// Default enum code according to
	// https://medium.com/@akhaku/protobuf-definition-best-practices-87f281576f31.
	AuthErrorCode_AUTH_ERROR_CODE_UNSPECIFIED AuthErrorCode = 0
	// Unexpected behavior occured in system.
	AuthErrorCode_AUTH_ERROR_CODE_INTERNAL AuthErrorCode = 1
	// Mandatory field is empty.
	AuthErrorCode_AUTH_ERROR_CODE_EMPTY_FIELD AuthErrorCode = 2
	// Unauthorized.
	AuthErrorCode_AUTH_ERROR_CODE_UNAUTHORIZED AuthErrorCode = 3
	// Some arguments must be invalid.
	AuthErrorCode_AUTH_ERROR_CODE_INVALID_ARGUMENT AuthErrorCode = 4
	// Account instance is empty or nil.
	AuthErrorCode_AUTH_ERROR_CODE_EMPTY_ACCOUNT AuthErrorCode = 5
	// Account's email is invalid.
	AuthErrorCode_AUTH_ERROR_CODE_INVALID_EMAIL AuthErrorCode = 6
	// Account's password is invalid.
	AuthErrorCode_AUTH_ERROR_CODE_INVALID_PASSWORD AuthErrorCode = 7
	// Account already exists.
	// The uniqueness of an account is represented by email.
	AuthErrorCode_AUTH_ERROR_CODE_ALREADY_EXISTS AuthErrorCode = 8
	// Account's credential is invalid.
	AuthErrorCode_AUTH_ERROR_CODE_INVALID_CREDENTIAL AuthErrorCode = 9
	// Data is not found.
	AuthErrorCode_AUTH_ERROR_CODE_NOT_FOUND AuthErrorCode = 10
)

// Enum value maps for AuthErrorCode.
var (
	AuthErrorCode_name = map[int32]string{
		0:  "AUTH_ERROR_CODE_UNSPECIFIED",
		1:  "AUTH_ERROR_CODE_INTERNAL",
		2:  "AUTH_ERROR_CODE_EMPTY_FIELD",
		3:  "AUTH_ERROR_CODE_UNAUTHORIZED",
		4:  "AUTH_ERROR_CODE_INVALID_ARGUMENT",
		5:  "AUTH_ERROR_CODE_EMPTY_ACCOUNT",
		6:  "AUTH_ERROR_CODE_INVALID_EMAIL",
		7:  "AUTH_ERROR_CODE_INVALID_PASSWORD",
		8:  "AUTH_ERROR_CODE_ALREADY_EXISTS",
		9:  "AUTH_ERROR_CODE_INVALID_CREDENTIAL",
		10: "AUTH_ERROR_CODE_NOT_FOUND",
	}
	AuthErrorCode_value = map[string]int32{
		"AUTH_ERROR_CODE_UNSPECIFIED":        0,
		"AUTH_ERROR_CODE_INTERNAL":           1,
		"AUTH_ERROR_CODE_EMPTY_FIELD":        2,
		"AUTH_ERROR_CODE_UNAUTHORIZED":       3,
		"AUTH_ERROR_CODE_INVALID_ARGUMENT":   4,
		"AUTH_ERROR_CODE_EMPTY_ACCOUNT":      5,
		"AUTH_ERROR_CODE_INVALID_EMAIL":      6,
		"AUTH_ERROR_CODE_INVALID_PASSWORD":   7,
		"AUTH_ERROR_CODE_ALREADY_EXISTS":     8,
		"AUTH_ERROR_CODE_INVALID_CREDENTIAL": 9,
		"AUTH_ERROR_CODE_NOT_FOUND":          10,
	}
)

func (x AuthErrorCode) Enum() *AuthErrorCode {
	p := new(AuthErrorCode)
	*p = x
	return p
}

func (x AuthErrorCode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AuthErrorCode) Descriptor() protoreflect.EnumDescriptor {
	return file_api_v1_auth_proto_enumTypes[0].Descriptor()
}

func (AuthErrorCode) Type() protoreflect.EnumType {
	return &file_api_v1_auth_proto_enumTypes[0]
}

func (x AuthErrorCode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AuthErrorCode.Descriptor instead.
func (AuthErrorCode) EnumDescriptor() ([]byte, []int) {
	return file_api_v1_auth_proto_rawDescGZIP(), []int{0}
}

// LoginRequest represents request for login.
type LoginRequest struct {
	state         protoimpl.MessageState
	Credential    *Credential `protobuf:"bytes,1,opt,name=credential,proto3" json:"credential,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LoginRequest) Reset() {
	*x = LoginRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_auth_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoginRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginRequest) ProtoMessage() {}

func (x *LoginRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_auth_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginRequest.ProtoReflect.Descriptor instead.
func (*LoginRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_auth_proto_rawDescGZIP(), []int{0}
}

func (x *LoginRequest) GetCredential() *Credential {
	if x != nil {
		return x.Credential
	}
	return nil
}

// LoginResponse represents response from login.
type LoginResponse struct {
	state         protoimpl.MessageState
	Data          *Token `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LoginResponse) Reset() {
	*x = LoginResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_auth_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoginResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginResponse) ProtoMessage() {}

func (x *LoginResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_auth_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginResponse.ProtoReflect.Descriptor instead.
func (*LoginResponse) Descriptor() ([]byte, []int) {
	return file_api_v1_auth_proto_rawDescGZIP(), []int{1}
}

func (x *LoginResponse) GetData() *Token {
	if x != nil {
		return x.Data
	}
	return nil
}

// RegisterAccountRequest represents request for account registration.
type RegisterAccountRequest struct {
	state         protoimpl.MessageState
	Account       *Account `protobuf:"bytes,1,opt,name=account,proto3" json:"account,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RegisterAccountRequest) Reset() {
	*x = RegisterAccountRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_auth_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterAccountRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterAccountRequest) ProtoMessage() {}

func (x *RegisterAccountRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_auth_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterAccountRequest.ProtoReflect.Descriptor instead.
func (*RegisterAccountRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_auth_proto_rawDescGZIP(), []int{2}
}

func (x *RegisterAccountRequest) GetAccount() *Account {
	if x != nil {
		return x.Account
	}
	return nil
}

// RegisterAccountResponse represents response for account registration.
type RegisterAccountResponse struct {
	state         protoimpl.MessageState
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RegisterAccountResponse) Reset() {
	*x = RegisterAccountResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_auth_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterAccountResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterAccountResponse) ProtoMessage() {}

func (x *RegisterAccountResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_auth_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterAccountResponse.ProtoReflect.Descriptor instead.
func (*RegisterAccountResponse) Descriptor() ([]byte, []int) {
	return file_api_v1_auth_proto_rawDescGZIP(), []int{3}
}

// Account represents account.
type Account struct {
	state         protoimpl.MessageState
	Id            string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId        string `protobuf:"bytes,2,opt,name=user_id,proto3" json:"user_id,omitempty"`
	Email         string `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	Password      string `protobuf:"bytes,4,opt,name=password,proto3" json:"password,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Account) Reset() {
	*x = Account{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_auth_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Account) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Account) ProtoMessage() {}

func (x *Account) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_auth_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Account.ProtoReflect.Descriptor instead.
func (*Account) Descriptor() ([]byte, []int) {
	return file_api_v1_auth_proto_rawDescGZIP(), []int{4}
}

func (x *Account) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Account) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *Account) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *Account) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

// Credential represents login credential.
type Credential struct {
	state         protoimpl.MessageState
	Email         string `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	Password      string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Credential) Reset() {
	*x = Credential{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_auth_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Credential) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Credential) ProtoMessage() {}

func (x *Credential) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_auth_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Credential.ProtoReflect.Descriptor instead.
func (*Credential) Descriptor() ([]byte, []int) {
	return file_api_v1_auth_proto_rawDescGZIP(), []int{5}
}

func (x *Credential) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *Credential) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

// Token represents token.
type Token struct {
	state                 protoimpl.MessageState
	AccessToken           string `protobuf:"bytes,1,opt,name=access_token,proto3" json:"access_token,omitempty"`
	RefreshToken          string `protobuf:"bytes,3,opt,name=refresh_token,proto3" json:"refresh_token,omitempty"`
	unknownFields         protoimpl.UnknownFields
	sizeCache             protoimpl.SizeCache
	AccessTokenExpiresIn  uint32 `protobuf:"varint,2,opt,name=access_token_expires_in,proto3" json:"access_token_expires_in,omitempty"`
	RefreshTokenExpiresIn uint32 `protobuf:"varint,4,opt,name=refresh_token_expires_in,proto3" json:"refresh_token_expires_in,omitempty"`
}

func (x *Token) Reset() {
	*x = Token{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_auth_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Token) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Token) ProtoMessage() {}

func (x *Token) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_auth_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Token.ProtoReflect.Descriptor instead.
func (*Token) Descriptor() ([]byte, []int) {
	return file_api_v1_auth_proto_rawDescGZIP(), []int{6}
}

func (x *Token) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

func (x *Token) GetAccessTokenExpiresIn() uint32 {
	if x != nil {
		return x.AccessTokenExpiresIn
	}
	return 0
}

func (x *Token) GetRefreshToken() string {
	if x != nil {
		return x.RefreshToken
	}
	return ""
}

func (x *Token) GetRefreshTokenExpiresIn() uint32 {
	if x != nil {
		return x.RefreshTokenExpiresIn
	}
	return 0
}

// AuthError represents message for any error happening in auth service.
type AuthError struct {
	state         protoimpl.MessageState
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
	ErrorCode     AuthErrorCode `protobuf:"varint,1,opt,name=error_code,json=errorCode,proto3,enum=api.v1.AuthErrorCode" json:"error_code,omitempty"`
}

func (x *AuthError) Reset() {
	*x = AuthError{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_auth_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthError) ProtoMessage() {}

func (x *AuthError) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_auth_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthError.ProtoReflect.Descriptor instead.
func (*AuthError) Descriptor() ([]byte, []int) {
	return file_api_v1_auth_proto_rawDescGZIP(), []int{7}
}

func (x *AuthError) GetErrorCode() AuthErrorCode {
	if x != nil {
		return x.ErrorCode
	}
	return AuthErrorCode_AUTH_ERROR_CODE_UNSPECIFIED
}

var File_api_v1_auth_proto protoreflect.FileDescriptor

var file_api_v1_auth_proto_rawDesc = []byte{
	0x0a, 0x11, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x06, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x62, 0x65, 0x68, 0x61,
	0x76, 0x69, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32,
	0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x47, 0x0a, 0x0c, 0x4c, 0x6f,
	0x67, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x37, 0x0a, 0x0a, 0x63, 0x72,
	0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69,
	0x61, 0x6c, 0x42, 0x03, 0xe0, 0x41, 0x02, 0x52, 0x0a, 0x63, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74,
	0x69, 0x61, 0x6c, 0x22, 0x32, 0x0a, 0x0d, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x21, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x43, 0x0a, 0x16, 0x52, 0x65, 0x67, 0x69, 0x73,
	0x74, 0x65, 0x72, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x29, 0x0a, 0x07, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x63, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x52, 0x07, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x19, 0x0a, 0x17,
	0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0xd8, 0x02, 0x0a, 0x07, 0x41, 0x63, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x12, 0x41, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x31, 0x92, 0x41, 0x28, 0x4a, 0x26, 0x22, 0x30, 0x31, 0x39, 0x31, 0x37, 0x61, 0x30, 0x63, 0x2d,
	0x34, 0x37, 0x35, 0x65, 0x2d, 0x37, 0x64, 0x34, 0x61, 0x2d, 0x39, 0x65, 0x63, 0x31, 0x2d, 0x61,
	0x35, 0x36, 0x64, 0x31, 0x34, 0x64, 0x37, 0x38, 0x35, 0x36, 0x39, 0x22, 0xe0, 0x41, 0x02, 0xe0,
	0x41, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x5c, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x42, 0x92, 0x41, 0x3c, 0x32, 0x09, 0x55, 0x73,
	0x65, 0x72, 0x27, 0x73, 0x20, 0x69, 0x64, 0x4a, 0x26, 0x22, 0x30, 0x31, 0x39, 0x31, 0x37, 0x61,
	0x30, 0x63, 0x2d, 0x63, 0x64, 0x66, 0x65, 0x2d, 0x37, 0x36, 0x39, 0x36, 0x2d, 0x39, 0x32, 0x35,
	0x39, 0x2d, 0x31, 0x63, 0x39, 0x38, 0x32, 0x31, 0x63, 0x38, 0x66, 0x64, 0x37, 0x33, 0x22, 0xa2,
	0x02, 0x06, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0xe0, 0x41, 0x02, 0x52, 0x07, 0x75, 0x73, 0x65,
	0x72, 0x5f, 0x69, 0x64, 0x12, 0x61, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x4b, 0x92, 0x41, 0x45, 0x32, 0x0c, 0x55, 0x73, 0x65, 0x72, 0x27, 0x73,
	0x20, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x4a, 0x12, 0x22, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x40, 0x64,
	0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x63, 0x6f, 0x6d, 0x22, 0x8a, 0x01, 0x20, 0x5e, 0x5b, 0x5c,
	0x77, 0x2d, 0x5c, 0x2e, 0x5d, 0x2b, 0x40, 0x28, 0x5b, 0x5c, 0x77, 0x2d, 0x5d, 0x2b, 0x5c, 0x2e,
	0x29, 0x2b, 0x5b, 0x5c, 0x77, 0x2d, 0x5d, 0x7b, 0x32, 0x2c, 0x34, 0x7d, 0x24, 0xe0, 0x41, 0x02,
	0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x49, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77,
	0x6f, 0x72, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x2d, 0x92, 0x41, 0x2a, 0x32, 0x0f,
	0x55, 0x73, 0x65, 0x72, 0x27, 0x73, 0x20, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x4a,
	0x0e, 0x22, 0x77, 0x65, 0x61, 0x6b, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x22, 0xa2,
	0x02, 0x06, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f,
	0x72, 0x64, 0x22, 0xbd, 0x01, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61,
	0x6c, 0x12, 0x61, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x4b, 0x92, 0x41, 0x45, 0x32, 0x0c, 0x55, 0x73, 0x65, 0x72, 0x27, 0x73, 0x20, 0x65, 0x6d,
	0x61, 0x69, 0x6c, 0x4a, 0x12, 0x22, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x40, 0x64, 0x6f, 0x6d, 0x61,
	0x69, 0x6e, 0x2e, 0x63, 0x6f, 0x6d, 0x22, 0x8a, 0x01, 0x20, 0x5e, 0x5b, 0x5c, 0x77, 0x2d, 0x5c,
	0x2e, 0x5d, 0x2b, 0x40, 0x28, 0x5b, 0x5c, 0x77, 0x2d, 0x5d, 0x2b, 0x5c, 0x2e, 0x29, 0x2b, 0x5b,
	0x5c, 0x77, 0x2d, 0x5d, 0x7b, 0x32, 0x2c, 0x34, 0x7d, 0x24, 0xe0, 0x41, 0x02, 0x52, 0x05, 0x65,
	0x6d, 0x61, 0x69, 0x6c, 0x12, 0x4c, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x30, 0x92, 0x41, 0x2a, 0x32, 0x0f, 0x55, 0x73, 0x65,
	0x72, 0x27, 0x73, 0x20, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x4a, 0x0e, 0x22, 0x77,
	0x65, 0x61, 0x6b, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x22, 0xa2, 0x02, 0x06, 0x73,
	0x74, 0x72, 0x69, 0x6e, 0x67, 0xe0, 0x41, 0x02, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f,
	0x72, 0x64, 0x22, 0xe7, 0x01, 0x0a, 0x05, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x2a, 0x0a, 0x0c,
	0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x06, 0xe0, 0x41, 0x02, 0xe0, 0x41, 0x03, 0x52, 0x0c, 0x61, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x40, 0x0a, 0x17, 0x61, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73,
	0x5f, 0x69, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x42, 0x06, 0xe0, 0x41, 0x02, 0xe0, 0x41,
	0x03, 0x52, 0x17, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f,
	0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x5f, 0x69, 0x6e, 0x12, 0x2c, 0x0a, 0x0d, 0x72, 0x65,
	0x66, 0x72, 0x65, 0x73, 0x68, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x06, 0xe0, 0x41, 0x02, 0xe0, 0x41, 0x03, 0x52, 0x0d, 0x72, 0x65, 0x66, 0x72, 0x65,
	0x73, 0x68, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x42, 0x0a, 0x18, 0x72, 0x65, 0x66, 0x72,
	0x65, 0x73, 0x68, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65,
	0x73, 0x5f, 0x69, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x42, 0x06, 0xe0, 0x41, 0x02, 0xe0,
	0x41, 0x03, 0x52, 0x18, 0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x5f, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x5f, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x5f, 0x69, 0x6e, 0x22, 0x41, 0x0a, 0x09,
	0x41, 0x75, 0x74, 0x68, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x34, 0x0a, 0x0a, 0x65, 0x72, 0x72,
	0x6f, 0x72, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x15, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x45, 0x72, 0x72, 0x6f, 0x72,
	0x43, 0x6f, 0x64, 0x65, 0x52, 0x09, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x2a,
	0x8e, 0x03, 0x0a, 0x0d, 0x41, 0x75, 0x74, 0x68, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64,
	0x65, 0x12, 0x1f, 0x0a, 0x1b, 0x41, 0x55, 0x54, 0x48, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x5f,
	0x43, 0x4f, 0x44, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44,
	0x10, 0x00, 0x12, 0x1c, 0x0a, 0x18, 0x41, 0x55, 0x54, 0x48, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52,
	0x5f, 0x43, 0x4f, 0x44, 0x45, 0x5f, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x4e, 0x41, 0x4c, 0x10, 0x01,
	0x12, 0x1f, 0x0a, 0x1b, 0x41, 0x55, 0x54, 0x48, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x5f, 0x43,
	0x4f, 0x44, 0x45, 0x5f, 0x45, 0x4d, 0x50, 0x54, 0x59, 0x5f, 0x46, 0x49, 0x45, 0x4c, 0x44, 0x10,
	0x02, 0x12, 0x20, 0x0a, 0x1c, 0x41, 0x55, 0x54, 0x48, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x5f,
	0x43, 0x4f, 0x44, 0x45, 0x5f, 0x55, 0x4e, 0x41, 0x55, 0x54, 0x48, 0x4f, 0x52, 0x49, 0x5a, 0x45,
	0x44, 0x10, 0x03, 0x12, 0x24, 0x0a, 0x20, 0x41, 0x55, 0x54, 0x48, 0x5f, 0x45, 0x52, 0x52, 0x4f,
	0x52, 0x5f, 0x43, 0x4f, 0x44, 0x45, 0x5f, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x5f, 0x41,
	0x52, 0x47, 0x55, 0x4d, 0x45, 0x4e, 0x54, 0x10, 0x04, 0x12, 0x21, 0x0a, 0x1d, 0x41, 0x55, 0x54,
	0x48, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x5f, 0x43, 0x4f, 0x44, 0x45, 0x5f, 0x45, 0x4d, 0x50,
	0x54, 0x59, 0x5f, 0x41, 0x43, 0x43, 0x4f, 0x55, 0x4e, 0x54, 0x10, 0x05, 0x12, 0x21, 0x0a, 0x1d,
	0x41, 0x55, 0x54, 0x48, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x5f, 0x43, 0x4f, 0x44, 0x45, 0x5f,
	0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x5f, 0x45, 0x4d, 0x41, 0x49, 0x4c, 0x10, 0x06, 0x12,
	0x24, 0x0a, 0x20, 0x41, 0x55, 0x54, 0x48, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x5f, 0x43, 0x4f,
	0x44, 0x45, 0x5f, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x5f, 0x50, 0x41, 0x53, 0x53, 0x57,
	0x4f, 0x52, 0x44, 0x10, 0x07, 0x12, 0x22, 0x0a, 0x1e, 0x41, 0x55, 0x54, 0x48, 0x5f, 0x45, 0x52,
	0x52, 0x4f, 0x52, 0x5f, 0x43, 0x4f, 0x44, 0x45, 0x5f, 0x41, 0x4c, 0x52, 0x45, 0x41, 0x44, 0x59,
	0x5f, 0x45, 0x58, 0x49, 0x53, 0x54, 0x53, 0x10, 0x08, 0x12, 0x26, 0x0a, 0x22, 0x41, 0x55, 0x54,
	0x48, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x5f, 0x43, 0x4f, 0x44, 0x45, 0x5f, 0x49, 0x4e, 0x56,
	0x41, 0x4c, 0x49, 0x44, 0x5f, 0x43, 0x52, 0x45, 0x44, 0x45, 0x4e, 0x54, 0x49, 0x41, 0x4c, 0x10,
	0x09, 0x12, 0x1d, 0x0a, 0x19, 0x41, 0x55, 0x54, 0x48, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x5f,
	0x43, 0x4f, 0x44, 0x45, 0x5f, 0x4e, 0x4f, 0x54, 0x5f, 0x46, 0x4f, 0x55, 0x4e, 0x44, 0x10, 0x0a,
	0x32, 0x8a, 0x02, 0x0a, 0x0b, 0x41, 0x75, 0x74, 0x68, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x68, 0x0a, 0x05, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x14, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x76, 0x31, 0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x15, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x32, 0x92, 0x41, 0x0d, 0x0a, 0x04, 0x41, 0x75, 0x74,
	0x68, 0x2a, 0x05, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1c, 0x3a, 0x0a,
	0x63, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x22, 0x0e, 0x2f, 0x76, 0x31, 0x2f,
	0x61, 0x75, 0x74, 0x68, 0x2f, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x54, 0x0a, 0x0f, 0x52, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1e, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x41,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x41,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x1a, 0x3b, 0x92, 0x41, 0x38, 0x12, 0x36, 0x54, 0x68, 0x69, 0x73, 0x20, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x20, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x73, 0x20, 0x61, 0x6c, 0x6c,
	0x20, 0x75, 0x73, 0x65, 0x20, 0x63, 0x61, 0x73, 0x65, 0x73, 0x20, 0x74, 0x6f, 0x20, 0x77, 0x6f,
	0x72, 0x6b, 0x20, 0x77, 0x69, 0x74, 0x68, 0x20, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x42, 0x8d, 0x02,
	0x92, 0x41, 0xcf, 0x01, 0x12, 0x95, 0x01, 0x0a, 0x08, 0x41, 0x75, 0x74, 0x68, 0x20, 0x41, 0x50,
	0x49, 0x22, 0x30, 0x0a, 0x0d, 0x49, 0x6e, 0x64, 0x72, 0x61, 0x20, 0x53, 0x61, 0x70, 0x75, 0x74,
	0x72, 0x61, 0x12, 0x1f, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x6e, 0x64, 0x72, 0x61, 0x73, 0x61, 0x70, 0x75,
	0x74, 0x72, 0x61, 0x2a, 0x50, 0x0a, 0x14, 0x42, 0x53, 0x44, 0x20, 0x33, 0x2d, 0x43, 0x6c, 0x61,
	0x75, 0x73, 0x65, 0x20, 0x4c, 0x69, 0x63, 0x65, 0x6e, 0x73, 0x65, 0x12, 0x38, 0x68, 0x74, 0x74,
	0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x69, 0x6e, 0x64, 0x72, 0x61, 0x73, 0x61, 0x70, 0x75, 0x74, 0x72, 0x61, 0x2f, 0x61, 0x72, 0x6a,
	0x75, 0x6e, 0x61, 0x2f, 0x62, 0x6c, 0x6f, 0x62, 0x2f, 0x6d, 0x61, 0x69, 0x6e, 0x2f, 0x4c, 0x49,
	0x43, 0x45, 0x4e, 0x53, 0x45, 0x32, 0x05, 0x31, 0x2e, 0x30, 0x2e, 0x30, 0x1a, 0x0e, 0x6c, 0x6f,
	0x63, 0x61, 0x6c, 0x68, 0x6f, 0x73, 0x74, 0x3a, 0x38, 0x30, 0x30, 0x30, 0x2a, 0x01, 0x01, 0x32,
	0x10, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x6a, 0x73, 0x6f,
	0x6e, 0x3a, 0x10, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x6a,
	0x73, 0x6f, 0x6e, 0x5a, 0x38, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x69, 0x6e, 0x64, 0x72, 0x61, 0x73, 0x61, 0x70, 0x75, 0x74, 0x72, 0x61, 0x2f, 0x61, 0x72, 0x6a,
	0x75, 0x6e, 0x61, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x61, 0x75, 0x74, 0x68,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x3b, 0x61, 0x70, 0x69, 0x76, 0x31, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_v1_auth_proto_rawDescOnce sync.Once
	file_api_v1_auth_proto_rawDescData = file_api_v1_auth_proto_rawDesc
)

func file_api_v1_auth_proto_rawDescGZIP() []byte {
	file_api_v1_auth_proto_rawDescOnce.Do(func() {
		file_api_v1_auth_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_v1_auth_proto_rawDescData)
	})
	return file_api_v1_auth_proto_rawDescData
}

var file_api_v1_auth_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_api_v1_auth_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_api_v1_auth_proto_goTypes = []any{
	(AuthErrorCode)(0),              // 0: api.v1.AuthErrorCode
	(*LoginRequest)(nil),            // 1: api.v1.LoginRequest
	(*LoginResponse)(nil),           // 2: api.v1.LoginResponse
	(*RegisterAccountRequest)(nil),  // 3: api.v1.RegisterAccountRequest
	(*RegisterAccountResponse)(nil), // 4: api.v1.RegisterAccountResponse
	(*Account)(nil),                 // 5: api.v1.Account
	(*Credential)(nil),              // 6: api.v1.Credential
	(*Token)(nil),                   // 7: api.v1.Token
	(*AuthError)(nil),               // 8: api.v1.AuthError
}
var file_api_v1_auth_proto_depIdxs = []int32{
	6, // 0: api.v1.LoginRequest.credential:type_name -> api.v1.Credential
	7, // 1: api.v1.LoginResponse.data:type_name -> api.v1.Token
	5, // 2: api.v1.RegisterAccountRequest.account:type_name -> api.v1.Account
	0, // 3: api.v1.AuthError.error_code:type_name -> api.v1.AuthErrorCode
	1, // 4: api.v1.AuthService.Login:input_type -> api.v1.LoginRequest
	3, // 5: api.v1.AuthService.RegisterAccount:input_type -> api.v1.RegisterAccountRequest
	2, // 6: api.v1.AuthService.Login:output_type -> api.v1.LoginResponse
	4, // 7: api.v1.AuthService.RegisterAccount:output_type -> api.v1.RegisterAccountResponse
	6, // [6:8] is the sub-list for method output_type
	4, // [4:6] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_api_v1_auth_proto_init() }
func file_api_v1_auth_proto_init() {
	if File_api_v1_auth_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_v1_auth_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*LoginRequest); i {
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
		file_api_v1_auth_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*LoginResponse); i {
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
		file_api_v1_auth_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*RegisterAccountRequest); i {
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
		file_api_v1_auth_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*RegisterAccountResponse); i {
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
		file_api_v1_auth_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*Account); i {
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
		file_api_v1_auth_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*Credential); i {
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
		file_api_v1_auth_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*Token); i {
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
		file_api_v1_auth_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*AuthError); i {
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
			RawDescriptor: file_api_v1_auth_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_v1_auth_proto_goTypes,
		DependencyIndexes: file_api_v1_auth_proto_depIdxs,
		EnumInfos:         file_api_v1_auth_proto_enumTypes,
		MessageInfos:      file_api_v1_auth_proto_msgTypes,
	}.Build()
	File_api_v1_auth_proto = out.File
	file_api_v1_auth_proto_rawDesc = nil
	file_api_v1_auth_proto_goTypes = nil
	file_api_v1_auth_proto_depIdxs = nil
}
