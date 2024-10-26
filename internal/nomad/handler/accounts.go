package handler

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5"

	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/erratic"
	"go.breu.io/quantm/internal/nomad/convert"
	authv1 "go.breu.io/quantm/internal/nomad/proto/ctrlplane/auth/v1"
	"go.breu.io/quantm/internal/nomad/proto/ctrlplane/auth/v1/authv1connect"
)

type (
	AccountService struct{}
)

func (s *AccountService) GetAccountByProviderAccountID(
	ctx context.Context,
	rqst *connect.Request[authv1.GetAccountByProviderAccountIDRequest],
) (*connect.Response[authv1.GetAccountByProviderAccountIDResponse], error) {
	params := convert.ProtoToGetAccountByProviderAccountIDParams(rqst.Msg)

	account, err := db.Queries().GetOAuthAccountByProviderAccountID(ctx, params)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, erratic.NewNotFoundError(
				"entity", "accounts",
				"provider_account_id", rqst.Msg.GetProviderAccountId(),
			).ToConnectError()
		}

		return nil, erratic.NewInternalServerError().DataBaseError(err).ToConnectError()
	}

	proto := &authv1.GetAccountByProviderAccountIDResponse{Account: convert.AccountToProto(&account)}

	return connect.NewResponse(proto), nil
}

func (s *AccountService) GetAccountsByUserID(
	ctx context.Context,
	req *connect.Request[authv1.GetAccountsByUserIDRequest],
) (*connect.Response[authv1.GetAccountsByUserIDResponse], error) {
	id := convert.ProtoToUUID(req.Msg.UserId)

	accounts, err := db.Queries().GetOAuthAccountsByUserID(ctx, id)
	if err != nil {
		return nil, erratic.NewInternalServerError().DataBaseError(err).ToProto().Err()
	}

	proto := make([]*authv1.Account, len(accounts))
	for i, account := range accounts {
		proto[i] = convert.AccountToProto(&account)
	}

	return connect.NewResponse(&authv1.GetAccountsByUserIDResponse{Accounts: proto}), nil
}

func (s *AccountService) CreateAccount(
	ctx context.Context,
	req *connect.Request[authv1.CreateAccountRequest],
) (*connect.Response[authv1.CreateAccountResponse], error) {
	params := convert.ProtoToCreateAccountParams(req.Msg)

	account, err := db.Queries().CreateOAuthAccount(ctx, params)
	if err != nil {
		return nil, erratic.NewInternalServerError().DataBaseError(err).ToProto().Err()
	}

	return connect.NewResponse(&authv1.CreateAccountResponse{Account: convert.AccountToProto(&account)}), nil
}

func (s *AccountService) GetAccountByID(
	ctx context.Context,
	req *connect.Request[authv1.GetAccountByIDRequest],
) (*connect.Response[authv1.GetAccountByIDResponse], error) {
	id := convert.ProtoToUUID(req.Msg.Id)

	account, err := db.Queries().GetOAuthAccountByID(ctx, id)
	if err != nil {
		return nil, erratic.NewNotFoundError("entity", "accounts", "id", req.Msg.GetId().Value).ToProto().Err()
	}

	return connect.NewResponse(&authv1.GetAccountByIDResponse{Account: convert.AccountToProto(&account)}), nil
}

func NewAccountSericeServiceHandler() (string, http.Handler) {
	return authv1connect.NewAccountServiceHandler(&AccountService{})
}
