package cast

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"go.breu.io/quantm/internal/db/entities"
	authv1 "go.breu.io/quantm/internal/proto/ctrlplane/auth/v1"
)

func OrgToProto(org *entities.Org) *authv1.Org {
	return &authv1.Org{
		Id:        org.ID.String(),
		CreatedAt: timestamppb.New(org.CreatedAt),
		UpdatedAt: timestamppb.New(org.UpdatedAt),
		Name:      org.Name,
		Slug:      org.Slug,
	}
}
