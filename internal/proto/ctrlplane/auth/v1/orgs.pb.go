// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.0
// 	protoc        (unknown)
// source: ctrlplane/auth/v1/orgs.proto

package authv1

import (
	v1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// OrgHooks contains the status of the hooks for an organization.
type OrgHooks struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Repo          v1.RepoHook            `protobuf:"varint,1,opt,name=repo,proto3,enum=ctrlplane.events.v1.RepoHook" json:"repo,omitempty"`
	Chat          v1.ChatHook            `protobuf:"varint,2,opt,name=chat,proto3,enum=ctrlplane.events.v1.ChatHook" json:"chat,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *OrgHooks) Reset() {
	*x = OrgHooks{}
	mi := &file_ctrlplane_auth_v1_orgs_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *OrgHooks) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OrgHooks) ProtoMessage() {}

func (x *OrgHooks) ProtoReflect() protoreflect.Message {
	mi := &file_ctrlplane_auth_v1_orgs_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OrgHooks.ProtoReflect.Descriptor instead.
func (*OrgHooks) Descriptor() ([]byte, []int) {
	return file_ctrlplane_auth_v1_orgs_proto_rawDescGZIP(), []int{0}
}

func (x *OrgHooks) GetRepo() v1.RepoHook {
	if x != nil {
		return x.Repo
	}
	return v1.RepoHook(0)
}

func (x *OrgHooks) GetChat() v1.ChatHook {
	if x != nil {
		return x.Chat
	}
	return v1.ChatHook(0)
}

// Represents an organization within the control plane.
type Org struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	CreatedAt     *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt     *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	Name          string                 `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	Domain        string                 `protobuf:"bytes,5,opt,name=domain,proto3" json:"domain,omitempty"`
	Slug          string                 `protobuf:"bytes,6,opt,name=slug,proto3" json:"slug,omitempty"`
	Hooks         *OrgHooks              `protobuf:"bytes,7,opt,name=hooks,proto3" json:"hooks,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Org) Reset() {
	*x = Org{}
	mi := &file_ctrlplane_auth_v1_orgs_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Org) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Org) ProtoMessage() {}

func (x *Org) ProtoReflect() protoreflect.Message {
	mi := &file_ctrlplane_auth_v1_orgs_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Org.ProtoReflect.Descriptor instead.
func (*Org) Descriptor() ([]byte, []int) {
	return file_ctrlplane_auth_v1_orgs_proto_rawDescGZIP(), []int{1}
}

func (x *Org) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Org) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Org) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *Org) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Org) GetDomain() string {
	if x != nil {
		return x.Domain
	}
	return ""
}

func (x *Org) GetSlug() string {
	if x != nil {
		return x.Slug
	}
	return ""
}

func (x *Org) GetHooks() *OrgHooks {
	if x != nil {
		return x.Hooks
	}
	return nil
}

// CreateOrgRequest is the request to create a new organization.
type CreateOrgRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Domain        string                 `protobuf:"bytes,2,opt,name=domain,proto3" json:"domain,omitempty"`
	Slug          string                 `protobuf:"bytes,3,opt,name=slug,proto3" json:"slug,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateOrgRequest) Reset() {
	*x = CreateOrgRequest{}
	mi := &file_ctrlplane_auth_v1_orgs_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateOrgRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateOrgRequest) ProtoMessage() {}

func (x *CreateOrgRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ctrlplane_auth_v1_orgs_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateOrgRequest.ProtoReflect.Descriptor instead.
func (*CreateOrgRequest) Descriptor() ([]byte, []int) {
	return file_ctrlplane_auth_v1_orgs_proto_rawDescGZIP(), []int{2}
}

func (x *CreateOrgRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateOrgRequest) GetDomain() string {
	if x != nil {
		return x.Domain
	}
	return ""
}

func (x *CreateOrgRequest) GetSlug() string {
	if x != nil {
		return x.Slug
	}
	return ""
}

// CreateOrgResponse is the response to creating a new organization.
type CreateOrgResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Org           *Org                   `protobuf:"bytes,1,opt,name=org,proto3" json:"org,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateOrgResponse) Reset() {
	*x = CreateOrgResponse{}
	mi := &file_ctrlplane_auth_v1_orgs_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateOrgResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateOrgResponse) ProtoMessage() {}

func (x *CreateOrgResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ctrlplane_auth_v1_orgs_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateOrgResponse.ProtoReflect.Descriptor instead.
func (*CreateOrgResponse) Descriptor() ([]byte, []int) {
	return file_ctrlplane_auth_v1_orgs_proto_rawDescGZIP(), []int{3}
}

func (x *CreateOrgResponse) GetOrg() *Org {
	if x != nil {
		return x.Org
	}
	return nil
}

// GetOrgByIDRequest is the request to retrieve an organization by its unique
// identifier.
type GetOrgByIDRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetOrgByIDRequest) Reset() {
	*x = GetOrgByIDRequest{}
	mi := &file_ctrlplane_auth_v1_orgs_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetOrgByIDRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetOrgByIDRequest) ProtoMessage() {}

func (x *GetOrgByIDRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ctrlplane_auth_v1_orgs_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetOrgByIDRequest.ProtoReflect.Descriptor instead.
func (*GetOrgByIDRequest) Descriptor() ([]byte, []int) {
	return file_ctrlplane_auth_v1_orgs_proto_rawDescGZIP(), []int{4}
}

func (x *GetOrgByIDRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

// GetOrgByIDResponse is the response to retrieving an organization by its unique
// identifier.
type GetOrgByIDResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Org           *Org                   `protobuf:"bytes,1,opt,name=org,proto3" json:"org,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetOrgByIDResponse) Reset() {
	*x = GetOrgByIDResponse{}
	mi := &file_ctrlplane_auth_v1_orgs_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetOrgByIDResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetOrgByIDResponse) ProtoMessage() {}

