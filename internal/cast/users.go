package cast

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"

	"go.breu.io/quantm/internal/db/entities"
	authv1 "go.breu.io/quantm/internal/proto/ctrlplane/auth/v1"
)

type (
	AuthOrg struct {
		ID        string           `json:"id"`
		CreatedAt time.Time        `json:"created_at"`
		UpdatedAt time.Time        `json:"updated_at"`
		Name      string           `json:"name"`
		Domain    string           `json:"domain"`
		Slug      string           `json:"slug"`
		Hooks     *authv1.OrgHooks `json:"hooks"`
	}
)

// UserToProto converts a User entity to its protobuf representation.
//
// It maps all fields from the User entity to corresponding fields in the authv1.User protobuf message.
func UserToProto(user *entities.User) *authv1.User {
	return &authv1.User{
		Id:         user.ID.String(),
		CreatedAt:  timestamppb.New(user.CreatedAt),
		UpdatedAt:  timestamppb.New(user.UpdatedAt),
		OrgId:      user.OrgID.String(),
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Email:      user.Email,
		Picture:    user.Picture,
		IsActive:   user.IsActive,
		IsVerified: user.IsVerified,
	}
}

// ProtoToUser converts a authv1.User protobuf message to a User entity.
//
// It maps all fields from the authv1.User protobuf message to corresponding fields in the User entity.
func ProtoToUser(proto *authv1.User) *entities.User {
	return &entities.User{
		ID:         uuid.MustParse(proto.GetId()),
		CreatedAt:  proto.GetCreatedAt().AsTime(),
		UpdatedAt:  proto.GetUpdatedAt().AsTime(),
		OrgID:      uuid.MustParse(proto.GetOrgId()),
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
		Lower:     proto.GetEmail(),
		Picture:   proto.GetPicture(),
		Password:  string(hashed),
	}
}

// ProtoToUpdateUserParams converts an UpdateUserRequest protobuf message to UpdateUserParams.
//
// It maps the user ID, first name, last name, email, and organization ID from the protobuf message to the corresponding
// fields in the UpdateUserParams.
func ProtoToUpdateUserParams(proto *authv1.UpdateUserRequest) entities.UpdateUserParams {
	return entities.UpdateUserParams{
		ID:        uuid.MustParse(proto.User.GetId()),
		FirstName: proto.User.GetFirstName(),
		LastName:  proto.User.GetLastName(),
		Lower:     proto.User.GetEmail(),
		OrgID:     uuid.MustParse(proto.User.GetOrgId()),
	}
}

// AuthUserQueryResponseToProto converts a user, accounts, teams, and org byte slices to an authv1.AuthUser protobuf message.
func AuthUserQueryResponseToProto(user, orgs, roles, accounts, teams []byte) (*authv1.AuthUser, error) {
	response := &authv1.AuthUser{}

	usr := &entities.User{}
	if err := json.Unmarshal(user, usr); err != nil {
		slog.Error("unmarshalling user", "error", err)
		return nil, err
	}

	response.User = UserToProto(usr)

	org := &AuthOrg{}
	if err := json.Unmarshal(orgs, org); err != nil {
		slog.Error("unmarshalling org", "error", err)
		return nil, err
	}

	response.Org = &authv1.Org{
		Id:        org.ID,
		CreatedAt: timestamppb.New(org.CreatedAt),
		UpdatedAt: timestamppb.New(org.UpdatedAt),
		Name:      org.Name,
		Domain:    org.Domain,
		Slug:      org.Slug,
		Hooks:     org.Hooks,
	}

	rls, err := BytesToStringSlice(roles)
	if err != nil {
		return response, err
	}

	response.Roles = rls

	if tms, err := BytesToTeamSliceProto(teams); err != nil {
		return response, err
	} else {
		response.Teams = tms
	}

	if acts, err := BytesToAccountSliceProto(accounts); err != nil {
		return response, err
	} else {
		response.Accounts = acts
	}

	return response, nil
}

// BytesToTeamSliceProto converts a byte slice representing a JSON array of Team proto messages to a slice of
// pointers to Team proto messages.
//
// It unmarshals the JSON data into a temporary slice of Team proto messages and then appends pointers to each
// element of the temporary slice to the target slice. This approach ensures that memory is allocated correctly for
// the structs and that the pointers are referencing the correct locations, preventing potential data loss.
//
// Note that since slices are reference types in Go, the target slice will be modified in place.
func BytesToTeamSliceProto(src []byte) ([]*authv1.Team, error) {
	response := make([]*authv1.Team, 0)

	deserialized := make([]entities.Team, 0)

	if err := json.Unmarshal(src, &deserialized); err != nil {
		return response, err // pg hack.
	}

	for idx := range deserialized {
		response = append(response, TeamToProto(&deserialized[idx]))
	}

	return response, nil
}

// BytesToAccountSliceProto converts a byte slice representing a JSON array of Account proto messages to a slice of
// pointers to Account proto messages.
//
// It unmarshals the JSON data into a temporary slice of Team proto messages and then appends pointers to each
// element of the temporary slice to the target slice. This approach ensures that memory is allocated correctly for
// the structs and that the pointers are referencing the correct locations, preventing potential data loss.
//
// Note that since slices are reference types in Go, the target slice will be modified in place.
func BytesToAccountSliceProto(src []byte) ([]*authv1.Account, error) {
	response := make([]*authv1.Account, 0)

	deserialized := make([]entities.OauthAccount, 0)

	if err := json.Unmarshal(src, &deserialized); err != nil {
		return response, err
	}

	for idx := range deserialized {
		response = append(response, AccountToProto(&deserialized[idx]))
	}

	return response, nil
}

func BytesToStringSlice(src []byte) ([]string, error) {
	var response []string
	if err := json.Unmarshal(src, &response); err != nil {
		return response, err
	}

	return response, nil
}
