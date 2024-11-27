// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        (unknown)
// source: ctrlplane/core/v1/chat.proto

package corev1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ChannelKind int32

const (
	ChannelKind_CHANNEL_KIND_UNSPECIFIED ChannelKind = 0
	ChannelKind_CHANNEL_KIND_DEFAULT     ChannelKind = 1
	ChannelKind_CHANNEL_KIND_USER        ChannelKind = 2
)

// Enum value maps for ChannelKind.
var (
	ChannelKind_name = map[int32]string{
		0: "CHANNEL_KIND_UNSPECIFIED",
		1: "CHANNEL_KIND_DEFAULT",
		2: "CHANNEL_KIND_USER",
	}
	ChannelKind_value = map[string]int32{
		"CHANNEL_KIND_UNSPECIFIED": 0,
		"CHANNEL_KIND_DEFAULT":     1,
		"CHANNEL_KIND_USER":        2,
	}
)

func (x ChannelKind) Enum() *ChannelKind {
	p := new(ChannelKind)
	*p = x
	return p
}

func (x ChannelKind) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ChannelKind) Descriptor() protoreflect.EnumDescriptor {
	return file_ctrlplane_core_v1_chat_proto_enumTypes[0].Descriptor()
}

func (ChannelKind) Type() protoreflect.EnumType {
	return &file_ctrlplane_core_v1_chat_proto_enumTypes[0]
}

func (x ChannelKind) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ChannelKind.Descriptor instead.
func (ChannelKind) EnumDescriptor() ([]byte, []int) {
	return file_ctrlplane_core_v1_chat_proto_rawDescGZIP(), []int{0}
}

type Channel struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Kind        ChannelKind `protobuf:"varint,1,opt,name=kind,proto3,enum=ctrlplane.core.v1.ChannelKind" json:"kind,omitempty"`
	Credentials string      `protobuf:"bytes,2,opt,name=credentials,proto3" json:"credentials,omitempty"`
	Id          string      `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
	Name        string      `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *Channel) Reset() {
	*x = Channel{}
	mi := &file_ctrlplane_core_v1_chat_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Channel) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Channel) ProtoMessage() {}

func (x *Channel) ProtoReflect() protoreflect.Message {
	mi := &file_ctrlplane_core_v1_chat_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Channel.ProtoReflect.Descriptor instead.
func (*Channel) Descriptor() ([]byte, []int) {
	return file_ctrlplane_core_v1_chat_proto_rawDescGZIP(), []int{0}
}

func (x *Channel) GetKind() ChannelKind {
	if x != nil {
		return x.Kind
	}
	return ChannelKind_CHANNEL_KIND_UNSPECIFIED
}

func (x *Channel) GetCredentials() string {
	if x != nil {
		return x.Credentials
	}
	return ""
}

func (x *Channel) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Channel) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

var File_ctrlplane_core_v1_chat_proto protoreflect.FileDescriptor

