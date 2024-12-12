// Copyright 2024 aggronmagi <czy463@163.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v5.28.3
// source: redis.proto

package redis

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RedisScript struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Lua   *string `protobuf:"bytes,1,opt,name=lua,proto3,oneof" json:"lua,omitempty"`
	Argv  *string `protobuf:"bytes,2,opt,name=argv,proto3,oneof" json:"argv,omitempty"`
	Reply *string `protobuf:"bytes,3,opt,name=reply,proto3,oneof" json:"reply,omitempty"`
}

func (x *RedisScript) Reset() {
	*x = RedisScript{}
	if protoimpl.UnsafeEnabled {
		mi := &file_redis_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RedisScript) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RedisScript) ProtoMessage() {}

func (x *RedisScript) ProtoReflect() protoreflect.Message {
	mi := &file_redis_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RedisScript.ProtoReflect.Descriptor instead.
func (*RedisScript) Descriptor() ([]byte, []int) {
	return file_redis_proto_rawDescGZIP(), []int{0}
}

func (x *RedisScript) GetLua() string {
	if x != nil && x.Lua != nil {
		return *x.Lua
	}
	return ""
}

func (x *RedisScript) GetArgv() string {
	if x != nil && x.Argv != nil {
		return *x.Argv
	}
	return ""
}

func (x *RedisScript) GetReply() string {
	if x != nil && x.Reply != nil {
		return *x.Reply
	}
	return ""
}

var file_redis_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.FileOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         21000,
		Name:          "redis.pkg_svc",
		Tag:           "bytes,21000,opt,name=pkg_svc",
		Filename:      "redis.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FileOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         21001,
		Name:          "redis.pkg_util",
		Tag:           "bytes,21001,opt,name=pkg_util",
		Filename:      "redis.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FileOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         21002,
		Name:          "redis.pkg_pb",
		Tag:           "bytes,21002,opt,name=pkg_pb",
		Filename:      "redis.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FileOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         21003,
		Name:          "redis.pkg_wpb",
		Tag:           "bytes,21003,opt,name=pkg_wpb",
		Filename:      "redis.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         21100,
		Name:          "redis.key",
		Tag:           "bytes,21100,opt,name=key",
		Filename:      "redis.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         21101,
		Name:          "redis.type",
		Tag:           "bytes,21101,opt,name=type",
		Filename:      "redis.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*int32)(nil),
		Field:         21102,
		Name:          "redis.key_size",
		Tag:           "varint,21102,opt,name=key_size",
		Filename:      "redis.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         21103,
		Name:          "redis.protobuf",
		Tag:           "varint,21103,opt,name=protobuf",
		Filename:      "redis.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         21104,
		Name:          "redis.wpb",
		Tag:           "varint,21104,opt,name=wpb",
		Filename:      "redis.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         21105,
		Name:          "redis.field",
		Tag:           "bytes,21105,opt,name=field",
		Filename:      "redis.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         21106,
		Name:          "redis.value",
		Tag:           "bytes,21106,opt,name=value",
		Filename:      "redis.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         21107,
		Name:          "redis.member",
		Tag:           "bytes,21107,opt,name=member",
		Filename:      "redis.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: ([]*RedisScript)(nil),
		Field:         21108,
		Name:          "redis.script",
		Tag:           "bytes,21108,rep,name=script",
		Filename:      "redis.proto",
	},
}

// Extension fields to descriptorpb.FileOptions.
var (
	// optional string pkg_svc = 21000;
	E_PkgSvc = &file_redis_proto_extTypes[0] // default: github.com/walleframe/svc_redis
	// optional string pkg_util = 21001;
	E_PkgUtil = &file_redis_proto_extTypes[1] // default: github.com/walleframe/walle/util
	// optional string pkg_pb = 21002;
	E_PkgPb = &file_redis_proto_extTypes[2] // protobuf message pkg: github.com/gogo/protobuf/proto
	// optional string pkg_wpb = 21003;
	E_PkgWpb = &file_redis_proto_extTypes[3] // wpb message, pkg: github.com/walleframe/walle/process/message
)

// Extension fields to descriptorpb.MessageOptions.
var (
	// optional string key = 21100;
	E_Key = &file_redis_proto_extTypes[4] // must set
	// optional string type = 21101;
	E_Type = &file_redis_proto_extTypes[5] // default: only keys. '!' means not key operation. string,hash,set,zset,lock
	// optional int32 key_size = 21102;
	E_KeySize = &file_redis_proto_extTypes[6] // default: 64
	// optional bool protobuf = 21103;
	E_Protobuf = &file_redis_proto_extTypes[7] // default: false
	// optional bool wpb = 21104;
	E_Wpb = &file_redis_proto_extTypes[8] // default: false
	// optional string field = 21105;
	E_Field = &file_redis_proto_extTypes[9] // use merged field
	// optional string value = 21106;
	E_Value = &file_redis_proto_extTypes[10] // use merged value
	// optional string member = 21107;
	E_Member = &file_redis_proto_extTypes[11] // use merged member for zset
	// repeated redis.RedisScript script = 21108;
	E_Script = &file_redis_proto_extTypes[12] // redis script support
)

