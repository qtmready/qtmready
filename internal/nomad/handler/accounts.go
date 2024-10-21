package handler

import (
	"context"

	timestamppb "google.golang.org/protobuf/types/known/timestamppb"

	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/erratic"
	authv1 "go.breu.io/quantm/internal/nomad/proto/ctrlplane/auth/v1"
	commonv1 "go.breu.io/quantm/internal/nomad/proto/ctrlplane/common/v1"
)

type (
	AccountService struct {
		authv1.UnimplementedAccountServiceServer
	}
)

func (s *AccountService) GetAccountByProviderAccountID(
	ctx context.Context, req *authv1.GetAccountByProviderAccountIDRequest,
) (*authv1.GetAccountByProviderAccountIDResponse, error) {
	provider := ""

	switch req.Provider {
	case authv1.Provider_PROVIDER_GITHUB:
		provider = "github"
	case authv1.Provider_PROVIDER_GOOGLE:
		provider = "google"
	case authv1.Provider_PROVIDER_UNSPECIFIED:
		return nil, erratic.NewBadRequestError()
	}

	resp, err := db.Queries().GetOAuthAccountsByProviderAccountID(
		ctx,
		entities.GetOAuthAccountsByProviderAccountIDParams{
			ProviderAccountID: req.GetProviderAccountId(),
			Provider:          provider,
		},
	)

	if err != nil {
		return nil, erratic.NewNotFoundError()
	}

	var prv authv1.Provider

	switch resp.Provider {
	case "github":
		prv = authv1.Provider_PROVIDER_GITHUB
	case "google":
		prv = authv1.Provider_PROVIDER_GOOGLE
	default:
		return nil, erratic.NewInternalServerError()
	}

	return &authv1.GetAccountByProviderAccountIDResponse{
		Account: &authv1.Account{
			Id:                &commonv1.UUID{Value: resp.ID.String()},
			CreatedAt:         timestamppb.New(resp.CreatedAt),
			UpdatedAt:         timestamppb.New(resp.UpdatedAt),
			ExpiresAt:         timestamppb.New(resp.ExpiresAt),
			UserId:            &commonv1.UUID{Value: resp.UserID.String()},
			Provider:          prv,
			ProviderAccountId: resp.ProviderAccountID,
			Kind:              resp.Type.String,
		},
	}, nil
}
