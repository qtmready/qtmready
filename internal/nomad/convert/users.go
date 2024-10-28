package convert

import (
	"encoding/json"

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

// AuthUserQueryToProto converts a user, accounts, teams, and org byte slices to an authv1.AuthUser protobuf message.
func AuthUserQueryToProto(user, accounts, teams, org []byte) (*authv1.AuthUser, error) {
	converted := authv1.AuthUser{
		Teams:    make([]*authv1.Team, 0),
		Accounts: make([]*authv1.Account, 0),
	}

	if err := json.Unmarshal(user, converted.User); err != nil {
		return nil, err
	}

	if err := BytesToSliceTeamProto(teams, converted.Teams); err != nil {
		return nil, err
	}

	if err := BytesToSliceAccountProto(accounts, converted.Accounts); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(org, converted.Org); err != nil {
		return nil, err
	}

	return &converted, nil
}

// BytesToSliceTeamProto converts a byte slice representing a JSON array of Team proto messages to a slice of
// pointers to Team proto messages.
//
// It unmarshals the JSON data into a temporary slice of Team proto messages and then appends pointers to each
// element of the temporary slice to the target slice. This approach ensures that memory is allocated correctly for
// the structs and that the pointers are referencing the correct locations, preventing potential data loss.
//
// Note that since slices are reference types in Go, the target slice will be modified in place.
func BytesToSliceTeamProto(src []byte, tgt []*authv1.Team) error {
	var deserialized []authv1.Team
	if err := json.Unmarshal(src, &deserialized); err != nil {
		return err
	}

	for idx := range deserialized {
		tgt = append(tgt, &deserialized[idx])
	}

	return nil
}

// BytesToSliceAccountProto converts a byte slice representing a JSON array of Account proto messages to a slice of
// pointers to Account proto messages.
//
// It unmarshals the JSON data into a temporary slice of Team proto messages and then appends pointers to each
// element of the temporary slice to the target slice. This approach ensures that memory is allocated correctly for
// the structs and that the pointers are referencing the correct locations, preventing potential data loss.
//
// Note that since slices are reference types in Go, the target slice will be modified in place.
func BytesToSliceAccountProto(src []byte, tgt []*authv1.Account) error {
	deserialized := make([]authv1.Account, 0)
	if err := json.Unmarshal(src, &deserialized); err != nil {
		return err
	}

	for idx := range deserialized {
		tgt = append(tgt, &deserialized[idx])
	}

	return nil
}