var file_ctrlplane_core_v1_chat_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2f, 0x63, 0x6f, 0x72, 0x65,
	0x2f, 0x76, 0x31, 0x2f, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x11,
	0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76,
	0x31, 0x22, 0x83, 0x01, 0x0a, 0x07, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x12, 0x32, 0x0a,
	0x04, 0x6b, 0x69, 0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1e, 0x2e, 0x63, 0x74,
	0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e,
	0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x4b, 0x69, 0x6e, 0x64, 0x52, 0x04, 0x6b, 0x69, 0x6e,
	0x64, 0x12, 0x20, 0x0a, 0x0b, 0x63, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x73,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69,
	0x61, 0x6c, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x2a, 0x5c, 0x0a, 0x0b, 0x43, 0x68, 0x61, 0x6e, 0x6e,
	0x65, 0x6c, 0x4b, 0x69, 0x6e, 0x64, 0x12, 0x1c, 0x0a, 0x18, 0x43, 0x48, 0x41, 0x4e, 0x4e, 0x45,
	0x4c, 0x5f, 0x4b, 0x49, 0x4e, 0x44, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49,
	0x45, 0x44, 0x10, 0x00, 0x12, 0x18, 0x0a, 0x14, 0x43, 0x48, 0x41, 0x4e, 0x4e, 0x45, 0x4c, 0x5f,
	0x4b, 0x49, 0x4e, 0x44, 0x5f, 0x44, 0x45, 0x46, 0x41, 0x55, 0x4c, 0x54, 0x10, 0x01, 0x12, 0x15,
	0x0a, 0x11, 0x43, 0x48, 0x41, 0x4e, 0x4e, 0x45, 0x4c, 0x5f, 0x4b, 0x49, 0x4e, 0x44, 0x5f, 0x55,
	0x53, 0x45, 0x52, 0x10, 0x02, 0x42, 0xc3, 0x01, 0x0a, 0x15, 0x63, 0x6f, 0x6d, 0x2e, 0x63, 0x74,
	0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x42,
	0x09, 0x43, 0x68, 0x61, 0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x39, 0x67, 0x6f,
	0x2e, 0x62, 0x72, 0x65, 0x75, 0x2e, 0x69, 0x6f, 0x2f, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x6d, 0x2f,
	0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63,
	0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31,
	0x3b, 0x63, 0x6f, 0x72, 0x65, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x43, 0x43, 0x58, 0xaa, 0x02, 0x11,
	0x43, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x43, 0x6f, 0x72, 0x65, 0x2e, 0x56,
	0x31, 0xca, 0x02, 0x11, 0x43, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x5c, 0x43, 0x6f,
	0x72, 0x65, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x1d, 0x43, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e,
	0x65, 0x5c, 0x43, 0x6f, 0x72, 0x65, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x13, 0x43, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e,
	0x65, 0x3a, 0x3a, 0x43, 0x6f, 0x72, 0x65, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_ctrlplane_core_v1_chat_proto_rawDescOnce sync.Once
	file_ctrlplane_core_v1_chat_proto_rawDescData = file_ctrlplane_core_v1_chat_proto_rawDesc
)

func file_ctrlplane_core_v1_chat_proto_rawDescGZIP() []byte {
	file_ctrlplane_core_v1_chat_proto_rawDescOnce.Do(func() {
		file_ctrlplane_core_v1_chat_proto_rawDescData = protoimpl.X.CompressGZIP(file_ctrlplane_core_v1_chat_proto_rawDescData)
	})
	return file_ctrlplane_core_v1_chat_proto_rawDescData
}

var file_ctrlplane_core_v1_chat_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_ctrlplane_core_v1_chat_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_ctrlplane_core_v1_chat_proto_goTypes = []any{
	(ChannelKind)(0), // 0: ctrlplane.core.v1.ChannelKind
	(*Channel)(nil),  // 1: ctrlplane.core.v1.Channel
}
var file_ctrlplane_core_v1_chat_proto_depIdxs = []int32{
	0, // 0: ctrlplane.core.v1.Channel.kind:type_name -> ctrlplane.core.v1.ChannelKind
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_ctrlplane_core_v1_chat_proto_init() }
func file_ctrlplane_core_v1_chat_proto_init() {
	if File_ctrlplane_core_v1_chat_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_ctrlplane_core_v1_chat_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_ctrlplane_core_v1_chat_proto_goTypes,
		DependencyIndexes: file_ctrlplane_core_v1_chat_proto_depIdxs,
		EnumInfos:         file_ctrlplane_core_v1_chat_proto_enumTypes,
		MessageInfos:      file_ctrlplane_core_v1_chat_proto_msgTypes,
	}.Build()
	File_ctrlplane_core_v1_chat_proto = out.File
	file_ctrlplane_core_v1_chat_proto_rawDesc = nil
	file_ctrlplane_core_v1_chat_proto_goTypes = nil
	file_ctrlplane_core_v1_chat_proto_depIdxs = nil
}
