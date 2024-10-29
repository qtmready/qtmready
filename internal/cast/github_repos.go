package cast

import (
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/db/entities"
	githubv1 "go.breu.io/quantm/internal/proto/hooks/github/v1"
)

// GithubRepoToProto converts an GithubRepo entity to its protobuf representation.
//
// It maps all fields from the User entity to corresponding fields in the githubv1.GithubRepo protobuf message.
func GithubRepoToProto(repo *entities.GithubRepo) *githubv1.GithubRepo {
	return &githubv1.GithubRepo{
		Id:             repo.ID.String(),
		CreatedAt:      timestamppb.New(repo.CreatedAt),
		UpdatedAt:      timestamppb.New(repo.UpdatedAt),
		Name:           repo.Name,
		RepoId:         db.ConvertPgTypeUUIDToUUID(repo.RepoID).String(),
		InstallationId: repo.InstallationID.String(),
		GithubId:       repo.GithubID,
		Url:            repo.Url,
		IsActive:       repo.IsActive.Bool,
	}
}

// ProtoToGithubRepo converts a githubv1.GithubRepo protobuf message to a Repo entity.
//
// It maps all fields from the githubv1.GithubRepo protobuf message to corresponding fields in the Repo entity.
func ProtoToGithubRepo(proto *githubv1.GithubRepo) *entities.GithubRepo {
	return &entities.GithubRepo{
		ID:             uuid.MustParse(proto.GetId()),
		CreatedAt:      proto.GetCreatedAt().AsTime(),
		UpdatedAt:      proto.GetUpdatedAt().AsTime(),
		Name:           proto.GetName(),
		RepoID:         db.ConvertStringToPgTypeUUID(proto.GetRepoId()), // Convert string to pgtype.UUID
		InstallationID: uuid.MustParse(proto.GetInstallationId()),
		GithubID:       proto.GetGithubId(),
		Url:            proto.GetUrl(),
		IsActive:       db.ConvertBoolToPgTypeBool(proto.GetIsActive()), // Convert bool to pgtype.Bool
	}
}
