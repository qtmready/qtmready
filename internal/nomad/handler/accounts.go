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
	ctx context.Context,
	rqst *authv1.GetAccountByProviderAccountIDRequest,
) (*authv1.GetAccountByProviderAccountIDResponse, error) {
	params := convert.ProtoToGetAccountByProviderAccountIDParams(rqst)

	account, err := db.Queries().GetOAuthAccountByProviderAccountID(ctx, params)
	if err != nil {
		return nil, erratic.NewNotFoundError(
			"entity", "accounts",
			"provider_account_id", rqst.GetProviderAccountId(),
		).ToProto()
	}

	proto := &authv1.GetAccountByProviderAccountIDResponse{Account: convert.AccountToProto(&account)}

	return proto, nil
}

func (s *AccountService) GetAccountsByUserID(
	ctx context.Context, req *authv1.GetAccountsByUserIDRequest,
) (*authv1.GetAccountsByUserIDResponse, error) {
	id := convert.ProtoToUUID(req.UserId)
	accounts, err := db.Queries().GetOAuthAccountsByUserID(ctx, id)
	if err != nil {
		return nil, erratic.NewInternalServerError().DataBaseError(err).ToProto()
	}

	proto := make([]*authv1.Account, len(accounts))
	for i, account := range accounts {
		proto[i] = convert.AccountToProto(&account)
	}

	return &authv1.GetAccountsByUserIDResponse{Accounts: proto}, nil
}

func (s *AccountService) CreateAccount(ctx context.Context, req *authv1.CreateAccountRequest) (*authv1.CreateAccountResponse, error) {
	params := convert.ProtoToCreateAccountParams(req)
	account, err := db.Queries().CreateOAuthAccount(ctx, params)
	if err != nil {
		return nil, erratic.NewInternalServerError().DataBaseError(err).ToProto()
	}

	return &authv1.CreateAccountResponse{Account: convert.AccountToProto(&account)}, nil
}

func (s *AccountService) GetAccountByID(ctx context.Context, req *authv1.GetAccountByIDRequest) (*authv1.GetAccountByIDResponse, error) {
	id := convert.ProtoToUUID(req.Id)
	account, err := db.Queries().GetOAuthAccountByID(ctx, id)
	if err != nil {
		return nil, erratic.NewNotFoundError("entity", "accounts", "id", req.GetId().Value).ToProto()
	}

	return &authv1.GetAccountByIDResponse{Account: convert.AccountToProto(&account)}, nil
}
