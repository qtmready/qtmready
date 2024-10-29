package cast

import (
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/db/entities"
	corev1 "go.breu.io/quantm/internal/proto/ctrlplane/core/v1"
)

// CoreRepoToProto converts an Repo entity to its protobuf representation.
func CoreRepoToProto(repo *entities.Repo) *corev1.Repo {
	return &corev1.Repo{
		Id:            repo.ID.String(),
		CreatedAt:     timestamppb.New(repo.CreatedAt),
		UpdatedAt:     timestamppb.New(repo.UpdatedAt),
		Name:          repo.Name,
		Provider:      repo.Hook,
		ProviderId:    repo.HookID,
		DefaultBranch: repo.DefaultBranch,
		Threshold:     repo.Threshold,
		StaleDuration: db.IntervalToDurationString(repo.StaleDuration),
		OrgId:         repo.OrgID.String(),
	}
}

// ProtoToCoreRepo converts a corev1.Repo protobuf message to a Repo entity.
//
// It maps all fields from the corev1.Repo protobuf message to corresponding fields in the Repo entity.
func ProtoToCoreRepo(proto *corev1.Repo) *entities.Repo {
	return &entities.Repo{
		ID:            uuid.MustParse(proto.GetId()),
		CreatedAt:     proto.GetCreatedAt().AsTime(),
		UpdatedAt:     proto.GetUpdatedAt().AsTime(),
		OrgID:         uuid.MustParse(proto.GetOrgId()),
		Name:          proto.GetName(),
		Hook:          proto.GetProvider(),
		HookID:        proto.GetProviderId(),
		DefaultBranch: proto.GetDefaultBranch(),
		Threshold:     proto.GetThreshold(),
		StaleDuration: db.StringToInterval(proto.GetStaleDuration()),
	}
}
