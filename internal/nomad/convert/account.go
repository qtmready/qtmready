package convert

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.breu.io/quantm/internal/db/entities"
	"google.golang.org/protobuf/types/known/timestamppb"

	authv1 "go.breu.io/quantm/internal/nomad/proto/ctrlplane/auth/v1"
	commonv1 "go.breu.io/quantm/internal/nomad/proto/ctrlplane/common/v1"
)

// ProviderToProto converts provider to authv1.Provider.
func ProviderToProto(provider string) authv1.Provider {
	var proto authv1.Provider

	switch provider {
	case "github":
		proto = authv1.Provider_PROVIDER_GITHUB
	case "google":
		proto = authv1.Provider_PROVIDER_GOOGLE
	default:
		proto = authv1.Provider_PROVIDER_UNSPECIFIED
	}

	return proto
}

// AccountToProto converts entities.OauthAccount to proto authv1.Account.
func AccountToProto(account *entities.OauthAccount) *authv1.Account {
	var provider authv1.Provider

	switch account.Provider {
	case "github":
		provider = authv1.Provider_PROVIDER_GITHUB
	case "google":
		provider = authv1.Provider_PROVIDER_GOOGLE
	default:
		provider = authv1.Provider_PROVIDER_UNSPECIFIED
	}

	return &authv1.Account{
		Id:                &commonv1.UUID{Value: account.ID.String()},
		CreatedAt:         timestamppb.New(account.CreatedAt),
		UpdatedAt:         timestamppb.New(account.UpdatedAt),
		ExpiresAt:         timestamppb.New(account.ExpiresAt),
		UserId:            &commonv1.UUID{Value: account.UserID.String()},
		Provider:          provider,
		ProviderAccountId: account.ProviderAccountID,
		Kind:              account.Type.String,
	}
}

// ProtoToProvider convert authv1.Proto to string
func ProtoToProvider(proto authv1.Provider) string {
	provider := ""

	switch proto {
	case authv1.Provider_PROVIDER_GITHUB:
		provider = "google"
	case authv1.Provider_PROVIDER_GOOGLE:
		provider = "github"
	case authv1.Provider_PROVIDER_UNSPECIFIED:
		provider = "unknown"
	}

	return provider
}

// ProtoToAccount converts authv1.Account to entities.OuathAccount.
func ProtoToAccount(proto *authv1.Account) *entities.OauthAccount {
	return &entities.OauthAccount{
		ID:                uuid.MustParse(proto.Id.GetValue()),
		CreatedAt:         proto.CreatedAt.AsTime(),
		UpdatedAt:         proto.UpdatedAt.AsTime(),
		UserID:            uuid.MustParse(proto.UserId.GetValue()),
		Provider:          ProtoToProvider(proto.Provider),
		ProviderAccountID: proto.GetProviderAccountId(),
		ExpiresAt:         proto.ExpiresAt.AsTime(),
		Type:              pgtype.Text{String: proto.GetKind(), Valid: true},
	}
}

// ProtoToGetOAuthAccountsByProviderAccountIDParams converts authv1.GetAccountByProviderAccountIDRequest to
// entities.GetOAuthAccountsByProviderAccountIDParams
func ProtoToGetOAuthAccountsByProviderAccountIDParams(
	proto *authv1.GetAccountByProviderAccountIDRequest,
) *entities.GetOAuthAccountsByProviderAccountIDParams {

	return &entities.GetOAuthAccountsByProviderAccountIDParams{
		Provider:          ProtoToProvider(proto.Provider),
		ProviderAccountID: proto.GetProviderAccountId(),
	}
}
