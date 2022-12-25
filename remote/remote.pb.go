// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: remote/remote.proto

package remote

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Node string `protobuf:"bytes,1,opt,name=node,proto3" json:"node,omitempty"`
	Id   string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *PID) Reset() {
	*x = PID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_remote_remote_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PID) ProtoMessage() {}

func (x *PID) ProtoReflect() protoreflect.Message {
	mi := &file_remote_remote_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PID.ProtoReflect.Descriptor instead.
func (*PID) Descriptor() ([]byte, []int) {
	return file_remote_remote_proto_rawDescGZIP(), []int{0}
}

func (x *PID) GetNode() string {
	if x != nil {
		return x.Node
	}
	return ""
}

func (x *PID) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type Error struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code   int32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg    string `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Detail string `protobuf:"bytes,3,opt,name=detail,proto3" json:"detail,omitempty"`
}

func (x *Error) Reset() {
	*x = Error{}
	if protoimpl.UnsafeEnabled {
		mi := &file_remote_remote_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Error) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Error) ProtoMessage() {}

func (x *Error) ProtoReflect() protoreflect.Message {
	mi := &file_remote_remote_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Error.ProtoReflect.Descriptor instead.
func (*Error) Descriptor() ([]byte, []int) {
	return file_remote_remote_proto_rawDescGZIP(), []int{1}
}

func (x *Error) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *Error) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *Error) GetDetail() string {
	if x != nil {
		return x.Detail
	}
	return ""
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	From      *PID       `protobuf:"bytes,1,opt,name=from,proto3" json:"from,omitempty"`
	RequestId int64      `protobuf:"varint,2,opt,name=requestId,proto3" json:"requestId,omitempty"`
	Data      *anypb.Any `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	Flag      int32      `protobuf:"varint,4,opt,name=flag,proto3" json:"flag,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_remote_remote_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_remote_remote_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_remote_remote_proto_rawDescGZIP(), []int{2}
}

func (x *Message) GetFrom() *PID {
	if x != nil {
		return x.From
	}
	return nil
}

func (x *Message) GetRequestId() int64 {
	if x != nil {
		return x.RequestId
	}
	return 0
}

func (x *Message) GetData() *anypb.Any {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *Message) GetFlag() int32 {
	if x != nil {
		return x.Flag
	}
	return 0
}

type OnMessageRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	To      *PID     `protobuf:"bytes,1,opt,name=to,proto3" json:"to,omitempty"`
	Message *Message `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *OnMessageRequest) Reset() {
	*x = OnMessageRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_remote_remote_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OnMessageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OnMessageRequest) ProtoMessage() {}

func (x *OnMessageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_remote_remote_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OnMessageRequest.ProtoReflect.Descriptor instead.
func (*OnMessageRequest) Descriptor() ([]byte, []int) {
	return file_remote_remote_proto_rawDescGZIP(), []int{3}
}

func (x *OnMessageRequest) GetTo() *PID {
	if x != nil {
		return x.To
	}
	return nil
}

func (x *OnMessageRequest) GetMessage() *Message {
	if x != nil {
		return x.Message
	}
	return nil
}

type OnMessageReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *OnMessageReply) Reset() {
	*x = OnMessageReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_remote_remote_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OnMessageReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OnMessageReply) ProtoMessage() {}

func (x *OnMessageReply) ProtoReflect() protoreflect.Message {
	mi := &file_remote_remote_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OnMessageReply.ProtoReflect.Descriptor instead.
func (*OnMessageReply) Descriptor() ([]byte, []int) {
	return file_remote_remote_proto_rawDescGZIP(), []int{4}
}

var File_remote_remote_proto protoreflect.FileDescriptor

var file_remote_remote_proto_rawDesc = []byte{
	0x0a, 0x13, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x2f, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x1a, 0x19, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61,
	0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x29, 0x0a, 0x03, 0x50, 0x49, 0x44, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x6f, 0x64, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x22, 0x45, 0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x12, 0x0a, 0x04,
	0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65,
	0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d,
	0x73, 0x67, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x22, 0x86, 0x01, 0x0a, 0x07, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1f, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x2e, 0x50, 0x49,
	0x44, 0x52, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x12, 0x1c, 0x0a, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x72, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x49, 0x64, 0x12, 0x28, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12,
	0x12, 0x0a, 0x04, 0x66, 0x6c, 0x61, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x66,
	0x6c, 0x61, 0x67, 0x22, 0x5a, 0x0a, 0x10, 0x4f, 0x6e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x02, 0x74, 0x6f, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x2e, 0x50, 0x49, 0x44,
	0x52, 0x02, 0x74, 0x6f, 0x12, 0x29, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x2e, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22,
	0x10, 0x0a, 0x0e, 0x4f, 0x6e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x70, 0x6c,
	0x79, 0x32, 0x49, 0x0a, 0x06, 0x52, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x12, 0x3f, 0x0a, 0x09, 0x4f,
	0x6e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x18, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74,
	0x65, 0x2e, 0x4f, 0x6e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x16, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x2e, 0x4f, 0x6e, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x42, 0x11, 0x5a, 0x0f,
	0x67, 0x6f, 0x2d, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x2f, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_remote_remote_proto_rawDescOnce sync.Once
	file_remote_remote_proto_rawDescData = file_remote_remote_proto_rawDesc
)

func file_remote_remote_proto_rawDescGZIP() []byte {
	file_remote_remote_proto_rawDescOnce.Do(func() {
		file_remote_remote_proto_rawDescData = protoimpl.X.CompressGZIP(file_remote_remote_proto_rawDescData)
	})
	return file_remote_remote_proto_rawDescData
}

var file_remote_remote_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_remote_remote_proto_goTypes = []interface{}{
	(*PID)(nil),              // 0: remote.PID
	(*Error)(nil),            // 1: remote.Error
	(*Message)(nil),          // 2: remote.Message
	(*OnMessageRequest)(nil), // 3: remote.OnMessageRequest
	(*OnMessageReply)(nil),   // 4: remote.OnMessageReply
	(*anypb.Any)(nil),        // 5: google.protobuf.Any
}
var file_remote_remote_proto_depIdxs = []int32{
	0, // 0: remote.Message.from:type_name -> remote.PID
	5, // 1: remote.Message.data:type_name -> google.protobuf.Any
	0, // 2: remote.OnMessageRequest.to:type_name -> remote.PID
	2, // 3: remote.OnMessageRequest.message:type_name -> remote.Message
	3, // 4: remote.Remote.OnMessage:input_type -> remote.OnMessageRequest
	4, // 5: remote.Remote.OnMessage:output_type -> remote.OnMessageReply
	5, // [5:6] is the sub-list for method output_type
	4, // [4:5] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_remote_remote_proto_init() }
func file_remote_remote_proto_init() {
	if File_remote_remote_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_remote_remote_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PID); i {
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
		file_remote_remote_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Error); i {
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
		file_remote_remote_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
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
		file_remote_remote_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OnMessageRequest); i {
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
		file_remote_remote_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OnMessageReply); i {
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
			RawDescriptor: file_remote_remote_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_remote_remote_proto_goTypes,
		DependencyIndexes: file_remote_remote_proto_depIdxs,
		MessageInfos:      file_remote_remote_proto_msgTypes,
	}.Build()
	File_remote_remote_proto = out.File
	file_remote_remote_proto_rawDesc = nil
	file_remote_remote_proto_goTypes = nil
	file_remote_remote_proto_depIdxs = nil
}