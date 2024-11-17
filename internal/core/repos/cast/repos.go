package reposcast

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/db/entities"
	corev1 "go.breu.io/quantm/internal/proto/ctrlplane/core/v1"
)

func RepoToProto(repo *entities.Repo) *corev1.Repo {
	return &corev1.Repo{
		Id:        repo.ID.String(),
		CreatedAt: timestamppb.New(repo.CreatedAt),
		UpdatedAt: timestamppb.New(repo.UpdatedAt),
		OrgId:     repo.OrgID.String(),
		Name:      repo.Name,
		// Hook:          repo.Hook,
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
