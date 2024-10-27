package convert

import (
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"

	"go.breu.io/quantm/internal/db/entities"
	authv1 "go.breu.io/quantm/internal/nomad/proto/ctrlplane/auth/v1"
)

// UserToProto converts a User entity to its protobuf representation.
//
// It maps all fields from the User entity to corresponding fields in the authv1.User protobuf message.
func UserToProto(user *entities.User) *authv1.User {
	return &authv1.User{
		Id:         UUIDToProto(user.ID),
		CreatedAt:  timestamppb.New(user.CreatedAt),
		UpdatedAt:  timestamppb.New(user.UpdatedAt),
		OrgId:      UUIDToProto(user.OrgID),
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Email:      user.Email,
		IsActive:   user.IsActive,
		IsVerified: user.IsVerified,
	}
}

// ProtoToUser converts a authv1.User protobuf message to a User entity.
//
// It maps all fields from the authv1.User protobuf message to corresponding fields in the User entity.
func ProtoToUser(proto *authv1.User) *entities.User {
	return &entities.User{
		ID:         ProtoToUUID(proto.GetId()),
		CreatedAt:  proto.GetCreatedAt().AsTime(),
		UpdatedAt:  proto.GetUpdatedAt().AsTime(),
		OrgID:      ProtoToUUID(proto.GetOrgId()),
		FirstName:  proto.GetFirstName(),
		LastName:   proto.GetLastName(),
		Email:      proto.GetEmail(),
		IsActive:   proto.GetIsActive(),
		IsVerified: proto.GetIsVerified(),
	}
}

// ProtoToCreateUserParams converts a CreateUserRequest protobuf message to CreateUserParams.
//
// It maps the first name, last name, and email from the protobuf message to the corresponding fields in the
// CreateUserParams. The password is hashed using bcrypt.DefaultCost.
//
// TODO: Implement actual password hashing using the provided password in the protobuf message.
func ProtoToCreateUserParams(proto *authv1.CreateUserRequest) entities.CreateUserParams {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(""), bcrypt.DefaultCost) // TODO: hash password

	return entities.CreateUserParams{
		FirstName: proto.GetFirstName(),
		LastName:  proto.GetLastName(),
		Email:     proto.GetEmail(),
		Password:  string(hashed),
	}
}

// ProtoToUpdateUserParams converts an UpdateUserRequest protobuf message to UpdateUserParams.
//
// It maps the user ID, first name, last name, email, and organization ID from the protobuf message to the corresponding
// fields in the UpdateUserParams.
func ProtoToUpdateUserParams(proto *authv1.UpdateUserRequest) entities.UpdateUserParams {
	return entities.UpdateUserParams{
		ID:        ProtoToUUID(proto.User.GetId()),
		FirstName: proto.User.GetFirstName(),
		LastName:  proto.User.GetLastName(),
		Lower:     proto.User.GetEmail(),
		OrgID:     ProtoToUUID(proto.User.GetOrgId()),
	}
}
