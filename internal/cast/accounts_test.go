package cast_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"go.breu.io/quantm/internal/cast"
	"go.breu.io/quantm/internal/db/entities"
	authv1 "go.breu.io/quantm/internal/proto/ctrlplane/auth/v1"
)

func TestAccountToProto(t *testing.T) {
	t.Parallel()

	now := time.Now()
	account := &entities.OauthAccount{
		ID:                uuid.New(),
		CreatedAt:         now,
		UpdatedAt:         now,
		ExpiresAt:         now.Add(time.Hour * 24),
		UserID:            uuid.New(),
		Provider:          cast.AuthProviderGoogle,
		ProviderAccountID: "1234567890",
		Type:              "user",
	}

	pb := cast.AccountToProto(account)

	assert.Equal(t, account.ID.String(), pb.GetId())
	assert.True(t, proto.Equal(pb.GetCreatedAt(), timestamppb.New(account.CreatedAt)))
	assert.True(t, proto.Equal(pb.GetUpdatedAt(), timestamppb.New(account.UpdatedAt)))
	assert.True(t, proto.Equal(pb.GetExpiresAt(), timestamppb.New(account.ExpiresAt)))
	assert.Equal(t, account.UserID.String(), pb.GetUserId())
	assert.Equal(t, authv1.AuthProvider_AUTH_PROVIDER_GOOGLE, pb.GetProvider())
	assert.Equal(t, account.ProviderAccountID, pb.GetProviderAccountId())
	assert.Equal(t, "user", pb.GetKind())
}

func TestProtoToAccount(t *testing.T) {
	t.Parallel()

	now := time.Now()
	pb := &authv1.Account{
		Id:                uuid.New().String(),
		CreatedAt:         timestamppb.New(now),
		UpdatedAt:         timestamppb.New(now),
		ExpiresAt:         timestamppb.New(now.Add(time.Hour * 24)),
		UserId:            uuid.New().String(),
		Provider:          authv1.AuthProvider_AUTH_PROVIDER_GOOGLE,
		ProviderAccountId: "1234567890",
		Kind:              "user",
	}

	acc := cast.ProtoToAccount(pb)

	assert.Equal(t, uuid.MustParse(pb.GetId()), acc.ID)
	assert.Equal(t, pb.GetCreatedAt().AsTime(), acc.CreatedAt)
	assert.Equal(t, pb.GetUpdatedAt().AsTime(), acc.UpdatedAt)
	assert.Equal(t, pb.GetExpiresAt().AsTime(), acc.ExpiresAt)
	assert.Equal(t, uuid.MustParse(pb.GetUserId()), acc.UserID)
	assert.Equal(t, cast.AuthProviderGoogle, acc.Provider)
	assert.Equal(t, pb.GetProviderAccountId(), acc.ProviderAccountID)
	assert.Equal(t, "user", acc.Type)
}

func TestProtoToGetAccountsByUserIDParams(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	req := &authv1.GetAccountsByUserIDRequest{UserId: id.String()}

	parsedID := cast.ProtoToGetAccountsByUserIDParams(req)

	assert.Equal(t, id, parsedID)
}

func TestProtoToCreateAccountParams(t *testing.T) {
	t.Parallel()

	now := time.Now()
	req := &authv1.CreateAccountRequest{
		UserId:            uuid.New().String(),
		Provider:          authv1.AuthProvider_AUTH_PROVIDER_GOOGLE,
		ProviderAccountId: "1234567890",
		ExpiresAt:         timestamppb.New(now),
		Kind:              "user",
	}

	params := cast.ProtoToCreateAccountParams(req)

	assert.Equal(t, uuid.MustParse(req.GetUserId()), params.UserID)
	assert.Equal(t, cast.AuthProviderGoogle, params.Provider)
	assert.Equal(t, req.GetProviderAccountId(), params.ProviderAccountID)
	assert.Equal(t, req.GetExpiresAt().AsTime(), params.ExpiresAt)
	assert.Equal(t, "user", params.Type)
}

func TestProtoToGetAccountByIDParams(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	req := &authv1.GetAccountByIDRequest{Id: id.String()}

	parsedID := cast.ProtoToGetAccountByIDParams(req)

	assert.Equal(t, id, parsedID)
}

func TestProtoToGetAccountByProviderAccountIDParams(t *testing.T) {
	t.Parallel()

	req := &authv1.GetAccountByProviderAccountIDRequest{
		Provider:          authv1.AuthProvider_AUTH_PROVIDER_GOOGLE,
		ProviderAccountId: "1234567890",
	}

	params := cast.ProtoToGetAccountByProviderAccountIDParams(req)

	assert.Equal(t, cast.AuthProviderGoogle, params.Provider)
	assert.Equal(t, req.GetProviderAccountId(), params.ProviderAccountID)
}