var File_redis_proto protoreflect.FileDescriptor

var file_redis_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x72, 0x65, 0x64, 0x69, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x72,
	0x65, 0x64, 0x69, 0x73, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x82, 0x01, 0x0a, 0x0b, 0x52, 0x65, 0x64, 0x69, 0x73,
	0x53, 0x63, 0x72, 0x69, 0x70, 0x74, 0x12, 0x1a, 0x0a, 0x03, 0x6c, 0x75, 0x61, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x03, 0x88, 0x01, 0x02, 0x48, 0x00, 0x52, 0x03, 0x6c, 0x75, 0x61, 0x88,
	0x01, 0x01, 0x12, 0x1c, 0x0a, 0x04, 0x61, 0x72, 0x67, 0x76, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x03, 0x88, 0x01, 0x02, 0x48, 0x01, 0x52, 0x04, 0x61, 0x72, 0x67, 0x76, 0x88, 0x01, 0x01,
	0x12, 0x1e, 0x0a, 0x05, 0x72, 0x65, 0x70, 0x6c, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x03, 0x88, 0x01, 0x02, 0x48, 0x02, 0x52, 0x05, 0x72, 0x65, 0x70, 0x6c, 0x79, 0x88, 0x01, 0x01,
	0x42, 0x06, 0x0a, 0x04, 0x5f, 0x6c, 0x75, 0x61, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x61, 0x72, 0x67,
	0x76, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x72, 0x65, 0x70, 0x6c, 0x79, 0x3a, 0x3f, 0x0a, 0x07, 0x70,
	0x6b, 0x67, 0x5f, 0x73, 0x76, 0x63, 0x12, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x4f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x18, 0x88, 0xa4, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03, 0x88, 0x01,
	0x02, 0x52, 0x06, 0x70, 0x6b, 0x67, 0x53, 0x76, 0x63, 0x88, 0x01, 0x01, 0x3a, 0x41, 0x0a, 0x08,
	0x70, 0x6b, 0x67, 0x5f, 0x75, 0x74, 0x69, 0x6c, 0x12, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x4f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x89, 0xa4, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03,
	0x88, 0x01, 0x02, 0x52, 0x07, 0x70, 0x6b, 0x67, 0x55, 0x74, 0x69, 0x6c, 0x88, 0x01, 0x01, 0x3a,
	0x3d, 0x0a, 0x06, 0x70, 0x6b, 0x67, 0x5f, 0x70, 0x62, 0x12, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x6c, 0x65,
	0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x8a, 0xa4, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x03, 0x88, 0x01, 0x02, 0x52, 0x05, 0x70, 0x6b, 0x67, 0x50, 0x62, 0x88, 0x01, 0x01, 0x3a, 0x3f,
	0x0a, 0x07, 0x70, 0x6b, 0x67, 0x5f, 0x77, 0x70, 0x62, 0x12, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x6c, 0x65,
	0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x8b, 0xa4, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x03, 0x88, 0x01, 0x02, 0x52, 0x06, 0x70, 0x6b, 0x67, 0x57, 0x70, 0x62, 0x88, 0x01, 0x01, 0x3a,
	0x3b, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xec, 0xa4, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x03, 0x88, 0x01, 0x02, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x88, 0x01, 0x01, 0x3a, 0x3d, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xed, 0xa4, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03, 0x88,
	0x01, 0x02, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x88, 0x01, 0x01, 0x3a, 0x44, 0x0a, 0x08, 0x6b,
	0x65, 0x79, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xee, 0xa4, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x42, 0x03, 0x88, 0x01, 0x02, 0x52, 0x07, 0x6b, 0x65, 0x79, 0x53, 0x69, 0x7a, 0x65, 0x88, 0x01,
	0x01, 0x3a, 0x45, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x12, 0x1f, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xef,
	0xa4, 0x01, 0x20, 0x01, 0x28, 0x08, 0x42, 0x03, 0x88, 0x01, 0x02, 0x52, 0x08, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x88, 0x01, 0x01, 0x3a, 0x3b, 0x0a, 0x03, 0x77, 0x70, 0x62, 0x12,
	0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x18, 0xf0, 0xa4, 0x01, 0x20, 0x01, 0x28, 0x08, 0x42, 0x03, 0x88, 0x01, 0x02, 0x52, 0x03, 0x77,
	0x70, 0x62, 0x88, 0x01, 0x01, 0x3a, 0x3f, 0x0a, 0x05, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x1f,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18,
	0xf1, 0xa4, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03, 0x88, 0x01, 0x02, 0x52, 0x05, 0x66, 0x69,
	0x65, 0x6c, 0x64, 0x88, 0x01, 0x01, 0x3a, 0x3f, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12,
	0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x18, 0xf2, 0xa4, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03, 0x88, 0x01, 0x02, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x88, 0x01, 0x01, 0x3a, 0x41, 0x0a, 0x06, 0x6d, 0x65, 0x6d, 0x62, 0x65,
	0x72, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x18, 0xf3, 0xa4, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03, 0x88, 0x01, 0x02, 0x52,
	0x06, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x88, 0x01, 0x01, 0x3a, 0x52, 0x0a, 0x06, 0x73, 0x63,
	0x72, 0x69, 0x70, 0x74, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xf4, 0xa4, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e,
	0x72, 0x65, 0x64, 0x69, 0x73, 0x2e, 0x52, 0x65, 0x64, 0x69, 0x73, 0x53, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x42, 0x03, 0x88, 0x01, 0x02, 0x52, 0x06, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x42, 0x2f,
	0x5a, 0x2d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x77, 0x61, 0x6c,
	0x6c, 0x65, 0x66, 0x72, 0x61, 0x6d, 0x65, 0x2f, 0x70, 0x62, 0x5f, 0x70, 0x6c, 0x75, 0x67, 0x69,
	0x6e, 0x73, 0x2f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x64, 0x2f, 0x72, 0x65, 0x64, 0x69, 0x73, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_redis_proto_rawDescOnce sync.Once
	file_redis_proto_rawDescData = file_redis_proto_rawDesc
)

