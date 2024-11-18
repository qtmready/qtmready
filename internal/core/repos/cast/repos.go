package reposcast

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/db/entities"
	corev1 "go.breu.io/quantm/internal/proto/ctrlplane/core/v1"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

func HookToProto(hook int32) eventsv1.RepoHook {
	v, ok := eventsv1.RepoHook_name[hook]
	if !ok {
		return eventsv1.RepoHook_REPO_HOOK_UNSPECIFIED
	}

	return eventsv1.RepoHook(eventsv1.RepoHook_value[v])
}

func RepoToProto(repo *entities.Repo) *corev1.Repo {
	h := HookToProto(repo.Hook)

	return &corev1.Repo{
		Id:            repo.ID.String(),
		CreatedAt:     timestamppb.New(repo.CreatedAt),
		UpdatedAt:     timestamppb.New(repo.UpdatedAt),
		OrgId:         repo.OrgID.String(),
		Name:          repo.Name,
		Hook:          h.String(),
		HookId:        repo.HookID.String(),
		DefaultBranch: repo.DefaultBranch,
		IsMonorepo:    repo.IsMonorepo,
		Threshold:     repo.Threshold,
		StaleDuration: db.IntervalToProto(repo.StaleDuration),
		Url:           repo.Url,
		IsActive:      repo.IsActive,
	}
}

func ReposToProto(repos []entities.Repo) []*corev1.Repo {
	protos := make([]*corev1.Repo, 0)
	for _, repo := range repos {
		protos = append(protos, RepoToProto(&repo))
	}

	return protos
}
