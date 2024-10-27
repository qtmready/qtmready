package convert_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/nomad/convert"
	authv1 "go.breu.io/quantm/internal/nomad/proto/ctrlplane/auth/v1"
)

func TestUserToProto(t *testing.T) {
	t.Parallel()

	now := time.Now()
	usr := &entities.User{
		ID:         uuid.New(),
		CreatedAt:  now,
		UpdatedAt:  now,
		OrgID:      uuid.New(),
		FirstName:  "John",
		LastName:   "Doe",
		Email:      "john.doe@example.com",
		IsActive:   true,
		IsVerified: false,
	}

	proto := convert.UserToProto(usr)

	if proto.GetId().GetValue() != usr.ID.String() {
		t.Errorf("Expected ID to be %s, got %s", usr.ID.String(), proto.GetId().GetValue())
	}

	if proto.GetOrgId().GetValue() != usr.OrgID.String() {
		t.Errorf("Expected OrgID to be %s, got %s", usr.OrgID.String(), proto.GetOrgId().GetValue())
	}

	if proto.GetFirstName() != usr.FirstName {
		t.Errorf("Expected FirstName to be %s, got %s", usr.FirstName, proto.GetFirstName())
	}

	if proto.GetLastName() != usr.LastName {
		t.Errorf("Expected LastName to be %s, got %s", usr.LastName, proto.GetLastName())
	}

	if proto.GetEmail() != usr.Email {
		t.Errorf("Expected Email to be %s, got %s", usr.Email, proto.GetEmail())
	}

	if proto.GetIsActive() != usr.IsActive {
		t.Errorf("Expected IsActive to be %v, got %v", usr.IsActive, proto.GetIsActive())
	}

	if proto.GetIsVerified() != usr.IsVerified {
		t.Errorf("Expected IsVerified to be %v, got %v", usr.IsVerified, proto.GetIsVerified())
	}
}

func TestProtoToUser(t *testing.T) {
	t.Parallel()

	now := time.Now()
	proto := &authv1.User{
		Id:         convert.UUIDToProto(uuid.New()),
		CreatedAt:  timestamppb.New(now),
		UpdatedAt:  timestamppb.New(now),
		OrgId:      convert.UUIDToProto(uuid.New()),
		FirstName:  "John",
		LastName:   "Doe",
		Email:      "john.doe@example.com",
		IsActive:   true,
		IsVerified: false,
	}

	usr := convert.ProtoToUser(proto)

	if usr.ID != convert.ProtoToUUID(proto.GetId()) {
		t.Errorf("Expected ID to be %s, got %s", proto.GetId().GetValue(), usr.ID.String())
	}

	if !usr.CreatedAt.Equal(proto.GetCreatedAt().AsTime()) {
		t.Errorf("Expected CreatedAt to be %s, got %s", proto.GetCreatedAt().AsTime().Format(time.RFC3339), usr.CreatedAt.Format(time.RFC3339))
	}

	if !usr.UpdatedAt.Equal(proto.GetUpdatedAt().AsTime()) {
		t.Errorf("Expected UpdatedAt to be %s, got %s", proto.GetUpdatedAt().AsTime().Format(time.RFC3339), usr.UpdatedAt.Format(time.RFC3339))
	}

	if usr.OrgID != convert.ProtoToUUID(proto.GetOrgId()) {
		t.Errorf("Expected OrgID to be %s, got %s", proto.GetOrgId().GetValue(), usr.OrgID.String())
	}

	if usr.FirstName != proto.GetFirstName() {
		t.Errorf("Expected FirstName to be %s, got %s", proto.GetFirstName(), usr.FirstName)
	}

	if usr.LastName != proto.GetLastName() {
		t.Errorf("Expected LastName to be %s, got %s", proto.GetLastName(), usr.LastName)
	}

	if usr.Email != proto.GetEmail() {
		t.Errorf("Expected Email to be %s, got %s", proto.GetEmail(), usr.Email)
	}

	if usr.IsActive != proto.GetIsActive() {
		t.Errorf("Expected IsActive to be %v, got %v", proto.GetIsActive(), usr.IsActive)
	}

	if usr.IsVerified != proto.GetIsVerified() {
		t.Errorf("Expected IsVerified to be %v, got %v", proto.GetIsVerified(), usr.IsVerified)
	}
}

func TestProtoToCreateUserParams(t *testing.T) {
	t.Parallel()

	req := &authv1.CreateUserRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}

	params := convert.ProtoToCreateUserParams(req)

	if params.FirstName != req.GetFirstName() {
		t.Errorf("Expected FirstName to be %s, got %s", req.GetFirstName(), params.FirstName)
	}

	if params.LastName != req.GetLastName() {
		t.Errorf("Expected LastName to be %s, got %s", req.GetLastName(), params.LastName)
	}

	if params.Email != req.GetEmail() {
		t.Errorf("Expected Email to be %s, got %s", req.GetEmail(), params.Email)
	}
}

func TestProtoToUpdateUserParams(t *testing.T) {
	t.Parallel()

	proto := &authv1.User{
		Id:        convert.UUIDToProto(uuid.New()),
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		OrgId:     convert.UUIDToProto(uuid.New()),
	}

	req := &authv1.UpdateUserRequest{
		User: proto,
	}

	params := convert.ProtoToUpdateUserParams(req)

	if params.ID != convert.ProtoToUUID(proto.GetId()) {
		t.Errorf("Expected ID to be %s, got %s", proto.GetId().GetValue(), params.ID.String())
	}

	if params.FirstName != proto.GetFirstName() {
		t.Errorf("Expected FirstName to be %s, got %s", proto.GetFirstName(), params.FirstName)
	}

	if params.LastName != proto.GetLastName() {
		t.Errorf("Expected LastName to be %s, got %s", proto.GetLastName(), params.LastName)
	}

	if params.Lower != proto.GetEmail() {
		t.Errorf("Expected Lower to be %s, got %s", proto.GetEmail(), params.Lower)
	}

	if params.OrgID != convert.ProtoToUUID(proto.GetOrgId()) {
		t.Errorf("Expected OrgID to be %s, got %s", proto.GetOrgId().GetValue(), params.OrgID.String())
	}
}
