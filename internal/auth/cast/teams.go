package authcast

import (
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"go.breu.io/quantm/internal/db/entities"
	authv1 "go.breu.io/quantm/internal/proto/ctrlplane/auth/v1"
)

// ProtoToTeam converts a protobuf team to its Team entity representation.
func ProtoToTeam(proto *authv1.Team) *entities.Team {
	return &entities.Team{
		ID:        uuid.MustParse(proto.GetId()),
		CreatedAt: proto.GetCreatedAt().AsTime(),
		UpdatedAt: proto.GetUpdatedAt().AsTime(),
		Name:      proto.GetName(),
		Slug:      proto.GetSlug(),
	}
}

// TeamToProto converts a Team entity to its protobuf representation.
func TeamToProto(team *entities.Team) *authv1.Team {
	return &authv1.Team{
		Id:        team.ID.String(),
		CreatedAt: timestamppb.New(team.CreatedAt),
		UpdatedAt: timestamppb.New(team.UpdatedAt),
		Name:      team.Name,
		Slug:      team.Slug,
	}
}

// ProtoToCreateTeamParams converts a protobuf CreateTeamRequest to a CreateTeamParams.
func ProtoToCreateTeamParams(proto *authv1.CreateTeamRequest) entities.CreateTeamParams {
	return entities.CreateTeamParams{
		OrgID: uuid.MustParse(proto.GetOrgId()),
		Name:  proto.GetName(),
	}
}

// GetTeamBySlugRowToProto converts a GetTeamBySlugRow entity to its protobuf representation.
func GetTeamBySlugRowToProto(team entities.GetTeamBySlugRow) *authv1.Team {
	return &authv1.Team{
		Id:   team.ID.String(),
		Name: team.Name,
	}
}
