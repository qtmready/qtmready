// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        (unknown)
// source: ctrlplane/core/v1/repos.proto

package corev1

import (
	v1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
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

// Represents repo within the control plane.
type Repo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	CreatedAt     *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt     *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	OrgId         string                 `protobuf:"bytes,4,opt,name=org_id,json=orgId,proto3" json:"org_id,omitempty"`
	Name          string                 `protobuf:"bytes,5,opt,name=name,proto3" json:"name,omitempty"`
	Hook          v1.RepoHook            `protobuf:"varint,6,opt,name=hook,proto3,enum=ctrlplane.events.v1.RepoHook" json:"hook,omitempty"`
	HookId        string                 `protobuf:"bytes,7,opt,name=hook_id,json=hookId,proto3" json:"hook_id,omitempty"`
	DefaultBranch string                 `protobuf:"bytes,8,opt,name=default_branch,json=defaultBranch,proto3" json:"default_branch,omitempty"`
	IsMonorepo    bool                   `protobuf:"varint,9,opt,name=is_monorepo,json=isMonorepo,proto3" json:"is_monorepo,omitempty"`
	Threshold     int32                  `protobuf:"varint,10,opt,name=threshold,proto3" json:"threshold,omitempty"`
	StaleDuration *durationpb.Duration   `protobuf:"bytes,11,opt,name=stale_duration,json=staleDuration,proto3" json:"stale_duration,omitempty"`
	Url           string                 `protobuf:"bytes,12,opt,name=url,proto3" json:"url,omitempty"`
	IsActive      bool                   `protobuf:"varint,13,opt,name=is_active,json=isActive,proto3" json:"is_active,omitempty"`
}

func (x *Repo) Reset() {
	*x = Repo{}
	mi := &file_ctrlplane_core_v1_repos_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Repo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Repo) ProtoMessage() {}

func (x *Repo) ProtoReflect() protoreflect.Message {
	mi := &file_ctrlplane_core_v1_repos_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Repo.ProtoReflect.Descriptor instead.
func (*Repo) Descriptor() ([]byte, []int) {
	return file_ctrlplane_core_v1_repos_proto_rawDescGZIP(), []int{0}
}

func (x *Repo) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Repo) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Repo) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *Repo) GetOrgId() string {
	if x != nil {
		return x.OrgId
	}
	return ""
}

func (x *Repo) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Repo) GetHook() v1.RepoHook {
	if x != nil {
		return x.Hook
	}
	return v1.RepoHook(0)
}

func (x *Repo) GetHookId() string {
	if x != nil {
		return x.HookId
	}
	return ""
}

func (x *Repo) GetDefaultBranch() string {
	if x != nil {
		return x.DefaultBranch
	}
	return ""
}

func (x *Repo) GetIsMonorepo() bool {
	if x != nil {
		return x.IsMonorepo
	}
	return false
}

func (x *Repo) GetThreshold() int32 {
	if x != nil {
		return x.Threshold
	}
	return 0
}

func (x *Repo) GetStaleDuration() *durationpb.Duration {
	if x != nil {
		return x.StaleDuration
	}
	return nil
}

func (x *Repo) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *Repo) GetIsActive() bool {
	if x != nil {
		return x.IsActive
	}
	return false
}

