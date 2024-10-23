package convert

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/timestamppb"

	"go.breu.io/quantm/internal/db/entities"
	authv1 "go.breu.io/quantm/internal/nomad/proto/ctrlplane/auth/v1"
)

const (
	AccountProviderUnknown = "unknown"
	AccountProviderGithub  = "github"
	AccountProviderGoogle  = "google"
)

func ProviderToProto(provider string) authv1.Provider {
	switch provider {
	case AccountProviderGithub:
		return authv1.Provider_PROVIDER_GITHUB
	case AccountProviderGoogle:
		return authv1.Provider_PROVIDER_GOOGLE
	default:
		return authv1.Provider_PROVIDER_UNSPECIFIED
	}
}

func AccountToProto(account *entities.OauthAccount) *authv1.Account {
	return &authv1.Account{
		Id:                UUIDToProto(account.ID),
		CreatedAt:         timestamppb.New(account.CreatedAt),
		UpdatedAt:         timestamppb.New(account.UpdatedAt),
		ExpiresAt:         timestamppb.New(account.ExpiresAt),
		UserId:            UUIDToProto(account.UserID),
		Provider:          ProviderToProto(account.Provider),
		ProviderAccountId: account.ProviderAccountID,
		Kind:              account.Type.String,
	}
}

func ProtoToProvider(proto authv1.Provider) string {
	switch proto {
	case authv1.Provider_PROVIDER_GITHUB:
		return AccountProviderGithub
	case authv1.Provider_PROVIDER_GOOGLE:
		return AccountProviderGoogle
	default:
		return AccountProviderUnknown
	}
}

func ProtoToAccount(proto *authv1.Account) *entities.OauthAccount {
	return &entities.OauthAccount{
		ID:                ProtoToUUID(proto.Id),
		CreatedAt:         proto.CreatedAt.AsTime(),
		UpdatedAt:         proto.UpdatedAt.AsTime(),
		UserID:            ProtoToUUID(proto.UserId),
		Provider:          ProtoToProvider(proto.Provider),
		ProviderAccountID: proto.ProviderAccountId,
		ExpiresAt:         proto.ExpiresAt.AsTime(),
		Type:              pgtype.Text{String: proto.Kind, Valid: true},
	}
}

func ProtoToGetAccountsByUserIDParams(proto *authv1.GetAccountsByUserIDRequest) uuid.UUID {
	return ProtoToUUID(proto.UserId)
}

func ProtoToCreateAccountParams(proto *authv1.CreateAccountRequest) entities.CreateOAuthAccountParams {
	return entities.CreateOAuthAccountParams{
		UserID:            ProtoToUUID(proto.UserId),
		Provider:          ProtoToProvider(proto.Provider),
		ProviderAccountID: proto.ProviderAccountId,
		ExpiresAt:         proto.ExpiresAt.AsTime(),
		Type:              pgtype.Text{String: proto.Kind, Valid: true},
	}
}

func ProtoToGetAccountByIDParams(proto *authv1.GetAccountByIDRequest) uuid.UUID {
	return ProtoToUUID(proto.Id)
}

func ProtoToGetAccountByProviderAccountIDParams(
	proto *authv1.GetAccountByProviderAccountIDRequest,
) entities.GetOAuthAccountByProviderAccountIDParams {
	return entities.GetOAuthAccountByProviderAccountIDParams{
		Provider:          ProtoToProvider(proto.Provider),
		ProviderAccountID: proto.ProviderAccountId,
	}
}
