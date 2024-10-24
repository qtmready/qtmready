package convert

import (
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"

	"go.breu.io/quantm/internal/db/entities"
	authv1 "go.breu.io/quantm/internal/nomad/proto/ctrlplane/auth/v1"
)

func UserToProto(user *entities.User) *authv1.User {
	return &authv1.User{
		Id:         UUIDToProto(user.ID),
		CreatedAt:  timestamppb.New(user.CreatedAt),
		UpdatedAt:  timestamppb.New(user.UpdatedAt),
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Email:      user.Email,
		IsActive:   user.IsActive,
		IsVerified: user.IsVerified,
	}
}

func ProtoToUser(proto *authv1.User) *entities.User {
	return &entities.User{
		ID:         ProtoToUUID(proto.GetId()),
		CreatedAt:  proto.GetCreatedAt().AsTime(),
		UpdatedAt:  proto.GetUpdatedAt().AsTime(),
		FirstName:  proto.GetFirstName(),
		LastName:   proto.GetLastName(),
		Email:      proto.GetEmail(),
		IsActive:   proto.GetIsActive(),
		IsVerified: proto.GetIsVerified(),
	}
}

func ProtoToCreateUserParams(proto *authv1.CreateUserRequest) entities.CreateUserParams {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(""), bcrypt.DefaultCost) // TODO: hash password

	return entities.CreateUserParams{
		FirstName: proto.GetFirstName(),
		LastName:  proto.GetLastName(),
		Email:     proto.GetEmail(),
		Password:  string(hashed),
	}
}
