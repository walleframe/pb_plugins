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
// source: mysql.proto

package mysql

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

var file_mysql_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.FileOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         20000,
		Name:          "mysql.db_name",
		Tag:           "bytes,20000,opt,name=db_name",
		Filename:      "mysql.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FileOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         20001,
		Name:          "mysql.db_charset",
		Tag:           "bytes,20001,opt,name=db_charset",
		Filename:      "mysql.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FileOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         20002,
		Name:          "mysql.db_collate",
		Tag:           "bytes,20002,opt,name=db_collate",
		Filename:      "mysql.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         20100,
		Name:          "mysql.tbl_name",
		Tag:           "bytes,20100,opt,name=tbl_name",
		Filename:      "mysql.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         20101,
		Name:          "mysql.ignore",
		Tag:           "varint,20101,opt,name=ignore",
		Filename:      "mysql.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         20102,
		Name:          "mysql.engine",
		Tag:           "bytes,20102,opt,name=engine",
		Filename:      "mysql.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         20103,
		Name:          "mysql.pks",
		Tag:           "bytes,20103,opt,name=pks",
		Filename:      "mysql.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         20104,
		Name:          "mysql.unique",
		Tag:           "bytes,20104,opt,name=unique",
		Filename:      "mysql.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         20105,
		Name:          "mysql.index",
		Tag:           "bytes,20105,opt,name=index",
		Filename:      "mysql.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         20106,
		Name:          "mysql.update",
		Tag:           "varint,20106,opt,name=update",
		Filename:      "mysql.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         20107,
		Name:          "mysql.upsert",
		Tag:           "varint,20107,opt,name=upsert",
		Filename:      "mysql.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         20109,
		Name:          "mysql.gen_ex",
		Tag:           "varint,20109,opt,name=gen_ex",
		Filename:      "mysql.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         20110,
		Name:          "mysql.tbl_opt",
		Tag:           "bytes,20110,opt,name=tbl_opt",
		Filename:      "mysql.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         20011,
		Name:          "mysql.tbl_charset",
		Tag:           "bytes,20011,opt,name=tbl_charset",
		Filename:      "mysql.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         20012,
		Name:          "mysql.tbl_collate",
		Tag:           "bytes,20012,opt,name=tbl_collate",
		Filename:      "mysql.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FieldOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         20200,
		Name:          "mysql.pk",
		Tag:           "varint,20200,opt,name=pk",
		Filename:      "mysql.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FieldOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         20201,
		Name:          "mysql.increment",
		Tag:           "varint,20201,opt,name=increment",
		Filename:      "mysql.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FieldOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         20202,
		Name:          "mysql.type",
		Tag:           "bytes,20202,opt,name=type",
		Filename:      "mysql.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FieldOptions)(nil),
		ExtensionType: (*int32)(nil),
		Field:         20203,
		Name:          "mysql.size",
		Tag:           "varint,20203,opt,name=size",
		Filename:      "mysql.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FieldOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         20204,
		Name:          "mysql.custom",
		Tag:           "varint,20204,opt,name=custom",
		Filename:      "mysql.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FieldOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         20205,
		Name:          "mysql.column",
		Tag:           "bytes,20205,opt,name=column",
		Filename:      "mysql.proto",
	},
}

// Extension fields to descriptorpb.FileOptions.
var (
	// optional string db_name = 20000;
	E_DbName = &file_mysql_proto_extTypes[0]
	// optional string db_charset = 20001;
	E_DbCharset = &file_mysql_proto_extTypes[1]
	// optional string db_collate = 20002;
	E_DbCollate = &file_mysql_proto_extTypes[2]
)

