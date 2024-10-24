package convert

import (
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/timestamppb"

	"go.breu.io/quantm/internal/db/entities"
	authv1 "go.breu.io/quantm/internal/nomad/proto/ctrlplane/auth/v1"
)

func UserToProto(user *entities.User) *authv1.User {
	return &authv1.User{
		Id:         UUIDToProto(user.ID),
		CreatedAt:  timestamppb.New(user.CreatedAt),
		UpdatedAt:  timestamppb.New(user.UpdatedAt),
		FirstName:  user.FirstName.String,
		LastName:   user.LastName.String,
		Email:      user.Email,
		Password:   user.Password.String,
		IsActive:   user.IsActive,
		IsVerified: user.IsVerified,
	}
}

func ProtoToUser(proto *authv1.User) *entities.User {
	return &entities.User{
		ID:         ProtoToUUID(proto.Id),
		CreatedAt:  proto.CreatedAt.AsTime(),
		UpdatedAt:  proto.UpdatedAt.AsTime(),
		FirstName:  pgtype.Text{String: proto.FirstName, Valid: true},
		LastName:   pgtype.Text{String: proto.LastName, Valid: true},
		Password:   pgtype.Text{String: proto.Password, Valid: true},
		Email:      proto.Email,
		IsActive:   proto.IsActive,
		IsVerified: proto.IsVerified,
	}
}

func ProtoToCreateUserParams(proto *authv1.CreateUserRequest) entities.CreateUserParams {
	return entities.CreateUserParams{
		FirstName: pgtype.Text{String: proto.FirstName, Valid: true},
		LastName:  pgtype.Text{String: proto.LastName, Valid: true},
		Email:     proto.Email,
		Password:  pgtype.Text{String: proto.Password, Valid: true},
	}
}
