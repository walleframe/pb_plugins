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
// 	protoc-gen-go v1.31.0
// 	protoc        v5.29.3
// source: gin.proto

package gin

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

type ApiPath struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Method string `protobuf:"bytes,1,opt,name=method,proto3" json:"method,omitempty"`
	Path   string `protobuf:"bytes,2,opt,name=path,proto3" json:"path,omitempty"`
}

func (x *ApiPath) Reset() {
	*x = ApiPath{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gin_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ApiPath) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ApiPath) ProtoMessage() {}

func (x *ApiPath) ProtoReflect() protoreflect.Message {
	mi := &file_gin_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ApiPath.ProtoReflect.Descriptor instead.
func (*ApiPath) Descriptor() ([]byte, []int) {
	return file_gin_proto_rawDescGZIP(), []int{0}
}

func (x *ApiPath) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

func (x *ApiPath) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

var file_gin_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.FileOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         23000,
		Name:          "gin.group",
		Tag:           "bytes,23000,opt,name=group",
		Filename:      "gin.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MethodOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         23200,
		Name:          "gin.get",
		Tag:           "bytes,23200,opt,name=get",
		Filename:      "gin.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MethodOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         23201,
		Name:          "gin.put",
		Tag:           "bytes,23201,opt,name=put",
		Filename:      "gin.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MethodOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         23202,
		Name:          "gin.post",
		Tag:           "bytes,23202,opt,name=post",
		Filename:      "gin.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MethodOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         23203,
		Name:          "gin.delete",
		Tag:           "bytes,23203,opt,name=delete",
		Filename:      "gin.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MethodOptions)(nil),
		ExtensionType: ([]*ApiPath)(nil),
		Field:         23204,
		Name:          "gin.custom",
		Tag:           "bytes,23204,rep,name=custom",
		Filename:      "gin.proto",
	},
}

// Extension fields to descriptorpb.FileOptions.
var (
	// optional string group = 23000;
	E_Group = &file_gin_proto_extTypes[0]
)

// Extension fields to descriptorpb.MethodOptions.
var (
	// optional string get = 23200;
	E_Get = &file_gin_proto_extTypes[1]
	// optional string put = 23201;
	E_Put = &file_gin_proto_extTypes[2]
	// optional string post = 23202;
	E_Post = &file_gin_proto_extTypes[3]
	// optional string delete = 23203;
	E_Delete = &file_gin_proto_extTypes[4]
	// repeated gin.ApiPath custom = 23204;
	E_Custom = &file_gin_proto_extTypes[5]
)

var File_gin_proto protoreflect.FileDescriptor

var file_gin_proto_rawDesc = []byte{
	0x0a, 0x09, 0x67, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x67, 0x69, 0x6e,
	0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x3f, 0x0a, 0x07, 0x41, 0x70, 0x69, 0x50, 0x61, 0x74, 0x68, 0x12, 0x1b, 0x0a,
	0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03, 0x88,
	0x01, 0x02, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x17, 0x0a, 0x04, 0x70, 0x61,
	0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03, 0x88, 0x01, 0x02, 0x52, 0x04, 0x70,
	0x61, 0x74, 0x68, 0x3a, 0x39, 0x0a, 0x05, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x1c, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46,
	0x69, 0x6c, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xd8, 0xb3, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x03, 0x88, 0x01, 0x02, 0x52, 0x05, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x3a, 0x37,
	0x0a, 0x03, 0x67, 0x65, 0x74, 0x12, 0x1e, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x4f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xa0, 0xb5, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03, 0x88,
	0x01, 0x02, 0x52, 0x03, 0x67, 0x65, 0x74, 0x3a, 0x37, 0x0a, 0x03, 0x70, 0x75, 0x74, 0x12, 0x1e,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xa1,
	0xb5, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03, 0x88, 0x01, 0x02, 0x52, 0x03, 0x70, 0x75, 0x74,
	0x3a, 0x39, 0x0a, 0x04, 0x70, 0x6f, 0x73, 0x74, 0x12, 0x1e, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x74, 0x68, 0x6f,
	0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xa2, 0xb5, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x03, 0x88, 0x01, 0x02, 0x52, 0x04, 0x70, 0x6f, 0x73, 0x74, 0x3a, 0x3d, 0x0a, 0x06, 0x64,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x1e, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x4f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xa3, 0xb5, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03, 0x88,
	0x01, 0x02, 0x52, 0x06, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x3a, 0x4b, 0x0a, 0x06, 0x63, 0x75,
	0x73, 0x74, 0x6f, 0x6d, 0x12, 0x1e, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x4f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x18, 0xa4, 0xb5, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x67,
	0x69, 0x6e, 0x2e, 0x41, 0x70, 0x69, 0x50, 0x61, 0x74, 0x68, 0x42, 0x03, 0x88, 0x01, 0x02, 0x52,
	0x06, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x42, 0x2d, 0x5a, 0x2b, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x66, 0x72, 0x61, 0x6d, 0x65,
	0x2f, 0x70, 0x62, 0x5f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x73, 0x2f, 0x65, 0x78, 0x74, 0x65,
	0x6e, 0x64, 0x2f, 0x67, 0x69, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_gin_proto_rawDescOnce sync.Once
	file_gin_proto_rawDescData = file_gin_proto_rawDesc
)

func file_gin_proto_rawDescGZIP() []byte {
	file_gin_proto_rawDescOnce.Do(func() {
		file_gin_proto_rawDescData = protoimpl.X.CompressGZIP(file_gin_proto_rawDescData)
	})
	return file_gin_proto_rawDescData
}

var file_gin_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_gin_proto_goTypes = []interface{}{
	(*ApiPath)(nil),                    // 0: gin.ApiPath
	(*descriptorpb.FileOptions)(nil),   // 1: google.protobuf.FileOptions
	(*descriptorpb.MethodOptions)(nil), // 2: google.protobuf.MethodOptions
}
var file_gin_proto_depIdxs = []int32{
	1, // 0: gin.group:extendee -> google.protobuf.FileOptions
	2, // 1: gin.get:extendee -> google.protobuf.MethodOptions
	2, // 2: gin.put:extendee -> google.protobuf.MethodOptions
	2, // 3: gin.post:extendee -> google.protobuf.MethodOptions
	2, // 4: gin.delete:extendee -> google.protobuf.MethodOptions
	2, // 5: gin.custom:extendee -> google.protobuf.MethodOptions
	0, // 6: gin.custom:type_name -> gin.ApiPath
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	6, // [6:7] is the sub-list for extension type_name
	0, // [0:6] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_gin_proto_init() }
func file_gin_proto_init() {
	if File_gin_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_gin_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ApiPath); i {
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
			RawDescriptor: file_gin_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 6,
			NumServices:   0,
		},
		GoTypes:           file_gin_proto_goTypes,
		DependencyIndexes: file_gin_proto_depIdxs,
		MessageInfos:      file_gin_proto_msgTypes,
		ExtensionInfos:    file_gin_proto_extTypes,
	}.Build()
	File_gin_proto = out.File
	file_gin_proto_rawDesc = nil
	file_gin_proto_goTypes = nil
	file_gin_proto_depIdxs = nil
}
