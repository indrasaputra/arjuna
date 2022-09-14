// user.proto defines service for user.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        (unknown)
// source: api/v1/user.proto

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

// UserErrorCode enumerates user error code.
type UserErrorCode int32

const (
	// Default enum code according to
	// https://medium.com/@akhaku/protobuf-definition-best-practices-87f281576f31.
	UserErrorCode_USER_ERROR_CODE_UNSPECIFIED UserErrorCode = 0
	// Unexpected behavior occured in system.
	UserErrorCode_USER_ERROR_CODE_INTERNAL UserErrorCode = 1
	// User instance is empty or nil.
	UserErrorCode_USER_ERROR_CODE_EMPTY_USER UserErrorCode = 2
	// User already exists.
	// The uniqueness of a user is represented by email.
	UserErrorCode_USER_ERROR_CODE_ALREADY_EXISTS UserErrorCode = 3
	// User's name is invalid.
	// Allowed characters are alphabet only.
	UserErrorCode_USER_ERROR_CODE_INVALID_NAME UserErrorCode = 4
	// User's email is invalid.
	UserErrorCode_USER_ERROR_CODE_INVALID_EMAIL UserErrorCode = 5
	// User' not found.
	UserErrorCode_USER_ERROR_CODE_NOT_FOUND UserErrorCode = 6
)

// Enum value maps for UserErrorCode.
var (
	UserErrorCode_name = map[int32]string{
		0: "USER_ERROR_CODE_UNSPECIFIED",
		1: "USER_ERROR_CODE_INTERNAL",
		2: "USER_ERROR_CODE_EMPTY_USER",
		3: "USER_ERROR_CODE_ALREADY_EXISTS",
		4: "USER_ERROR_CODE_INVALID_NAME",
		5: "USER_ERROR_CODE_INVALID_EMAIL",
		6: "USER_ERROR_CODE_NOT_FOUND",
	}
	UserErrorCode_value = map[string]int32{
		"USER_ERROR_CODE_UNSPECIFIED":    0,
		"USER_ERROR_CODE_INTERNAL":       1,
		"USER_ERROR_CODE_EMPTY_USER":     2,
		"USER_ERROR_CODE_ALREADY_EXISTS": 3,
		"USER_ERROR_CODE_INVALID_NAME":   4,
		"USER_ERROR_CODE_INVALID_EMAIL":  5,
		"USER_ERROR_CODE_NOT_FOUND":      6,
	}
)

func (x UserErrorCode) Enum() *UserErrorCode {
	p := new(UserErrorCode)
	*p = x
	return p
}