// Request to create a org's core repo.
type CreateRepoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name          string               `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Hook          v1.RepoHook          `protobuf:"varint,2,opt,name=hook,proto3,enum=ctrlplane.events.v1.RepoHook" json:"hook,omitempty"`
	HookId        string               `protobuf:"bytes,3,opt,name=hook_id,json=hookId,proto3" json:"hook_id,omitempty"`
	DefaultBranch string               `protobuf:"bytes,4,opt,name=default_branch,json=defaultBranch,proto3" json:"default_branch,omitempty"`
	IsMonorepo    bool                 `protobuf:"varint,5,opt,name=is_monorepo,json=isMonorepo,proto3" json:"is_monorepo,omitempty"`
	Threshold     int32                `protobuf:"varint,6,opt,name=threshold,proto3" json:"threshold,omitempty"`
	StaleDuration *durationpb.Duration `protobuf:"bytes,7,opt,name=stale_duration,json=staleDuration,proto3" json:"stale_duration,omitempty"`
	OrgId         string               `protobuf:"bytes,8,opt,name=org_id,json=orgId,proto3" json:"org_id,omitempty"`
}

func (x *CreateRepoRequest) Reset() {
	*x = CreateRepoRequest{}
	mi := &file_ctrlplane_core_v1_repos_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateRepoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateRepoRequest) ProtoMessage() {}

func (x *CreateRepoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ctrlplane_core_v1_repos_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateRepoRequest.ProtoReflect.Descriptor instead.
func (*CreateRepoRequest) Descriptor() ([]byte, []int) {
	return file_ctrlplane_core_v1_repos_proto_rawDescGZIP(), []int{1}
}

func (x *CreateRepoRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateRepoRequest) GetHook() v1.RepoHook {
	if x != nil {
		return x.Hook
	}
	return v1.RepoHook(0)
}

func (x *CreateRepoRequest) GetHookId() string {
	if x != nil {
		return x.HookId
	}
	return ""
}

func (x *CreateRepoRequest) GetDefaultBranch() string {
	if x != nil {
		return x.DefaultBranch
	}
	return ""
}

func (x *CreateRepoRequest) GetIsMonorepo() bool {
	if x != nil {
		return x.IsMonorepo
	}
	return false
}

func (x *CreateRepoRequest) GetThreshold() int32 {
	if x != nil {
		return x.Threshold
	}
	return 0
}

func (x *CreateRepoRequest) GetStaleDuration() *durationpb.Duration {
	if x != nil {
		return x.StaleDuration
	}
	return nil
}

func (x *CreateRepoRequest) GetOrgId() string {
	if x != nil {
		return x.OrgId
	}
	return ""
}

// Response containing org's core repo.
type CreateRepoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Repo *Repo `protobuf:"bytes,1,opt,name=repo,proto3" json:"repo,omitempty"`
}

func (x *CreateRepoResponse) Reset() {
	*x = CreateRepoResponse{}
	mi := &file_ctrlplane_core_v1_repos_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateRepoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateRepoResponse) ProtoMessage() {}

func (x *CreateRepoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ctrlplane_core_v1_repos_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateRepoResponse.ProtoReflect.Descriptor instead.
func (*CreateRepoResponse) Descriptor() ([]byte, []int) {
	return file_ctrlplane_core_v1_repos_proto_rawDescGZIP(), []int{2}
}

func (x *CreateRepoResponse) GetRepo() *Repo {
	if x != nil {
		return x.Repo
	}
	return nil
}

// Request to get org's core repo by id.
type GetRepoByIDRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetRepoByIDRequest) Reset() {
	*x = GetRepoByIDRequest{}
	mi := &file_ctrlplane_core_v1_repos_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetRepoByIDRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRepoByIDRequest) ProtoMessage() {}

func (x *GetRepoByIDRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ctrlplane_core_v1_repos_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRepoByIDRequest.ProtoReflect.Descriptor instead.
func (*GetRepoByIDRequest) Descriptor() ([]byte, []int) {
	return file_ctrlplane_core_v1_repos_proto_rawDescGZIP(), []int{3}
}

func (x *GetRepoByIDRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

// Response to get org's core repo.
type GetRepoByIDResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Repo *Repo `protobuf:"bytes,1,opt,name=repo,proto3" json:"repo,omitempty"`
}

func (x *GetRepoByIDResponse) Reset() {
	*x = GetRepoByIDResponse{}
	mi := &file_ctrlplane_core_v1_repos_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetRepoByIDResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRepoByIDResponse) ProtoMessage() {}

func (x *GetRepoByIDResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ctrlplane_core_v1_repos_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRepoByIDResponse.ProtoReflect.Descriptor instead.
func (*GetRepoByIDResponse) Descriptor() ([]byte, []int) {
	return file_ctrlplane_core_v1_repos_proto_rawDescGZIP(), []int{4}
}

func (x *GetRepoByIDResponse) GetRepo() *Repo {
	if x != nil {
		return x.Repo
	}
	return nil
}

// Request to get org's core repo by org_id.
type GetOrgReposByOrgIDRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OrgId string `protobuf:"bytes,1,opt,name=org_id,json=orgId,proto3" json:"org_id,omitempty"`
}

func (x *GetOrgReposByOrgIDRequest) Reset() {
	*x = GetOrgReposByOrgIDRequest{}
	mi := &file_ctrlplane_core_v1_repos_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetOrgReposByOrgIDRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetOrgReposByOrgIDRequest) ProtoMessage() {}

func (x *GetOrgReposByOrgIDRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ctrlplane_core_v1_repos_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetOrgReposByOrgIDRequest.ProtoReflect.Descriptor instead.
func (*GetOrgReposByOrgIDRequest) Descriptor() ([]byte, []int) {
	return file_ctrlplane_core_v1_repos_proto_rawDescGZIP(), []int{5}
}

func (x *GetOrgReposByOrgIDRequest) GetOrgId() string {
	if x != nil {
		return x.OrgId
	}
	return ""
}

// Response to get org's core repo.
type GetOrgReposByOrgIDResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Repo *Repo `protobuf:"bytes,1,opt,name=repo,proto3" json:"repo,omitempty"`
}

func (x *GetOrgReposByOrgIDResponse) Reset() {
	*x = GetOrgReposByOrgIDResponse{}
	mi := &file_ctrlplane_core_v1_repos_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetOrgReposByOrgIDResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetOrgReposByOrgIDResponse) ProtoMessage() {}

func (x *GetOrgReposByOrgIDResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ctrlplane_core_v1_repos_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetOrgReposByOrgIDResponse.ProtoReflect.Descriptor instead.
func (*GetOrgReposByOrgIDResponse) Descriptor() ([]byte, []int) {
	return file_ctrlplane_core_v1_repos_proto_rawDescGZIP(), []int{6}
}

func (x *GetOrgReposByOrgIDResponse) GetRepo() *Repo {
	if x != nil {
		return x.Repo
	}
	return nil
}

type RepoExtended struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	CreatedAt     *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt     *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	OrgId         string                 `protobuf:"bytes,4,opt,name=org_id,json=orgId,proto3" json:"org_id,omitempty"`
	Name          string                 `protobuf:"bytes,5,opt,name=name,proto3" json:"name,omitempty"`
	Hook          v1.RepoHook            `protobuf:"varint,6,opt,name=hook,proto3,enum=ctrlplane.events.v1.RepoHook" json:"hook,omitempty"`
	HookId        string                 `protobuf:"bytes,7,opt,name=hook_id,json=hookId,proto3" json:"hook_id,omitempty"`
	DefaultBranch string                 `protobuf:"bytes,8,opt,name=default_branch,json=defaultBranch,proto3" json:"default_branch,omitempty"`
	IsMonorepo    bool                   `protobuf:"varint,9,opt,name=is_monorepo,json=isMonorepo,proto3" json:"is_monorepo,omitempty"`
	Threshold     int32                  `protobuf:"varint,10,opt,name=threshold,proto3" json:"threshold,omitempty"`
	StaleDuration *durationpb.Duration   `protobuf:"bytes,11,opt,name=stale_duration,json=staleDuration,proto3" json:"stale_duration,omitempty"`
	Url           string                 `protobuf:"bytes,12,opt,name=url,proto3" json:"url,omitempty"`
	IsActive      bool                   `protobuf:"varint,13,opt,name=is_active,json=isActive,proto3" json:"is_active,omitempty"`
	ChatEnabled   bool                   `protobuf:"varint,14,opt,name=chat_enabled,json=chatEnabled,proto3" json:"chat_enabled,omitempty"`
	ChannelName   string                 `protobuf:"bytes,15,opt,name=channel_name,json=channelName,proto3" json:"channel_name,omitempty"`
}

func (x *RepoExtended) Reset() {
	*x = RepoExtended{}
	mi := &file_ctrlplane_core_v1_repos_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RepoExtended) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RepoExtended) ProtoMessage() {}

func (x *RepoExtended) ProtoReflect() protoreflect.Message {
	mi := &file_ctrlplane_core_v1_repos_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RepoExtended.ProtoReflect.Descriptor instead.
func (*RepoExtended) Descriptor() ([]byte, []int) {
	return file_ctrlplane_core_v1_repos_proto_rawDescGZIP(), []int{7}
}

func (x *RepoExtended) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *RepoExtended) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *RepoExtended) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *RepoExtended) GetOrgId() string {
	if x != nil {
		return x.OrgId
	}
	return ""
}

func (x *RepoExtended) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *RepoExtended) GetHook() v1.RepoHook {
	if x != nil {
		return x.Hook
	}
	return v1.RepoHook(0)
}

func (x *RepoExtended) GetHookId() string {
	if x != nil {
		return x.HookId
	}
	return ""
}

func (x *RepoExtended) GetDefaultBranch() string {
	if x != nil {
		return x.DefaultBranch
	}
	return ""
}

func (x *RepoExtended) GetIsMonorepo() bool {
	if x != nil {
		return x.IsMonorepo
	}
	return false
}

func (x *RepoExtended) GetThreshold() int32 {
	if x != nil {
		return x.Threshold
	}
	return 0
}

func (x *RepoExtended) GetStaleDuration() *durationpb.Duration {
	if x != nil {
		return x.StaleDuration
	}
	return nil
}

func (x *RepoExtended) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *RepoExtended) GetIsActive() bool {
	if x != nil {
		return x.IsActive
	}
	return false
}

func (x *RepoExtended) GetChatEnabled() bool {
	if x != nil {
		return x.ChatEnabled
	}
	return false
}

func (x *RepoExtended) GetChannelName() string {
	if x != nil {
		return x.ChannelName
	}
	return ""
}

type ListReposResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Repos []*RepoExtended `protobuf:"bytes,1,rep,name=repos,proto3" json:"repos,omitempty"`
}

func (x *ListReposResponse) Reset() {
	*x = ListReposResponse{}
	mi := &file_ctrlplane_core_v1_repos_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListReposResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListReposResponse) ProtoMessage() {}

func (x *ListReposResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ctrlplane_core_v1_repos_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListReposResponse.ProtoReflect.Descriptor instead.
func (*ListReposResponse) Descriptor() ([]byte, []int) {
	return file_ctrlplane_core_v1_repos_proto_rawDescGZIP(), []int{8}
}

func (x *ListReposResponse) GetRepos() []*RepoExtended {
	if x != nil {
		return x.Repos
	}
	return nil
}

var File_ctrlplane_core_v1_repos_proto protoreflect.FileDescriptor

var file_ctrlplane_core_v1_repos_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2f, 0x63, 0x6f, 0x72, 0x65,
	0x2f, 0x76, 0x31, 0x2f, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x11, 0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e,
	0x76, 0x31, 0x1a, 0x1f, 0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2f, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x68, 0x6f, 0x6f, 0x6b, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xda, 0x03, 0x0a, 0x04, 0x52, 0x65, 0x70, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64,
	0x5f, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74,
	0x12, 0x15, 0x0a, 0x06, 0x6f, 0x72, 0x67, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x6f, 0x72, 0x67, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x31, 0x0a, 0x04, 0x68,
	0x6f, 0x6f, 0x6b, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1d, 0x2e, 0x63, 0x74, 0x72, 0x6c,
	0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e,
	0x52, 0x65, 0x70, 0x6f, 0x48, 0x6f, 0x6f, 0x6b, 0x52, 0x04, 0x68, 0x6f, 0x6f, 0x6b, 0x12, 0x17,
	0x0a, 0x07, 0x68, 0x6f, 0x6f, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x68, 0x6f, 0x6f, 0x6b, 0x49, 0x64, 0x12, 0x25, 0x0a, 0x0e, 0x64, 0x65, 0x66, 0x61, 0x75,
	0x6c, 0x74, 0x5f, 0x62, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0d, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x42, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x12, 0x1f,
	0x0a, 0x0b, 0x69, 0x73, 0x5f, 0x6d, 0x6f, 0x6e, 0x6f, 0x72, 0x65, 0x70, 0x6f, 0x18, 0x09, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x0a, 0x69, 0x73, 0x4d, 0x6f, 0x6e, 0x6f, 0x72, 0x65, 0x70, 0x6f, 0x12,
	0x1c, 0x0a, 0x09, 0x74, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x18, 0x0a, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x09, 0x74, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x12, 0x40, 0x0a,
	0x0e, 0x73, 0x74, 0x61, 0x6c, 0x65, 0x5f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x0d, 0x73, 0x74, 0x61, 0x6c, 0x65, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72,
	0x6c, 0x12, 0x1b, 0x0a, 0x09, 0x69, 0x73, 0x5f, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x18, 0x0d,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x69, 0x73, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x22, 0xb2,
	0x02, 0x0a, 0x11, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x31, 0x0a, 0x04, 0x68, 0x6f, 0x6f, 0x6b,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1d, 0x2e, 0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61,
	0x6e, 0x65, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x70,
	0x6f, 0x48, 0x6f, 0x6f, 0x6b, 0x52, 0x04, 0x68, 0x6f, 0x6f, 0x6b, 0x12, 0x17, 0x0a, 0x07, 0x68,
	0x6f, 0x6f, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x68, 0x6f,
	0x6f, 0x6b, 0x49, 0x64, 0x12, 0x25, 0x0a, 0x0e, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x5f,
	0x62, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x64, 0x65,
	0x66, 0x61, 0x75, 0x6c, 0x74, 0x42, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x12, 0x1f, 0x0a, 0x0b, 0x69,
	0x73, 0x5f, 0x6d, 0x6f, 0x6e, 0x6f, 0x72, 0x65, 0x70, 0x6f, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x0a, 0x69, 0x73, 0x4d, 0x6f, 0x6e, 0x6f, 0x72, 0x65, 0x70, 0x6f, 0x12, 0x1c, 0x0a, 0x09,
	0x74, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x09, 0x74, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x12, 0x40, 0x0a, 0x0e, 0x73, 0x74,
	0x61, 0x6c, 0x65, 0x5f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0d, 0x73,
	0x74, 0x61, 0x6c, 0x65, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x15, 0x0a, 0x06,
	0x6f, 0x72, 0x67, 0x5f, 0x69, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6f, 0x72,
	0x67, 0x49, 0x64, 0x22, 0x41, 0x0a, 0x12, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x70,
	0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2b, 0x0a, 0x04, 0x72, 0x65, 0x70,
	0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c,
	0x61, 0x6e, 0x65, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x70, 0x6f,
	0x52, 0x04, 0x72, 0x65, 0x70, 0x6f, 0x22, 0x24, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x52, 0x65, 0x70,
	0x6f, 0x42, 0x79, 0x49, 0x44, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x42, 0x0a, 0x13,
	0x47, 0x65, 0x74, 0x52, 0x65, 0x70, 0x6f, 0x42, 0x79, 0x49, 0x44, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x2b, 0x0a, 0x04, 0x72, 0x65, 0x70, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x17, 0x2e, 0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x63, 0x6f,
	0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x52, 0x04, 0x72, 0x65, 0x70, 0x6f,
	0x22, 0x32, 0x0a, 0x19, 0x47, 0x65, 0x74, 0x4f, 0x72, 0x67, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x42,
	0x79, 0x4f, 0x72, 0x67, 0x49, 0x44, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x15, 0x0a,
	0x06, 0x6f, 0x72, 0x67, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6f,
	0x72, 0x67, 0x49, 0x64, 0x22, 0x49, 0x0a, 0x1a, 0x47, 0x65, 0x74, 0x4f, 0x72, 0x67, 0x52, 0x65,
	0x70, 0x6f, 0x73, 0x42, 0x79, 0x4f, 0x72, 0x67, 0x49, 0x44, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x2b, 0x0a, 0x04, 0x72, 0x65, 0x70, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x17, 0x2e, 0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x63, 0x6f, 0x72,
	0x65, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x52, 0x04, 0x72, 0x65, 0x70, 0x6f, 0x22,
	0xa8, 0x04, 0x0a, 0x0c, 0x52, 0x65, 0x70, 0x6f, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x64, 0x65, 0x64,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x75,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x75, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x15, 0x0a, 0x06, 0x6f, 0x72, 0x67, 0x5f, 0x69, 0x64,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6f, 0x72, 0x67, 0x49, 0x64, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x31, 0x0a, 0x04, 0x68, 0x6f, 0x6f, 0x6b, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x1d, 0x2e, 0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x48, 0x6f, 0x6f, 0x6b, 0x52, 0x04,
	0x68, 0x6f, 0x6f, 0x6b, 0x12, 0x17, 0x0a, 0x07, 0x68, 0x6f, 0x6f, 0x6b, 0x5f, 0x69, 0x64, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x68, 0x6f, 0x6f, 0x6b, 0x49, 0x64, 0x12, 0x25, 0x0a,
	0x0e, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x5f, 0x62, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x42, 0x72,
	0x61, 0x6e, 0x63, 0x68, 0x12, 0x1f, 0x0a, 0x0b, 0x69, 0x73, 0x5f, 0x6d, 0x6f, 0x6e, 0x6f, 0x72,
	0x65, 0x70, 0x6f, 0x18, 0x09, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x69, 0x73, 0x4d, 0x6f, 0x6e,
	0x6f, 0x72, 0x65, 0x70, 0x6f, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f,
	0x6c, 0x64, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x74, 0x68, 0x72, 0x65, 0x73, 0x68,
	0x6f, 0x6c, 0x64, 0x12, 0x40, 0x0a, 0x0e, 0x73, 0x74, 0x61, 0x6c, 0x65, 0x5f, 0x64, 0x75, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0d, 0x73, 0x74, 0x61, 0x6c, 0x65, 0x44, 0x75, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x0c, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x1b, 0x0a, 0x09, 0x69, 0x73, 0x5f, 0x61, 0x63,
	0x74, 0x69, 0x76, 0x65, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x69, 0x73, 0x41, 0x63,
	0x74, 0x69, 0x76, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x68, 0x61, 0x74, 0x5f, 0x65, 0x6e, 0x61,
	0x62, 0x6c, 0x65, 0x64, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x63, 0x68, 0x61, 0x74,
	0x45, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x68, 0x61, 0x6e, 0x6e,
	0x65, 0x6c, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63,
	0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x4a, 0x0a, 0x11, 0x4c, 0x69,
	0x73, 0x74, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x35, 0x0a, 0x05, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f,
	0x2e, 0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x64, 0x65, 0x64, 0x52,
	0x05, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x32, 0x84, 0x03, 0x0a, 0x0b, 0x52, 0x65, 0x70, 0x6f, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x59, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x52, 0x65, 0x70, 0x6f, 0x12, 0x24, 0x2e, 0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65,
	0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52,
	0x65, 0x70, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x63, 0x74, 0x72,
	0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x5c, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x52, 0x65, 0x70, 0x6f, 0x42, 0x79, 0x49, 0x44,
	0x12, 0x25, 0x2e, 0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x63, 0x6f, 0x72,
	0x65, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x70, 0x6f, 0x42, 0x79, 0x49, 0x44,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x26, 0x2e, 0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c,
	0x61, 0x6e, 0x65, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x52,
	0x65, 0x70, 0x6f, 0x42, 0x79, 0x49, 0x44, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x71, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x4f, 0x72, 0x67, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x42, 0x79,
	0x4f, 0x72, 0x67, 0x49, 0x44, 0x12, 0x2c, 0x2e, 0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e,
	0x65, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x4f, 0x72, 0x67,
	0x52, 0x65, 0x70, 0x6f, 0x73, 0x42, 0x79, 0x4f, 0x72, 0x67, 0x49, 0x44, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x2d, 0x2e, 0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e,
	0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x4f, 0x72, 0x67, 0x52, 0x65,
	0x70, 0x6f, 0x73, 0x42, 0x79, 0x4f, 0x72, 0x67, 0x49, 0x44, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x49, 0x0a, 0x09, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x12,
	0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x24, 0x2e, 0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c,
	0x61, 0x6e, 0x65, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74,
	0x52, 0x65, 0x70, 0x6f, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0xc4, 0x01,
	0x0a, 0x15, 0x63, 0x6f, 0x6d, 0x2e, 0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e,
	0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x42, 0x0a, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x50, 0x72,
	0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x39, 0x67, 0x6f, 0x2e, 0x62, 0x72, 0x65, 0x75, 0x2e, 0x69,
	0x6f, 0x2f, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x6d, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61,
	0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e,
	0x65, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x3b, 0x63, 0x6f, 0x72, 0x65, 0x76, 0x31,
	0xa2, 0x02, 0x03, 0x43, 0x43, 0x58, 0xaa, 0x02, 0x11, 0x43, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61,
	0x6e, 0x65, 0x2e, 0x43, 0x6f, 0x72, 0x65, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x11, 0x43, 0x74, 0x72,
	0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x5c, 0x43, 0x6f, 0x72, 0x65, 0x5c, 0x56, 0x31, 0xe2, 0x02,
	0x1d, 0x43, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x5c, 0x43, 0x6f, 0x72, 0x65, 0x5c,
	0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02,
	0x13, 0x43, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x3a, 0x3a, 0x43, 0x6f, 0x72, 0x65,
	0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_ctrlplane_core_v1_repos_proto_rawDescOnce sync.Once
	file_ctrlplane_core_v1_repos_proto_rawDescData = file_ctrlplane_core_v1_repos_proto_rawDesc
)

func file_ctrlplane_core_v1_repos_proto_rawDescGZIP() []byte {
	file_ctrlplane_core_v1_repos_proto_rawDescOnce.Do(func() {
		file_ctrlplane_core_v1_repos_proto_rawDescData = protoimpl.X.CompressGZIP(file_ctrlplane_core_v1_repos_proto_rawDescData)
	})
	return file_ctrlplane_core_v1_repos_proto_rawDescData
}

var file_ctrlplane_core_v1_repos_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_ctrlplane_core_v1_repos_proto_goTypes = []any{
	(*Repo)(nil),                       // 0: ctrlplane.core.v1.Repo
	(*CreateRepoRequest)(nil),          // 1: ctrlplane.core.v1.CreateRepoRequest
	(*CreateRepoResponse)(nil),         // 2: ctrlplane.core.v1.CreateRepoResponse
	(*GetRepoByIDRequest)(nil),         // 3: ctrlplane.core.v1.GetRepoByIDRequest
	(*GetRepoByIDResponse)(nil),        // 4: ctrlplane.core.v1.GetRepoByIDResponse
	(*GetOrgReposByOrgIDRequest)(nil),  // 5: ctrlplane.core.v1.GetOrgReposByOrgIDRequest
	(*GetOrgReposByOrgIDResponse)(nil), // 6: ctrlplane.core.v1.GetOrgReposByOrgIDResponse
	(*RepoExtended)(nil),               // 7: ctrlplane.core.v1.RepoExtended
	(*ListReposResponse)(nil),          // 8: ctrlplane.core.v1.ListReposResponse
	(*timestamppb.Timestamp)(nil),      // 9: google.protobuf.Timestamp
	(v1.RepoHook)(0),                   // 10: ctrlplane.events.v1.RepoHook
	(*durationpb.Duration)(nil),        // 11: google.protobuf.Duration
	(*emptypb.Empty)(nil),              // 12: google.protobuf.Empty
}
var file_ctrlplane_core_v1_repos_proto_depIdxs = []int32{
	9,  // 0: ctrlplane.core.v1.Repo.created_at:type_name -> google.protobuf.Timestamp
	9,  // 1: ctrlplane.core.v1.Repo.updated_at:type_name -> google.protobuf.Timestamp
	10, // 2: ctrlplane.core.v1.Repo.hook:type_name -> ctrlplane.events.v1.RepoHook
	11, // 3: ctrlplane.core.v1.Repo.stale_duration:type_name -> google.protobuf.Duration
	10, // 4: ctrlplane.core.v1.CreateRepoRequest.hook:type_name -> ctrlplane.events.v1.RepoHook
	11, // 5: ctrlplane.core.v1.CreateRepoRequest.stale_duration:type_name -> google.protobuf.Duration
	0,  // 6: ctrlplane.core.v1.CreateRepoResponse.repo:type_name -> ctrlplane.core.v1.Repo
	0,  // 7: ctrlplane.core.v1.GetRepoByIDResponse.repo:type_name -> ctrlplane.core.v1.Repo
	0,  // 8: ctrlplane.core.v1.GetOrgReposByOrgIDResponse.repo:type_name -> ctrlplane.core.v1.Repo
	9,  // 9: ctrlplane.core.v1.RepoExtended.created_at:type_name -> google.protobuf.Timestamp
	9,  // 10: ctrlplane.core.v1.RepoExtended.updated_at:type_name -> google.protobuf.Timestamp
	10, // 11: ctrlplane.core.v1.RepoExtended.hook:type_name -> ctrlplane.events.v1.RepoHook
	11, // 12: ctrlplane.core.v1.RepoExtended.stale_duration:type_name -> google.protobuf.Duration
	7,  // 13: ctrlplane.core.v1.ListReposResponse.repos:type_name -> ctrlplane.core.v1.RepoExtended
	1,  // 14: ctrlplane.core.v1.RepoService.CreateRepo:input_type -> ctrlplane.core.v1.CreateRepoRequest
	3,  // 15: ctrlplane.core.v1.RepoService.GetRepoByID:input_type -> ctrlplane.core.v1.GetRepoByIDRequest
	5,  // 16: ctrlplane.core.v1.RepoService.GetOrgReposByOrgID:input_type -> ctrlplane.core.v1.GetOrgReposByOrgIDRequest
	12, // 17: ctrlplane.core.v1.RepoService.ListRepos:input_type -> google.protobuf.Empty
	2,  // 18: ctrlplane.core.v1.RepoService.CreateRepo:output_type -> ctrlplane.core.v1.CreateRepoResponse
	4,  // 19: ctrlplane.core.v1.RepoService.GetRepoByID:output_type -> ctrlplane.core.v1.GetRepoByIDResponse
	6,  // 20: ctrlplane.core.v1.RepoService.GetOrgReposByOrgID:output_type -> ctrlplane.core.v1.GetOrgReposByOrgIDResponse
	8,  // 21: ctrlplane.core.v1.RepoService.ListRepos:output_type -> ctrlplane.core.v1.ListReposResponse
	18, // [18:22] is the sub-list for method output_type
	14, // [14:18] is the sub-list for method input_type
	14, // [14:14] is the sub-list for extension type_name
	14, // [14:14] is the sub-list for extension extendee
	0,  // [0:14] is the sub-list for field type_name
}

func init() { file_ctrlplane_core_v1_repos_proto_init() }
func file_ctrlplane_core_v1_repos_proto_init() {
	if File_ctrlplane_core_v1_repos_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_ctrlplane_core_v1_repos_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_ctrlplane_core_v1_repos_proto_goTypes,
		DependencyIndexes: file_ctrlplane_core_v1_repos_proto_depIdxs,
		MessageInfos:      file_ctrlplane_core_v1_repos_proto_msgTypes,
	}.Build()
	File_ctrlplane_core_v1_repos_proto = out.File
	file_ctrlplane_core_v1_repos_proto_rawDesc = nil
	file_ctrlplane_core_v1_repos_proto_goTypes = nil
	file_ctrlplane_core_v1_repos_proto_depIdxs = nil
}
