package cast

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/db/entities"
	corev1 "go.breu.io/quantm/internal/proto/ctrlplane/core/v1"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

// HookToProto converts an int32 representation of a RepoHook to a RepoHook proto.
func HookToProto(hook int32) eventsv1.RepoHook {
	v, ok := eventsv1.RepoHook_name[hook]
	if !ok {
		return eventsv1.RepoHook_REPO_HOOK_UNSPECIFIED
	}

	return eventsv1.RepoHook(eventsv1.RepoHook_value[v])
}

// RepoToProto converts a Repo entity to a Repo proto.
func RepoToProto(repo *entities.Repo) *corev1.Repo {
	return &corev1.Repo{
		Id:            repo.ID.String(),
		CreatedAt:     timestamppb.New(repo.CreatedAt),
		UpdatedAt:     timestamppb.New(repo.UpdatedAt),
		OrgId:         repo.OrgID.String(),
		Name:          repo.Name,
		Hook:          HookToProto(repo.Hook),
		HookId:        repo.HookID.String(),
		DefaultBranch: repo.DefaultBranch,
		IsMonorepo:    repo.IsMonorepo,
		Threshold:     repo.Threshold,
		StaleDuration: db.IntervalToProto(repo.StaleDuration),
		Url:           repo.Url,
		IsActive:      repo.IsActive,
	}
}

// RepoToProto converts a Repo entity to a Repo proto.
func RepoMesgingToProto(repo *entities.ListReposRow) *corev1.RepoHasMesging {
	return &corev1.RepoHasMesging{
		Id:            repo.ID.String(),
		CreatedAt:     timestamppb.New(repo.CreatedAt),
		UpdatedAt:     timestamppb.New(repo.UpdatedAt),
		OrgId:         repo.OrgID.String(),
		Name:          repo.Name,
		Hook:          HookToProto(repo.Hook),
		HookId:        repo.HookID.String(),
		DefaultBranch: repo.DefaultBranch,
		IsMonorepo:    repo.IsMonorepo,
		Threshold:     repo.Threshold,
		StaleDuration: db.IntervalToProto(repo.StaleDuration),
		Url:           repo.Url,
		IsActive:      repo.IsActive,
		HasMesging:    repo.HasMesging,
	}
}

// ReposToProto converts a slice of Repo entities to a slice of Repo protos.
func RepoHasMesgingToProto(repos []entities.ListReposRow) []*corev1.RepoHasMesging {
	protos := make([]*corev1.RepoHasMesging, 0)
	for _, repo := range repos {
		protos = append(protos, RepoMesgingToProto(&repo))
	}

	return protos
}