// Extension fields to descriptorpb.MessageOptions.
var (
	// optional string tbl_name = 20100;
	E_TblName = &file_mysql_proto_extTypes[3]
	// optional bool ignore = 20101;
	E_Ignore = &file_mysql_proto_extTypes[4]
	// optional string engine = 20102;
	E_Engine = &file_mysql_proto_extTypes[5]
	// optional string pks = 20103;
	E_Pks = &file_mysql_proto_extTypes[6]
	// optional string unique = 20104;
	E_Unique = &file_mysql_proto_extTypes[7]
	// optional string index = 20105;
	E_Index = &file_mysql_proto_extTypes[8]
	// optional bool update = 20106;
	E_Update = &file_mysql_proto_extTypes[9]
	// optional bool upsert = 20107;
	E_Upsert = &file_mysql_proto_extTypes[10]
	// optional bool gen_ex = 20109;
	E_GenEx = &file_mysql_proto_extTypes[11]
	// optional string tbl_opt = 20110;
	E_TblOpt = &file_mysql_proto_extTypes[12]
	// optional string tbl_charset = 20011;
	E_TblCharset = &file_mysql_proto_extTypes[13]
	// optional string tbl_collate = 20012;
	E_TblCollate = &file_mysql_proto_extTypes[14]
)

// Extension fields to descriptorpb.FieldOptions.
var (
	// optional bool pk = 20200;
	E_Pk = &file_mysql_proto_extTypes[15]
	// optional bool increment = 20201;
	E_Increment = &file_mysql_proto_extTypes[16]
	// optional string type = 20202;
	E_Type = &file_mysql_proto_extTypes[17]
	// optional int32 size = 20203;
	E_Size = &file_mysql_proto_extTypes[18]
	// optional bool custom = 20204;
	E_Custom = &file_mysql_proto_extTypes[19]
	// optional string column = 20205;
	E_Column = &file_mysql_proto_extTypes[20]
)

var File_mysql_proto protoreflect.FileDescriptor

