// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        (unknown)
// source: ctrlplane/common/v1/enums.proto

package commonv1

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

type AuthProvider int32

const (
	AuthProvider_AUTH_PROVIDER_UNSPECIFIED AuthProvider = 0
	AuthProvider_AUTH_PROVIDER_GOOGLE      AuthProvider = 1
	AuthProvider_AUTH_PROVIDER_GITHUB      AuthProvider = 2
)

// Enum value maps for AuthProvider.
var (
	AuthProvider_name = map[int32]string{
		0: "AUTH_PROVIDER_UNSPECIFIED",
		1: "AUTH_PROVIDER_GOOGLE",
		2: "AUTH_PROVIDER_GITHUB",
	}
	AuthProvider_value = map[string]int32{
		"AUTH_PROVIDER_UNSPECIFIED": 0,
		"AUTH_PROVIDER_GOOGLE":      1,
		"AUTH_PROVIDER_GITHUB":      2,
	}
)

func (x AuthProvider) Enum() *AuthProvider {
	p := new(AuthProvider)
	*p = x
	return p
}

func (x AuthProvider) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AuthProvider) Descriptor() protoreflect.EnumDescriptor {
	return file_ctrlplane_common_v1_enums_proto_enumTypes[0].Descriptor()
}

func (AuthProvider) Type() protoreflect.EnumType {
	return &file_ctrlplane_common_v1_enums_proto_enumTypes[0]
}

func (x AuthProvider) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AuthProvider.Descriptor instead.
func (AuthProvider) EnumDescriptor() ([]byte, []int) {
	return file_ctrlplane_common_v1_enums_proto_rawDescGZIP(), []int{0}
}

var File_ctrlplane_common_v1_enums_proto protoreflect.FileDescriptor

var file_ctrlplane_common_v1_enums_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2f, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x6e, 0x75, 0x6d, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x13, 0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2a, 0x61, 0x0a, 0x0c, 0x41, 0x75, 0x74, 0x68, 0x50, 0x72,
	0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x12, 0x1d, 0x0a, 0x19, 0x41, 0x55, 0x54, 0x48, 0x5f, 0x50,
	0x52, 0x4f, 0x56, 0x49, 0x44, 0x45, 0x52, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46,
	0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x18, 0x0a, 0x14, 0x41, 0x55, 0x54, 0x48, 0x5f, 0x50, 0x52,
	0x4f, 0x56, 0x49, 0x44, 0x45, 0x52, 0x5f, 0x47, 0x4f, 0x4f, 0x47, 0x4c, 0x45, 0x10, 0x01, 0x12,
	0x18, 0x0a, 0x14, 0x41, 0x55, 0x54, 0x48, 0x5f, 0x50, 0x52, 0x4f, 0x56, 0x49, 0x44, 0x45, 0x52,
	0x5f, 0x47, 0x49, 0x54, 0x48, 0x55, 0x42, 0x10, 0x02, 0x42, 0xd8, 0x01, 0x0a, 0x17, 0x63, 0x6f,
	0x6d, 0x2e, 0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x42, 0x0a, 0x45, 0x6e, 0x75, 0x6d, 0x73, 0x50, 0x72, 0x6f, 0x74,
	0x6f, 0x50, 0x01, 0x5a, 0x43, 0x67, 0x6f, 0x2e, 0x62, 0x72, 0x65, 0x75, 0x2e, 0x69, 0x6f, 0x2f,
	0x71, 0x75, 0x61, 0x6e, 0x74, 0x6d, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f,
	0x6e, 0x6f, 0x6d, 0x61, 0x64, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x74, 0x72, 0x6c,
	0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x3b,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x43, 0x43, 0x58, 0xaa, 0x02,
	0x13, 0x43, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x13, 0x43, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65,
	0x5c, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x1f, 0x43, 0x74, 0x72,
	0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x5c, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5c, 0x56, 0x31,
	0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x15, 0x43,
	0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x3a, 0x3a, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_ctrlplane_common_v1_enums_proto_rawDescOnce sync.Once
	file_ctrlplane_common_v1_enums_proto_rawDescData = file_ctrlplane_common_v1_enums_proto_rawDesc
)

func file_ctrlplane_common_v1_enums_proto_rawDescGZIP() []byte {
	file_ctrlplane_common_v1_enums_proto_rawDescOnce.Do(func() {
		file_ctrlplane_common_v1_enums_proto_rawDescData = protoimpl.X.CompressGZIP(file_ctrlplane_common_v1_enums_proto_rawDescData)
	})
	return file_ctrlplane_common_v1_enums_proto_rawDescData
}

var file_ctrlplane_common_v1_enums_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_ctrlplane_common_v1_enums_proto_goTypes = []any{
	(AuthProvider)(0), // 0: ctrlplane.common.v1.AuthProvider
}
var file_ctrlplane_common_v1_enums_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_ctrlplane_common_v1_enums_proto_init() }
func file_ctrlplane_common_v1_enums_proto_init() {
	if File_ctrlplane_common_v1_enums_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_ctrlplane_common_v1_enums_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_ctrlplane_common_v1_enums_proto_goTypes,
		DependencyIndexes: file_ctrlplane_common_v1_enums_proto_depIdxs,
		EnumInfos:         file_ctrlplane_common_v1_enums_proto_enumTypes,
	}.Build()
	File_ctrlplane_common_v1_enums_proto = out.File
	file_ctrlplane_common_v1_enums_proto_rawDesc = nil
	file_ctrlplane_common_v1_enums_proto_goTypes = nil
	file_ctrlplane_common_v1_enums_proto_depIdxs = nil
}