func file_redis_proto_rawDescGZIP() []byte {
	file_redis_proto_rawDescOnce.Do(func() {
		file_redis_proto_rawDescData = protoimpl.X.CompressGZIP(file_redis_proto_rawDescData)
	})
	return file_redis_proto_rawDescData
}

var file_redis_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_redis_proto_goTypes = []interface{}{
	(*RedisScript)(nil),                 // 0: redis.RedisScript
	(*descriptorpb.FileOptions)(nil),    // 1: google.protobuf.FileOptions
	(*descriptorpb.MessageOptions)(nil), // 2: google.protobuf.MessageOptions
}
var file_redis_proto_depIdxs = []int32{
	1,  // 0: redis.pkg_svc:extendee -> google.protobuf.FileOptions
	1,  // 1: redis.pkg_util:extendee -> google.protobuf.FileOptions
	1,  // 2: redis.pkg_pb:extendee -> google.protobuf.FileOptions
	1,  // 3: redis.pkg_wpb:extendee -> google.protobuf.FileOptions
	2,  // 4: redis.key:extendee -> google.protobuf.MessageOptions
	2,  // 5: redis.type:extendee -> google.protobuf.MessageOptions
	2,  // 6: redis.key_size:extendee -> google.protobuf.MessageOptions
	2,  // 7: redis.protobuf:extendee -> google.protobuf.MessageOptions
	2,  // 8: redis.wpb:extendee -> google.protobuf.MessageOptions
	2,  // 9: redis.field:extendee -> google.protobuf.MessageOptions
	2,  // 10: redis.value:extendee -> google.protobuf.MessageOptions
	2,  // 11: redis.member:extendee -> google.protobuf.MessageOptions
	2,  // 12: redis.script:extendee -> google.protobuf.MessageOptions
	0,  // 13: redis.script:type_name -> redis.RedisScript
	14, // [14:14] is the sub-list for method output_type
	14, // [14:14] is the sub-list for method input_type
	13, // [13:14] is the sub-list for extension type_name
	0,  // [0:13] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_redis_proto_init() }
func file_redis_proto_init() {
	if File_redis_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_redis_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RedisScript); i {
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
	file_redis_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_redis_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 13,
			NumServices:   0,
		},
		GoTypes:           file_redis_proto_goTypes,
		DependencyIndexes: file_redis_proto_depIdxs,
		MessageInfos:      file_redis_proto_msgTypes,
		ExtensionInfos:    file_redis_proto_extTypes,
	}.Build()
	File_redis_proto = out.File
	file_redis_proto_rawDesc = nil
	file_redis_proto_goTypes = nil
	file_redis_proto_depIdxs = nil
}