var file_mysql_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x6d, 0x79, 0x73, 0x71, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x6d,
	0x79, 0x73, 0x71, 0x6c, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x3a, 0x3f, 0x0a, 0x07, 0x64, 0x62, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18,
	0xa0, 0x9c, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03, 0x88, 0x01, 0x02, 0x52, 0x06, 0x64, 0x62,
	0x4e, 0x61, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x3a, 0x45, 0x0a, 0x0a, 0x64, 0x62, 0x5f, 0x63, 0x68,
	0x61, 0x72, 0x73, 0x65, 0x74, 0x12, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x4f, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x18, 0xa1, 0x9c, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03, 0x88, 0x01, 0x02,
	0x52, 0x09, 0x64, 0x62, 0x43, 0x68, 0x61, 0x72, 0x73, 0x65, 0x74, 0x88, 0x01, 0x01, 0x3a, 0x45,
	0x0a, 0x0a, 0x64, 0x62, 0x5f, 0x63, 0x6f, 0x6c, 0x6c, 0x61, 0x74, 0x65, 0x12, 0x1c, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46,
	0x69, 0x6c, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xa2, 0x9c, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x03, 0x88, 0x01, 0x02, 0x52, 0x09, 0x64, 0x62, 0x43, 0x6f, 0x6c, 0x6c, 0x61,
	0x74, 0x65, 0x88, 0x01, 0x01, 0x3a, 0x44, 0x0a, 0x08, 0x74, 0x62, 0x6c, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x18, 0x84, 0x9d, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03, 0x88, 0x01, 0x02, 0x52,
	0x07, 0x74, 0x62, 0x6c, 0x4e, 0x61, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x3a, 0x41, 0x0a, 0x06, 0x69,
	0x67, 0x6e, 0x6f, 0x72, 0x65, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x85, 0x9d, 0x01, 0x20, 0x01, 0x28, 0x08, 0x42, 0x03,
	0x88, 0x01, 0x02, 0x52, 0x06, 0x69, 0x67, 0x6e, 0x6f, 0x72, 0x65, 0x88, 0x01, 0x01, 0x3a, 0x41,
	0x0a, 0x06, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x86, 0x9d, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x03, 0x88, 0x01, 0x02, 0x52, 0x06, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x88, 0x01,
	0x01, 0x3a, 0x3b, 0x0a, 0x03, 0x70, 0x6b, 0x73, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x87, 0x9d, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x03, 0x88, 0x01, 0x02, 0x52, 0x03, 0x70, 0x6b, 0x73, 0x88, 0x01, 0x01, 0x3a, 0x41,
	0x0a, 0x06, 0x75, 0x6e, 0x69, 0x71, 0x75, 0x65, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x88, 0x9d, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x03, 0x88, 0x01, 0x02, 0x52, 0x06, 0x75, 0x6e, 0x69, 0x71, 0x75, 0x65, 0x88, 0x01,
	0x01, 0x3a, 0x3f, 0x0a, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x89, 0x9d, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x03, 0x88, 0x01, 0x02, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x88,
	0x01, 0x01, 0x3a, 0x41, 0x0a, 0x06, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x1f, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x8a, 0x9d,
	0x01, 0x20, 0x01, 0x28, 0x08, 0x42, 0x03, 0x88, 0x01, 0x02, 0x52, 0x06, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x88, 0x01, 0x01, 0x3a, 0x41, 0x0a, 0x06, 0x75, 0x70, 0x73, 0x65, 0x72, 0x74, 0x12,
	0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x18, 0x8b, 0x9d, 0x01, 0x20, 0x01, 0x28, 0x08, 0x42, 0x03, 0x88, 0x01, 0x02, 0x52, 0x06, 0x75,
	0x70, 0x73, 0x65, 0x72, 0x74, 0x88, 0x01, 0x01, 0x3a, 0x40, 0x0a, 0x06, 0x67, 0x65, 0x6e, 0x5f,
	0x65, 0x78, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x18, 0x8d, 0x9d, 0x01, 0x20, 0x01, 0x28, 0x08, 0x42, 0x03, 0x88, 0x01, 0x02,
	0x52, 0x05, 0x67, 0x65, 0x6e, 0x45, 0x78, 0x88, 0x01, 0x01, 0x3a, 0x42, 0x0a, 0x07, 0x74, 0x62,
	0x6c, 0x5f, 0x6f, 0x70, 0x74, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x8e, 0x9d, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03,
	0x88, 0x01, 0x02, 0x52, 0x06, 0x74, 0x62, 0x6c, 0x4f, 0x70, 0x74, 0x88, 0x01, 0x01, 0x3a, 0x4a,
	0x0a, 0x0b, 0x74, 0x62, 0x6c, 0x5f, 0x63, 0x68, 0x61, 0x72, 0x73, 0x65, 0x74, 0x12, 0x1f, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xab,
	0x9c, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03, 0x88, 0x01, 0x02, 0x52, 0x0a, 0x74, 0x62, 0x6c,
	0x43, 0x68, 0x61, 0x72, 0x73, 0x65, 0x74, 0x88, 0x01, 0x01, 0x3a, 0x4a, 0x0a, 0x0b, 0x74, 0x62,
	0x6c, 0x5f, 0x63, 0x6f, 0x6c, 0x6c, 0x61, 0x74, 0x65, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xac, 0x9c, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x03, 0x88, 0x01, 0x02, 0x52, 0x0a, 0x74, 0x62, 0x6c, 0x43, 0x6f, 0x6c, 0x6c,
	0x61, 0x74, 0x65, 0x88, 0x01, 0x01, 0x3a, 0x37, 0x0a, 0x02, 0x70, 0x6b, 0x12, 0x1d, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46,
	0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xe8, 0x9d, 0x01, 0x20,
	0x01, 0x28, 0x08, 0x42, 0x03, 0x88, 0x01, 0x02, 0x52, 0x02, 0x70, 0x6b, 0x88, 0x01, 0x01, 0x3a,
	0x45, 0x0a, 0x09, 0x69, 0x6e, 0x63, 0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x1d, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46,
	0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xe9, 0x9d, 0x01, 0x20,
	0x01, 0x28, 0x08, 0x42, 0x03, 0x88, 0x01, 0x02, 0x52, 0x09, 0x69, 0x6e, 0x63, 0x72, 0x65, 0x6d,
	0x65, 0x6e, 0x74, 0x88, 0x01, 0x01, 0x3a, 0x3b, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x1d,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xea, 0x9d,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03, 0x88, 0x01, 0x02, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x88, 0x01, 0x01, 0x3a, 0x3b, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x1d, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69,
	0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xeb, 0x9d, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x42, 0x03, 0x88, 0x01, 0x02, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x88, 0x01, 0x01,
	0x3a, 0x3f, 0x0a, 0x06, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x12, 0x1d, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x65,
	0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xec, 0x9d, 0x01, 0x20, 0x01, 0x28,
	0x08, 0x42, 0x03, 0x88, 0x01, 0x02, 0x52, 0x06, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x88, 0x01,
	0x01, 0x3a, 0x3f, 0x0a, 0x06, 0x63, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x12, 0x1d, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69,
	0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xed, 0x9d, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x03, 0x88, 0x01, 0x02, 0x52, 0x06, 0x63, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x88,
	0x01, 0x01, 0x42, 0x2f, 0x5a, 0x2d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x66, 0x72, 0x61, 0x6d, 0x65, 0x2f, 0x70, 0x62, 0x5f, 0x70,
	0x6c, 0x75, 0x67, 0x69, 0x6e, 0x73, 0x2f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x64, 0x2f, 0x6d, 0x79,
	0x73, 0x71, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_mysql_proto_goTypes = []interface{}{
	(*descriptorpb.FileOptions)(nil),    // 0: google.protobuf.FileOptions
	(*descriptorpb.MessageOptions)(nil), // 1: google.protobuf.MessageOptions
	(*descriptorpb.FieldOptions)(nil),   // 2: google.protobuf.FieldOptions
}
var file_mysql_proto_depIdxs = []int32{
	0,  // 0: mysql.db_name:extendee -> google.protobuf.FileOptions
	0,  // 1: mysql.db_charset:extendee -> google.protobuf.FileOptions
	0,  // 2: mysql.db_collate:extendee -> google.protobuf.FileOptions
	1,  // 3: mysql.tbl_name:extendee -> google.protobuf.MessageOptions
	1,  // 4: mysql.ignore:extendee -> google.protobuf.MessageOptions
	1,  // 5: mysql.engine:extendee -> google.protobuf.MessageOptions
	1,  // 6: mysql.pks:extendee -> google.protobuf.MessageOptions
	1,  // 7: mysql.unique:extendee -> google.protobuf.MessageOptions
	1,  // 8: mysql.index:extendee -> google.protobuf.MessageOptions
	1,  // 9: mysql.update:extendee -> google.protobuf.MessageOptions
	1,  // 10: mysql.upsert:extendee -> google.protobuf.MessageOptions
	1,  // 11: mysql.gen_ex:extendee -> google.protobuf.MessageOptions
	1,  // 12: mysql.tbl_opt:extendee -> google.protobuf.MessageOptions
	1,  // 13: mysql.tbl_charset:extendee -> google.protobuf.MessageOptions
	1,  // 14: mysql.tbl_collate:extendee -> google.protobuf.MessageOptions
	2,  // 15: mysql.pk:extendee -> google.protobuf.FieldOptions
	2,  // 16: mysql.increment:extendee -> google.protobuf.FieldOptions
	2,  // 17: mysql.type:extendee -> google.protobuf.FieldOptions
	2,  // 18: mysql.size:extendee -> google.protobuf.FieldOptions
	2,  // 19: mysql.custom:extendee -> google.protobuf.FieldOptions
	2,  // 20: mysql.column:extendee -> google.protobuf.FieldOptions
	21, // [21:21] is the sub-list for method output_type
	21, // [21:21] is the sub-list for method input_type
	21, // [21:21] is the sub-list for extension type_name
	0,  // [0:21] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_mysql_proto_init() }
func file_mysql_proto_init() {
	if File_mysql_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_mysql_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 21,
			NumServices:   0,
		},
		GoTypes:           file_mysql_proto_goTypes,
		DependencyIndexes: file_mysql_proto_depIdxs,
		ExtensionInfos:    file_mysql_proto_extTypes,
	}.Build()
	File_mysql_proto = out.File
	file_mysql_proto_rawDesc = nil
	file_mysql_proto_goTypes = nil
	file_mysql_proto_depIdxs = nil
}