func (x *GetOrgByIDResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ctrlplane_auth_v1_orgs_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetOrgByIDResponse.ProtoReflect.Descriptor instead.
func (*GetOrgByIDResponse) Descriptor() ([]byte, []int) {
	return file_ctrlplane_auth_v1_orgs_proto_rawDescGZIP(), []int{5}
}

func (x *GetOrgByIDResponse) GetOrg() *Org {
	if x != nil {
		return x.Org
	}
	return nil
}

// SetOrgHooksRequest is the request to set the hooks for an organization.
type SetOrgHooksRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	OrgId         string                 `protobuf:"bytes,1,opt,name=org_id,json=orgId,proto3" json:"org_id,omitempty"`
	Hooks         *OrgHooks              `protobuf:"bytes,2,opt,name=hooks,proto3" json:"hooks,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SetOrgHooksRequest) Reset() {
	*x = SetOrgHooksRequest{}
	mi := &file_ctrlplane_auth_v1_orgs_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SetOrgHooksRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetOrgHooksRequest) ProtoMessage() {}

func (x *SetOrgHooksRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ctrlplane_auth_v1_orgs_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetOrgHooksRequest.ProtoReflect.Descriptor instead.
func (*SetOrgHooksRequest) Descriptor() ([]byte, []int) {
	return file_ctrlplane_auth_v1_orgs_proto_rawDescGZIP(), []int{6}
}

func (x *SetOrgHooksRequest) GetOrgId() string {
	if x != nil {
		return x.OrgId
	}
	return ""
}

func (x *SetOrgHooksRequest) GetHooks() *OrgHooks {
	if x != nil {
		return x.Hooks
	}
	return nil
}

var File_ctrlplane_auth_v1_orgs_proto protoreflect.FileDescriptor

var file_ctrlplane_auth_v1_orgs_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2f, 0x61, 0x75, 0x74, 0x68,
	0x2f, 0x76, 0x31, 0x2f, 0x6f, 0x72, 0x67, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x11,
	0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x76,
	0x31, 0x1a, 0x1f, 0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2f, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x68, 0x6f, 0x6f, 0x6b, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x70, 0x0a, 0x08, 0x4f, 0x72, 0x67, 0x48, 0x6f, 0x6f, 0x6b, 0x73, 0x12, 0x31, 0x0a, 0x04,
	0x72, 0x65, 0x70, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1d, 0x2e, 0x63, 0x74, 0x72,
	0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x76, 0x31,
	0x2e, 0x52, 0x65, 0x70, 0x6f, 0x48, 0x6f, 0x6f, 0x6b, 0x52, 0x04, 0x72, 0x65, 0x70, 0x6f, 0x12,
	0x31, 0x0a, 0x04, 0x63, 0x68, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1d, 0x2e,
	0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73,
	0x2e, 0x76, 0x31, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x48, 0x6f, 0x6f, 0x6b, 0x52, 0x04, 0x63, 0x68,
	0x61, 0x74, 0x22, 0xfe, 0x01, 0x0a, 0x03, 0x4f, 0x72, 0x67, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64,
	0x5f, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x12, 0x12, 0x0a, 0x04,
	0x73, 0x6c, 0x75, 0x67, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x73, 0x6c, 0x75, 0x67,
	0x12, 0x31, 0x0a, 0x05, 0x68, 0x6f, 0x6f, 0x6b, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1b, 0x2e, 0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x61, 0x75, 0x74, 0x68,
	0x2e, 0x76, 0x31, 0x2e, 0x4f, 0x72, 0x67, 0x48, 0x6f, 0x6f, 0x6b, 0x73, 0x52, 0x05, 0x68, 0x6f,
	0x6f, 0x6b, 0x73, 0x22, 0x52, 0x0a, 0x10, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4f, 0x72, 0x67,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x64,
	0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x6f, 0x6d,
	0x61, 0x69, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x6c, 0x75, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x73, 0x6c, 0x75, 0x67, 0x22, 0x3d, 0x0a, 0x11, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x4f, 0x72, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x28, 0x0a, 0x03,
	0x6f, 0x72, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x63, 0x74, 0x72, 0x6c,
	0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x4f, 0x72,
	0x67, 0x52, 0x03, 0x6f, 0x72, 0x67, 0x22, 0x23, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x4f, 0x72, 0x67,
	0x42, 0x79, 0x49, 0x44, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x3e, 0x0a, 0x12, 0x47,
	0x65, 0x74, 0x4f, 0x72, 0x67, 0x42, 0x79, 0x49, 0x44, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x28, 0x0a, 0x03, 0x6f, 0x72, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16,
	0x2e, 0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e,
	0x76, 0x31, 0x2e, 0x4f, 0x72, 0x67, 0x52, 0x03, 0x6f, 0x72, 0x67, 0x22, 0x5e, 0x0a, 0x12, 0x53,
	0x65, 0x74, 0x4f, 0x72, 0x67, 0x48, 0x6f, 0x6f, 0x6b, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x15, 0x0a, 0x06, 0x6f, 0x72, 0x67, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x6f, 0x72, 0x67, 0x49, 0x64, 0x12, 0x31, 0x0a, 0x05, 0x68, 0x6f, 0x6f, 0x6b,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c,
	0x61, 0x6e, 0x65, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x4f, 0x72, 0x67, 0x48,
	0x6f, 0x6f, 0x6b, 0x73, 0x52, 0x05, 0x68, 0x6f, 0x6f, 0x6b, 0x73, 0x32, 0x8d, 0x02, 0x0a, 0x0a,
	0x4f, 0x72, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x56, 0x0a, 0x09, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x4f, 0x72, 0x67, 0x12, 0x23, 0x2e, 0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c,
	0x61, 0x6e, 0x65, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x4f, 0x72, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x63,
	0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x76, 0x31,
	0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4f, 0x72, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x59, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x4f, 0x72, 0x67, 0x42, 0x79, 0x49, 0x44,
	0x12, 0x24, 0x2e, 0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x61, 0x75, 0x74,
	0x68, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x4f, 0x72, 0x67, 0x42, 0x79, 0x49, 0x44, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61,
	0x6e, 0x65, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x4f, 0x72,
	0x67, 0x42, 0x79, 0x49, 0x44, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4c, 0x0a,
	0x0b, 0x53, 0x65, 0x74, 0x4f, 0x72, 0x67, 0x48, 0x6f, 0x6f, 0x6b, 0x73, 0x12, 0x25, 0x2e, 0x63,
	0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x76, 0x31,
	0x2e, 0x53, 0x65, 0x74, 0x4f, 0x72, 0x67, 0x48, 0x6f, 0x6f, 0x6b, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0xc3, 0x01, 0x0a, 0x15,
	0x63, 0x6f, 0x6d, 0x2e, 0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x61, 0x75,
	0x74, 0x68, 0x2e, 0x76, 0x31, 0x42, 0x09, 0x4f, 0x72, 0x67, 0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x50, 0x01, 0x5a, 0x39, 0x67, 0x6f, 0x2e, 0x62, 0x72, 0x65, 0x75, 0x2e, 0x69, 0x6f, 0x2f, 0x71,
	0x75, 0x61, 0x6e, 0x74, 0x6d, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2f, 0x61,
	0x75, 0x74, 0x68, 0x2f, 0x76, 0x31, 0x3b, 0x61, 0x75, 0x74, 0x68, 0x76, 0x31, 0xa2, 0x02, 0x03,
	0x43, 0x41, 0x58, 0xaa, 0x02, 0x11, 0x43, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e,
	0x41, 0x75, 0x74, 0x68, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x11, 0x43, 0x74, 0x72, 0x6c, 0x70, 0x6c,
	0x61, 0x6e, 0x65, 0x5c, 0x41, 0x75, 0x74, 0x68, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x1d, 0x43, 0x74,
	0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x5c, 0x41, 0x75, 0x74, 0x68, 0x5c, 0x56, 0x31, 0x5c,
	0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x13, 0x43, 0x74,
	0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x3a, 0x3a, 0x41, 0x75, 0x74, 0x68, 0x3a, 0x3a, 0x56,
	0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_ctrlplane_auth_v1_orgs_proto_rawDescOnce sync.Once
	file_ctrlplane_auth_v1_orgs_proto_rawDescData = file_ctrlplane_auth_v1_orgs_proto_rawDesc
)

func file_ctrlplane_auth_v1_orgs_proto_rawDescGZIP() []byte {
	file_ctrlplane_auth_v1_orgs_proto_rawDescOnce.Do(func() {
		file_ctrlplane_auth_v1_orgs_proto_rawDescData = protoimpl.X.CompressGZIP(file_ctrlplane_auth_v1_orgs_proto_rawDescData)
	})
	return file_ctrlplane_auth_v1_orgs_proto_rawDescData
}

var file_ctrlplane_auth_v1_orgs_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_ctrlplane_auth_v1_orgs_proto_goTypes = []any{
	(*OrgHooks)(nil),              // 0: ctrlplane.auth.v1.OrgHooks
	(*Org)(nil),                   // 1: ctrlplane.auth.v1.Org
	(*CreateOrgRequest)(nil),      // 2: ctrlplane.auth.v1.CreateOrgRequest
	(*CreateOrgResponse)(nil),     // 3: ctrlplane.auth.v1.CreateOrgResponse
	(*GetOrgByIDRequest)(nil),     // 4: ctrlplane.auth.v1.GetOrgByIDRequest
	(*GetOrgByIDResponse)(nil),    // 5: ctrlplane.auth.v1.GetOrgByIDResponse
	(*SetOrgHooksRequest)(nil),    // 6: ctrlplane.auth.v1.SetOrgHooksRequest
	(v1.RepoHook)(0),              // 7: ctrlplane.events.v1.RepoHook
	(v1.ChatHook)(0),              // 8: ctrlplane.events.v1.ChatHook
	(*timestamppb.Timestamp)(nil), // 9: google.protobuf.Timestamp
	(*emptypb.Empty)(nil),         // 10: google.protobuf.Empty
}
var file_ctrlplane_auth_v1_orgs_proto_depIdxs = []int32{
	7,  // 0: ctrlplane.auth.v1.OrgHooks.repo:type_name -> ctrlplane.events.v1.RepoHook
	8,  // 1: ctrlplane.auth.v1.OrgHooks.chat:type_name -> ctrlplane.events.v1.ChatHook
	9,  // 2: ctrlplane.auth.v1.Org.created_at:type_name -> google.protobuf.Timestamp
	9,  // 3: ctrlplane.auth.v1.Org.updated_at:type_name -> google.protobuf.Timestamp
	0,  // 4: ctrlplane.auth.v1.Org.hooks:type_name -> ctrlplane.auth.v1.OrgHooks
	1,  // 5: ctrlplane.auth.v1.CreateOrgResponse.org:type_name -> ctrlplane.auth.v1.Org
	1,  // 6: ctrlplane.auth.v1.GetOrgByIDResponse.org:type_name -> ctrlplane.auth.v1.Org
	0,  // 7: ctrlplane.auth.v1.SetOrgHooksRequest.hooks:type_name -> ctrlplane.auth.v1.OrgHooks
	2,  // 8: ctrlplane.auth.v1.OrgService.CreateOrg:input_type -> ctrlplane.auth.v1.CreateOrgRequest
	4,  // 9: ctrlplane.auth.v1.OrgService.GetOrgByID:input_type -> ctrlplane.auth.v1.GetOrgByIDRequest
	6,  // 10: ctrlplane.auth.v1.OrgService.SetOrgHooks:input_type -> ctrlplane.auth.v1.SetOrgHooksRequest
	3,  // 11: ctrlplane.auth.v1.OrgService.CreateOrg:output_type -> ctrlplane.auth.v1.CreateOrgResponse
	5,  // 12: ctrlplane.auth.v1.OrgService.GetOrgByID:output_type -> ctrlplane.auth.v1.GetOrgByIDResponse
	10, // 13: ctrlplane.auth.v1.OrgService.SetOrgHooks:output_type -> google.protobuf.Empty
	11, // [11:14] is the sub-list for method output_type
	8,  // [8:11] is the sub-list for method input_type
	8,  // [8:8] is the sub-list for extension type_name
	8,  // [8:8] is the sub-list for extension extendee
	0,  // [0:8] is the sub-list for field type_name
}

func init() { file_ctrlplane_auth_v1_orgs_proto_init() }
func file_ctrlplane_auth_v1_orgs_proto_init() {
	if File_ctrlplane_auth_v1_orgs_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_ctrlplane_auth_v1_orgs_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_ctrlplane_auth_v1_orgs_proto_goTypes,
		DependencyIndexes: file_ctrlplane_auth_v1_orgs_proto_depIdxs,
		MessageInfos:      file_ctrlplane_auth_v1_orgs_proto_msgTypes,
	}.Build()
	File_ctrlplane_auth_v1_orgs_proto = out.File
	file_ctrlplane_auth_v1_orgs_proto_rawDesc = nil
	file_ctrlplane_auth_v1_orgs_proto_goTypes = nil
	file_ctrlplane_auth_v1_orgs_proto_depIdxs = nil
}
