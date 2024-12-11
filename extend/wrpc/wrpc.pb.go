// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v5.28.3
// source: wrpc.proto

package wrpc

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var file_wrpc_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.FileOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         22000,
		Name:          "wrpc.stub_client",
		Tag:           "varint,22000,opt,name=stub_client",
		Filename:      "wrpc.proto",
	},
	{
		ExtendedType:  (*descriptorpb.ServiceOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         22100,
		Name:          "wrpc.stub_service",
		Tag:           "varint,22100,opt,name=stub_service",
		Filename:      "wrpc.proto",
	},
}

// Extension fields to descriptorpb.FileOptions.
var (
	// optional bool stub_client = 22000;
	E_StubClient = &file_wrpc_proto_extTypes[0] // Generate client stubs.
)

// Extension fields to descriptorpb.ServiceOptions.
var (
	// optional bool stub_service = 22100;
	E_StubService = &file_wrpc_proto_extTypes[1] // Generate client stubs.
)

var File_wrpc_proto protoreflect.FileDescriptor

var file_wrpc_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x77, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x77, 0x72,
	0x70, 0x63, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x3a, 0x3f, 0x0a, 0x0b, 0x73, 0x74, 0x75, 0x62, 0x5f, 0x63, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x12, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x18, 0xf0, 0xab, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x73, 0x74, 0x75, 0x62, 0x43,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x3a, 0x44, 0x0a, 0x0c, 0x73, 0x74, 0x75, 0x62, 0x5f, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xd4, 0xac, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b,
	0x73, 0x74, 0x75, 0x62, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x42, 0x2e, 0x5a, 0x2c, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x66,
	0x72, 0x61, 0x6d, 0x65, 0x2f, 0x70, 0x62, 0x5f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x73, 0x2f,
	0x65, 0x78, 0x74, 0x65, 0x6e, 0x64, 0x2f, 0x77, 0x72, 0x70, 0x63,
}

var file_wrpc_proto_goTypes = []interface{}{
	(*descriptorpb.FileOptions)(nil),    // 0: google.protobuf.FileOptions
	(*descriptorpb.ServiceOptions)(nil), // 1: google.protobuf.ServiceOptions
}
var file_wrpc_proto_depIdxs = []int32{
	0, // 0: wrpc.stub_client:extendee -> google.protobuf.FileOptions
	1, // 1: wrpc.stub_service:extendee -> google.protobuf.ServiceOptions
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	0, // [0:2] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_wrpc_proto_init() }
func file_wrpc_proto_init() {
	if File_wrpc_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_wrpc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 2,
			NumServices:   0,
		},
		GoTypes:           file_wrpc_proto_goTypes,
		DependencyIndexes: file_wrpc_proto_depIdxs,
		ExtensionInfos:    file_wrpc_proto_extTypes,
	}.Build()
	File_wrpc_proto = out.File
	file_wrpc_proto_rawDesc = nil
	file_wrpc_proto_goTypes = nil
	file_wrpc_proto_depIdxs = nil
}