func (x UserErrorCode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (UserErrorCode) Descriptor() protoreflect.EnumDescriptor {
	return file_api_v1_user_proto_enumTypes[0].Descriptor()
}

func (UserErrorCode) Type() protoreflect.EnumType {
	return &file_api_v1_user_proto_enumTypes[0]
}

func (x UserErrorCode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use UserErrorCode.Descriptor instead.
func (UserErrorCode) EnumDescriptor() ([]byte, []int) {
	return file_api_v1_user_proto_rawDescGZIP(), []int{0}
}

// RegisterUserRequest represents request for register user.
type RegisterUserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// user represents user data.
	User *User `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
}

func (x *RegisterUserRequest) Reset() {
	*x = RegisterUserRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_user_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterUserRequest) ProtoMessage() {}

func (x *RegisterUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_user_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterUserRequest.ProtoReflect.Descriptor instead.
func (*RegisterUserRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_user_proto_rawDescGZIP(), []int{0}
}

func (x *RegisterUserRequest) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

// RegisterUserResponse represents response from register user.
type RegisterUserResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// data represents user.
	Data *User `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *RegisterUserResponse) Reset() {
	*x = RegisterUserResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_user_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterUserResponse) ProtoMessage() {}

func (x *RegisterUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_user_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterUserResponse.ProtoReflect.Descriptor instead.
func (*RegisterUserResponse) Descriptor() ([]byte, []int) {
	return file_api_v1_user_proto_rawDescGZIP(), []int{1}
}

func (x *RegisterUserResponse) GetData() *User {
	if x != nil {
		return x.Data
	}
	return nil
}

// DeleteUserRequest represents request for delete user.
type DeleteUserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// id represents user's id.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *DeleteUserRequest) Reset() {
	*x = DeleteUserRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_user_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteUserRequest) ProtoMessage() {}

func (x *DeleteUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_user_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteUserRequest.ProtoReflect.Descriptor instead.
func (*DeleteUserRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_user_proto_rawDescGZIP(), []int{2}
}

func (x *DeleteUserRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

// DeleteUserResponse represents response from delete user.
type DeleteUserResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeleteUserResponse) Reset() {
	*x = DeleteUserResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_user_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteUserResponse) ProtoMessage() {}

func (x *DeleteUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_user_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteUserResponse.ProtoReflect.Descriptor instead.
func (*DeleteUserResponse) Descriptor() ([]byte, []int) {
	return file_api_v1_user_proto_rawDescGZIP(), []int{3}
}

// GetAllUsersRequest represents request for get all users.
type GetAllUsersRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// limit specifies how many users to retrieve in a single call.
	Limit uint32 `protobuf:"varint,1,opt,name=limit,proto3" json:"limit,omitempty"`
}

func (x *GetAllUsersRequest) Reset() {
	*x = GetAllUsersRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_user_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAllUsersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAllUsersRequest) ProtoMessage() {}

func (x *GetAllUsersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_user_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAllUsersRequest.ProtoReflect.Descriptor instead.
func (*GetAllUsersRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_user_proto_rawDescGZIP(), []int{4}
}

func (x *GetAllUsersRequest) GetLimit() uint32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

// GetAllUsersResponse represents response from get all users.
type GetAllUsersResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// data represents an array of user data.
	Data []*User `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
}

func (x *GetAllUsersResponse) Reset() {
	*x = GetAllUsersResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_user_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAllUsersResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAllUsersResponse) ProtoMessage() {}

func (x *GetAllUsersResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_user_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAllUsersResponse.ProtoReflect.Descriptor instead.
func (*GetAllUsersResponse) Descriptor() ([]byte, []int) {
	return file_api_v1_user_proto_rawDescGZIP(), []int{5}
}

func (x *GetAllUsersResponse) GetData() []*User {
	if x != nil {
		return x.Data
	}
	return nil
}

// User represents a user data.
type User struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// id represents a user's id.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// email represents a user's email.
	Email string `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	// password represents the user's password.
	Password string `protobuf:"bytes,3,opt,name=password,proto3" json:"password,omitempty"`
	// name represents a user's name.
	Name string `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	// created_at represents when the user was registered.
	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// updated_at represents when the user was last updated.
	UpdatedAt *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

func (x *User) Reset() {
	*x = User{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_user_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_user_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_api_v1_user_proto_rawDescGZIP(), []int{6}
}

func (x *User) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *User) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *User) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *User) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *User) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *User) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

// UserError represents message for any error happening in user.
type UserError struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// error_code represents specific and unique error code for user.
	ErrorCode UserErrorCode `protobuf:"varint,1,opt,name=error_code,json=errorCode,proto3,enum=api.v1.UserErrorCode" json:"error_code,omitempty"`
}

func (x *UserError) Reset() {
	*x = UserError{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_user_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserError) ProtoMessage() {}

func (x *UserError) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_user_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserError.ProtoReflect.Descriptor instead.
func (*UserError) Descriptor() ([]byte, []int) {
	return file_api_v1_user_proto_rawDescGZIP(), []int{7}
}

func (x *UserError) GetErrorCode() UserErrorCode {
	if x != nil {
		return x.ErrorCode
	}
	return UserErrorCode_USER_ERROR_CODE_UNSPECIFIED
}

var File_api_v1_user_proto protoreflect.FileDescriptor

var file_api_v1_user_proto_rawDesc = []byte{
	0x0a, 0x11, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x06, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x62, 0x65, 0x68, 0x61,
	0x76, 0x69, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76,
	0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x37, 0x0a, 0x13, 0x52,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x20, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0c, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x04,
	0x75, 0x73, 0x65, 0x72, 0x22, 0x38, 0x0a, 0x14, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72,
	0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x20, 0x0a, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x76, 0x31, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x23,
	0x0a, 0x11, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x22, 0x14, 0x0a, 0x12, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x55, 0x73, 0x65,
	0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x2a, 0x0a, 0x12, 0x47, 0x65, 0x74,
	0x41, 0x6c, 0x6c, 0x55, 0x73, 0x65, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05,
	0x6c, 0x69, 0x6d, 0x69, 0x74, 0x22, 0x37, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x55,
	0x73, 0x65, 0x72, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x20, 0x0a, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x76, 0x31, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x88,
	0x03, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x04, 0xe2, 0x41, 0x01, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x67, 0x0a,
	0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x51, 0x92, 0x41,
	0x4e, 0x32, 0x0c, 0x75, 0x73, 0x65, 0x72, 0x27, 0x73, 0x20, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x4a,
	0x13, 0x22, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x40, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x6e, 0x79, 0x2e,
	0x63, 0x6f, 0x6d, 0x22, 0x8a, 0x01, 0x20, 0x5e, 0x5b, 0x5c, 0x77, 0x2d, 0x5c, 0x2e, 0x5d, 0x2b,
	0x40, 0x28, 0x5b, 0x5c, 0x77, 0x2d, 0x5d, 0x2b, 0x5c, 0x2e, 0x29, 0x2b, 0x5b, 0x5c, 0x77, 0x2d,
	0x5d, 0x7b, 0x32, 0x2c, 0x34, 0x7d, 0x24, 0xd2, 0x01, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x52,
	0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x3d, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f,
	0x72, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x21, 0x92, 0x41, 0x1a, 0x32, 0x0f, 0x75,
	0x73, 0x65, 0x72, 0x27, 0x73, 0x20, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0xa2, 0x02,
	0x06, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0xe2, 0x41, 0x01, 0x04, 0x52, 0x08, 0x70, 0x61, 0x73,
	0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x40, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x2c, 0x92, 0x41, 0x29, 0x32, 0x0b, 0x75, 0x73, 0x65, 0x72, 0x27, 0x73,
	0x20, 0x6e, 0x61, 0x6d, 0x65, 0x4a, 0x14, 0x22, 0x5a, 0x6c, 0x61, 0x74, 0x61, 0x6e, 0x20, 0x49,
	0x62, 0x72, 0x61, 0x68, 0x69, 0x6d, 0x6f, 0x76, 0x69, 0x63, 0x22, 0x78, 0xff, 0x01, 0x80, 0x01,
	0x01, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x3f, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x42, 0x04, 0xe2, 0x41, 0x01, 0x03, 0x52, 0x09, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x3f, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x42, 0x04, 0xe2, 0x41, 0x01, 0x03, 0x52, 0x09,
	0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0x41, 0x0a, 0x09, 0x55, 0x73, 0x65,
	0x72, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x34, 0x0a, 0x0a, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f,
	0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x15, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x76, 0x31, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64,
	0x65, 0x52, 0x09, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x2a, 0xf6, 0x01, 0x0a,
	0x0d, 0x55, 0x73, 0x65, 0x72, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x1f,
	0x0a, 0x1b, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x5f, 0x43, 0x4f, 0x44,
	0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12,
	0x1c, 0x0a, 0x18, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x5f, 0x43, 0x4f,
	0x44, 0x45, 0x5f, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x4e, 0x41, 0x4c, 0x10, 0x01, 0x12, 0x1e, 0x0a,
	0x1a, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x5f, 0x43, 0x4f, 0x44, 0x45,
	0x5f, 0x45, 0x4d, 0x50, 0x54, 0x59, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x10, 0x02, 0x12, 0x22, 0x0a,
	0x1e, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x5f, 0x43, 0x4f, 0x44, 0x45,
	0x5f, 0x41, 0x4c, 0x52, 0x45, 0x41, 0x44, 0x59, 0x5f, 0x45, 0x58, 0x49, 0x53, 0x54, 0x53, 0x10,
	0x03, 0x12, 0x20, 0x0a, 0x1c, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x5f,
	0x43, 0x4f, 0x44, 0x45, 0x5f, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x5f, 0x4e, 0x41, 0x4d,
	0x45, 0x10, 0x04, 0x12, 0x21, 0x0a, 0x1d, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x45, 0x52, 0x52, 0x4f,
	0x52, 0x5f, 0x43, 0x4f, 0x44, 0x45, 0x5f, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x5f, 0x45,
	0x4d, 0x41, 0x49, 0x4c, 0x10, 0x05, 0x12, 0x1d, 0x0a, 0x19, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x45,
	0x52, 0x52, 0x4f, 0x52, 0x5f, 0x43, 0x4f, 0x44, 0x45, 0x5f, 0x4e, 0x4f, 0x54, 0x5f, 0x46, 0x4f,
	0x55, 0x4e, 0x44, 0x10, 0x06, 0x32, 0xb0, 0x02, 0x0a, 0x12, 0x55, 0x73, 0x65, 0x72, 0x43, 0x6f,
	0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x82, 0x01, 0x0a,
	0x0c, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x55, 0x73, 0x65, 0x72, 0x12, 0x1b, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x55,
	0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x55, 0x73, 0x65, 0x72,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x37, 0x92, 0x41, 0x14, 0x0a, 0x04, 0x55,
	0x73, 0x65, 0x72, 0x2a, 0x0c, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x55, 0x73, 0x65,
	0x72, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1a, 0x22, 0x12, 0x2f, 0x76, 0x31, 0x2f, 0x75, 0x73, 0x65,
	0x72, 0x73, 0x2f, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x3a, 0x04, 0x75, 0x73, 0x65,
	0x72, 0x1a, 0x94, 0x01, 0x92, 0x41, 0x90, 0x01, 0x12, 0x8d, 0x01, 0x54, 0x68, 0x69, 0x73, 0x20,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x20, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x73,
	0x20, 0x62, 0x61, 0x73, 0x69, 0x63, 0x20, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x20, 0x6f,
	0x72, 0x20, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2d, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x69, 0x6e, 0x67,
	0x20, 0x75, 0x73, 0x65, 0x20, 0x63, 0x61, 0x73, 0x65, 0x73, 0x20, 0x74, 0x6f, 0x20, 0x77, 0x6f,
	0x72, 0x6b, 0x20, 0x77, 0x69, 0x74, 0x68, 0x20, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x41, 0x20, 0x75,
	0x73, 0x65, 0x72, 0x20, 0x69, 0x73, 0x20, 0x72, 0x65, 0x70, 0x72, 0x65, 0x73, 0x65, 0x6e, 0x74,
	0x65, 0x64, 0x20, 0x62, 0x79, 0x20, 0x61, 0x6e, 0x20, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x20, 0x61,
	0x73, 0x20, 0x69, 0x74, 0x73, 0x20, 0x75, 0x6e, 0x69, 0x71, 0x75, 0x65, 0x20, 0x69, 0x64, 0x65,
	0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x2e, 0x32, 0xeb, 0x01, 0x0a, 0x1a, 0x55, 0x73, 0x65,
	0x72, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x70, 0x0a, 0x0a, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x55, 0x73, 0x65, 0x72, 0x12, 0x19, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1a, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x2b, 0x92, 0x41,
	0x12, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x72, 0x2a, 0x0a, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x55,
	0x73, 0x65, 0x72, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x10, 0x2a, 0x0e, 0x2f, 0x76, 0x31, 0x2f, 0x75,
	0x73, 0x65, 0x72, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x1a, 0x5b, 0x92, 0x41, 0x58, 0x12, 0x56,
	0x49, 0x74, 0x20, 0x69, 0x73, 0x20, 0x74, 0x68, 0x65, 0x20, 0x73, 0x61, 0x6d, 0x65, 0x20, 0x61,
	0x73, 0x20, 0x55, 0x73, 0x65, 0x72, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x20, 0x62, 0x75,
	0x74, 0x20, 0x73, 0x68, 0x6f, 0x75, 0x6c, 0x64, 0x20, 0x62, 0x65, 0x20, 0x75, 0x73, 0x65, 0x64,
	0x20, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x6c, 0x79, 0x20, 0x61, 0x6e, 0x64, 0x20,
	0x6e, 0x6f, 0x74, 0x20, 0x65, 0x78, 0x70, 0x6f, 0x73, 0x65, 0x64, 0x20, 0x74, 0x6f, 0x20, 0x70,
	0x75, 0x62, 0x6c, 0x69, 0x63, 0x2e, 0x32, 0xdb, 0x01, 0x0a, 0x10, 0x55, 0x73, 0x65, 0x72, 0x51,
	0x75, 0x65, 0x72, 0x79, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x6f, 0x0a, 0x0b, 0x47,
	0x65, 0x74, 0x41, 0x6c, 0x6c, 0x55, 0x73, 0x65, 0x72, 0x73, 0x12, 0x1a, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x55, 0x73, 0x65, 0x72, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e,
	0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x55, 0x73, 0x65, 0x72, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x27, 0x92, 0x41, 0x13, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x72, 0x2a, 0x0b,
	0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x55, 0x73, 0x65, 0x72, 0x73, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x0b, 0x12, 0x09, 0x2f, 0x76, 0x31, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x1a, 0x56, 0x92, 0x41,
	0x53, 0x12, 0x51, 0x54, 0x68, 0x69, 0x73, 0x20, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x20,
	0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x73, 0x20, 0x62, 0x61, 0x73, 0x69, 0x63, 0x20, 0x71,
	0x75, 0x65, 0x72, 0x79, 0x20, 0x6f, 0x72, 0x20, 0x64, 0x61, 0x74, 0x61, 0x2d, 0x72, 0x65, 0x74,
	0x72, 0x69, 0x65, 0x76, 0x69, 0x6e, 0x67, 0x20, 0x75, 0x73, 0x65, 0x20, 0x63, 0x61, 0x73, 0x65,
	0x73, 0x20, 0x74, 0x6f, 0x20, 0x77, 0x6f, 0x72, 0x6b, 0x20, 0x77, 0x69, 0x74, 0x68, 0x20, 0x75,
	0x73, 0x65, 0x72, 0x2e, 0x42, 0x85, 0x02, 0x5a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x6e, 0x64, 0x72, 0x61, 0x73, 0x61, 0x70, 0x75, 0x74, 0x72, 0x61,
	0x2f, 0x61, 0x72, 0x6a, 0x75, 0x6e, 0x61, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x76, 0x31, 0x3b, 0x61, 0x70, 0x69, 0x76, 0x31, 0x92, 0x41, 0xcf, 0x01, 0x12, 0x95, 0x01,
	0x0a, 0x08, 0x55, 0x73, 0x65, 0x72, 0x20, 0x41, 0x50, 0x49, 0x22, 0x30, 0x0a, 0x0d, 0x49, 0x6e,
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
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x6a, 0x73, 0x6f, 0x6e, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_v1_user_proto_rawDescOnce sync.Once
	file_api_v1_user_proto_rawDescData = file_api_v1_user_proto_rawDesc
)

func file_api_v1_user_proto_rawDescGZIP() []byte {
	file_api_v1_user_proto_rawDescOnce.Do(func() {
		file_api_v1_user_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_v1_user_proto_rawDescData)
	})
	return file_api_v1_user_proto_rawDescData
}

var file_api_v1_user_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_api_v1_user_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_api_v1_user_proto_goTypes = []interface{}{
	(UserErrorCode)(0),            // 0: api.v1.UserErrorCode
	(*RegisterUserRequest)(nil),   // 1: api.v1.RegisterUserRequest
	(*RegisterUserResponse)(nil),  // 2: api.v1.RegisterUserResponse
	(*DeleteUserRequest)(nil),     // 3: api.v1.DeleteUserRequest
	(*DeleteUserResponse)(nil),    // 4: api.v1.DeleteUserResponse
	(*GetAllUsersRequest)(nil),    // 5: api.v1.GetAllUsersRequest
	(*GetAllUsersResponse)(nil),   // 6: api.v1.GetAllUsersResponse
	(*User)(nil),                  // 7: api.v1.User
	(*UserError)(nil),             // 8: api.v1.UserError
	(*timestamppb.Timestamp)(nil), // 9: google.protobuf.Timestamp
}
var file_api_v1_user_proto_depIdxs = []int32{
	7, // 0: api.v1.RegisterUserRequest.user:type_name -> api.v1.User
	7, // 1: api.v1.RegisterUserResponse.data:type_name -> api.v1.User
	7, // 2: api.v1.GetAllUsersResponse.data:type_name -> api.v1.User
	9, // 3: api.v1.User.created_at:type_name -> google.protobuf.Timestamp
	9, // 4: api.v1.User.updated_at:type_name -> google.protobuf.Timestamp
	0, // 5: api.v1.UserError.error_code:type_name -> api.v1.UserErrorCode
	1, // 6: api.v1.UserCommandService.RegisterUser:input_type -> api.v1.RegisterUserRequest
	3, // 7: api.v1.UserCommandInternalService.DeleteUser:input_type -> api.v1.DeleteUserRequest
	5, // 8: api.v1.UserQueryService.GetAllUsers:input_type -> api.v1.GetAllUsersRequest
	2, // 9: api.v1.UserCommandService.RegisterUser:output_type -> api.v1.RegisterUserResponse
	4, // 10: api.v1.UserCommandInternalService.DeleteUser:output_type -> api.v1.DeleteUserResponse
	6, // 11: api.v1.UserQueryService.GetAllUsers:output_type -> api.v1.GetAllUsersResponse
	9, // [9:12] is the sub-list for method output_type
	6, // [6:9] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_api_v1_user_proto_init() }
func file_api_v1_user_proto_init() {
	if File_api_v1_user_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_v1_user_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterUserRequest); i {
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
		file_api_v1_user_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterUserResponse); i {
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
		file_api_v1_user_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteUserRequest); i {
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
		file_api_v1_user_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteUserResponse); i {
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
		file_api_v1_user_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAllUsersRequest); i {
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
		file_api_v1_user_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAllUsersResponse); i {
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
		file_api_v1_user_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*User); i {
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
		file_api_v1_user_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserError); i {
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
			RawDescriptor: file_api_v1_user_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   3,
		},
		GoTypes:           file_api_v1_user_proto_goTypes,
		DependencyIndexes: file_api_v1_user_proto_depIdxs,
		EnumInfos:         file_api_v1_user_proto_enumTypes,
		MessageInfos:      file_api_v1_user_proto_msgTypes,
	}.Build()
	File_api_v1_user_proto = out.File
	file_api_v1_user_proto_rawDesc = nil
	file_api_v1_user_proto_goTypes = nil
	file_api_v1_user_proto_depIdxs = nil
}
