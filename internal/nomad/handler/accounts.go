package handler

import (
	"context"

	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/erratic"
	"go.breu.io/quantm/internal/nomad/convert"
	authv1 "go.breu.io/quantm/internal/nomad/proto/ctrlplane/auth/v1"
)

type (
	AccountService struct {
		authv1.UnimplementedAccountServiceServer
	}
)

func (s *AccountService) GetAccountByProviderAccountID(
	ctx context.Context, rqst *authv1.GetAccountByProviderAccountIDRequest,
) (*authv1.GetAccountByProviderAccountIDResponse, error) {
	account, err := db.Queries().GetOAuthAccountsByProviderAccountID(
		ctx, *convert.ProtoToGetOAuthAccountsByProviderAccountIDParams(rqst),
	)

	if err != nil {
		return nil, erratic.NewNotFoundError("account_id", rqst.GetProviderAccountId())
	}

	return &authv1.GetAccountByProviderAccountIDResponse{Account: convert.AccountToProto(&account)}, nil
}
