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

// ProviderToProto converts an account provider string to its protobuf representation.
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

// AccountToProto converts an OauthAccount entity to its protobuf representation.
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

// ProtoToProvider converts a protobuf provider to its string representation.
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

// ProtoToAccount converts a protobuf account to its OauthAccount entity representation.
func ProtoToAccount(proto *authv1.Account) *entities.OauthAccount {
	return &entities.OauthAccount{
		ID:                ProtoToUUID(proto.GetId()),
		CreatedAt:         proto.GetCreatedAt().AsTime(),
		UpdatedAt:         proto.GetUpdatedAt().AsTime(),
		UserID:            ProtoToUUID(proto.GetUserId()),
		Provider:          ProtoToProvider(proto.GetProvider()),
		ProviderAccountID: proto.GetProviderAccountId(),
		ExpiresAt:         proto.GetExpiresAt().AsTime(),
		Type:              pgtype.Text{String: proto.GetKind(), Valid: true},
	}
}

// ProtoToGetAccountsByUserIDParams converts a protobuf GetAccountsByUserIDRequest to a UUID.
func ProtoToGetAccountsByUserIDParams(proto *authv1.GetAccountsByUserIDRequest) uuid.UUID {
	return ProtoToUUID(proto.GetUserId())
}

// ProtoToCreateAccountParams converts a protobuf CreateAccountRequest to an entities.CreateOAuthAccountParams.
func ProtoToCreateAccountParams(proto *authv1.CreateAccountRequest) entities.CreateOAuthAccountParams {
	return entities.CreateOAuthAccountParams{
		UserID:            ProtoToUUID(proto.GetUserId()),
		Provider:          ProtoToProvider(proto.GetProvider()),
		ProviderAccountID: proto.GetProviderAccountId(),
		ExpiresAt:         proto.GetExpiresAt().AsTime(),
		Type:              pgtype.Text{String: proto.GetKind(), Valid: true},
	}
}

// ProtoToGetAccountByIDParams converts a protobuf GetAccountByIDRequest to a UUID.
func ProtoToGetAccountByIDParams(proto *authv1.GetAccountByIDRequest) uuid.UUID {
	return ProtoToUUID(proto.GetId())
}

// ProtoToGetAccountByProviderAccountIDParams converts a protobuf GetAccountByProviderAccountIDRequest to an
// entities.GetOAuthAccountByProviderAccountIDParams.
func ProtoToGetAccountByProviderAccountIDParams(
	proto *authv1.GetAccountByProviderAccountIDRequest,
) entities.GetOAuthAccountByProviderAccountIDParams {
	return entities.GetOAuthAccountByProviderAccountIDParams{
		Provider:          ProtoToProvider(proto.GetProvider()),
		ProviderAccountID: proto.GetProviderAccountId(),
	}
}
