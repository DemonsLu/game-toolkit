// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: message.proto

package common

import (
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoimpl"
	"reflect"
	"sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ID int32

const (
	ID__         ID = 0
	ID_FooReqID  ID = 1
	ID_FooRespID ID = 2
)

// Enum value maps for ID.
var (
	ID_name = map[int32]string{
		0: "_",
		1: "FooReqID",
		2: "FooRespID",
	}
	ID_value = map[string]int32{
		"_":         0,
		"FooReqID":  1,
		"FooRespID": 2,
	}
)

func (x ID) Enum() *ID {
	p := new(ID)
	*p = x
	return p
}

func (x ID) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ID) Descriptor() protoreflect.EnumDescriptor {
	return file_message_proto_enumTypes[0].Descriptor()
}

func (ID) Type() protoreflect.EnumType {
	return &file_message_proto_enumTypes[0]
}

func (x ID) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ID.Descriptor instead.
func (ID) EnumDescriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{0}
}

type FooReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Bar string `protobuf:"bytes,1,opt,name=Bar,proto3" json:"Bar,omitempty"`
	Buz int32  `protobuf:"varint,2,opt,name=Buz,proto3" json:"Buz,omitempty"`
}

func (x *FooReq) Reset() {
	*x = FooReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FooReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FooReq) ProtoMessage() {}

func (x *FooReq) ProtoReflect() protoreflect.Message {
	mi := &file_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FooReq.ProtoReflect.Descriptor instead.
func (*FooReq) Descriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{0}
}

func (x *FooReq) GetBar() string {
	if x != nil {
		return x.Bar
	}
	return ""
}

func (x *FooReq) GetBuz() int32 {
	if x != nil {
		return x.Buz
	}
	return 0
}

type FooResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    int32  `protobuf:"varint,1,opt,name=Code,proto3" json:"Code,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=Message,proto3" json:"Message,omitempty"`
}

func (x *FooResp) Reset() {
	*x = FooResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FooResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FooResp) ProtoMessage() {}

func (x *FooResp) ProtoReflect() protoreflect.Message {
	mi := &file_message_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FooResp.ProtoReflect.Descriptor instead.
func (*FooResp) Descriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{1}
}

func (x *FooResp) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *FooResp) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_message_proto protoreflect.FileDescriptor

var file_message_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x06, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x22, 0x2c, 0x0a, 0x06, 0x46, 0x6f, 0x6f, 0x52, 0x65,
	0x71, 0x12, 0x10, 0x0a, 0x03, 0x42, 0x61, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x42, 0x61, 0x72, 0x12, 0x10, 0x0a, 0x03, 0x42, 0x75, 0x7a, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x03, 0x42, 0x75, 0x7a, 0x22, 0x37, 0x0a, 0x07, 0x46, 0x6f, 0x6f, 0x52, 0x65, 0x73, 0x70,
	0x12, 0x12, 0x0a, 0x04, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04,
	0x43, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2a, 0x28,
	0x0a, 0x02, 0x49, 0x44, 0x12, 0x05, 0x0a, 0x01, 0x5f, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x46,
	0x6f, 0x6f, 0x52, 0x65, 0x71, 0x49, 0x44, 0x10, 0x01, 0x12, 0x0d, 0x0a, 0x09, 0x46, 0x6f, 0x6f,
	0x52, 0x65, 0x73, 0x70, 0x49, 0x44, 0x10, 0x02, 0x42, 0x09, 0x5a, 0x07, 0x2f, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_message_proto_rawDescOnce sync.Once
	file_message_proto_rawDescData = file_message_proto_rawDesc
)

func file_message_proto_rawDescGZIP() []byte {
	file_message_proto_rawDescOnce.Do(func() {
		file_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_message_proto_rawDescData)
	})
	return file_message_proto_rawDescData
}

var file_message_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_message_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_message_proto_goTypes = []interface{}{
	(ID)(0),         // 0: common.ID
	(*FooReq)(nil),  // 1: common.FooReq
	(*FooResp)(nil), // 2: common.FooResp
}
var file_message_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_message_proto_init() }
func file_message_proto_init() {
	if File_message_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FooReq); i {
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
		file_message_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FooResp); i {
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
			RawDescriptor: file_message_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_message_proto_goTypes,
		DependencyIndexes: file_message_proto_depIdxs,
		EnumInfos:         file_message_proto_enumTypes,
		MessageInfos:      file_message_proto_msgTypes,
	}.Build()
	File_message_proto = out.File
	file_message_proto_rawDesc = nil
	file_message_proto_goTypes = nil
	file_message_proto_depIdxs = nil
}
