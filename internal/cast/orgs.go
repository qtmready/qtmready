package cast

import (
	"encoding/json"
	"log/slog"

	"google.golang.org/protobuf/types/known/timestamppb"

	"go.breu.io/quantm/internal/db/entities"
	authv1 "go.breu.io/quantm/internal/proto/ctrlplane/auth/v1"
	commonv1 "go.breu.io/quantm/internal/proto/ctrlplane/common/v1"
)

func OrgToProto(org *entities.Org) *authv1.Org {
	unmarshalled := &authv1.OrgHooks{}

	if err := json.Unmarshal(org.Hooks, unmarshalled); err != nil {
		slog.Error("failed to unmarshal org hooks", "error", err)
	}

	hooks := authv1.OrgHooks{
		Repo:      commonv1.RepoHook(unmarshalled.Repo),
		Messaging: commonv1.MessagingHook(unmarshalled.Messaging),
	}

	return &authv1.Org{
		Id:        org.ID.String(),
		CreatedAt: timestamppb.New(org.CreatedAt),
		UpdatedAt: timestamppb.New(org.UpdatedAt),
		Name:      org.Name,
		Slug:      org.Slug,
		Domain:    org.Domain,
		Hooks:     &hooks,
	}
}
