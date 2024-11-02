package authnmd

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"go.breu.io/quantm/internal/cast"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/erratic"
	authv1 "go.breu.io/quantm/internal/proto/ctrlplane/auth/v1"
	"go.breu.io/quantm/internal/proto/ctrlplane/auth/v1/authv1connect"
)

type (
	AccountService struct {
		authv1connect.UnimplementedAccountServiceHandler
	}
)

func (s *AccountService) GetAccountByProviderAccountID(
	ctx context.Context,
	req *connect.Request[authv1.GetAccountByProviderAccountIDRequest],
) (*connect.Response[authv1.GetAccountByProviderAccountIDResponse], error) {
	params := cast.ProtoToGetAccountByProviderAccountIDParams(req.Msg)

	account, err := db.Queries().GetOAuthAccountByProviderAccountID(ctx, params)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, erratic.NewNotFoundError("entity", "accounts", "provider_id", req.Msg.GetProviderAccountId()).ToConnectError()
		}

		return nil, erratic.NewInternalServerError().DataBaseError(err).ToConnectError()
	}

	proto := &authv1.GetAccountByProviderAccountIDResponse{Account: cast.AccountToProto(&account)}

	return connect.NewResponse(proto), nil
}

func (s *AccountService) GetAccountsByUserID(
	ctx context.Context,
	req *connect.Request[authv1.GetAccountsByUserIDRequest],
) (*connect.Response[authv1.GetAccountsByUserIDResponse], error) {
	id := uuid.MustParse(req.Msg.GetUserId())

	accounts, err := db.Queries().GetOAuthAccountsByUserID(ctx, id)
	if err != nil {
		return nil, erratic.NewInternalServerError().DataBaseError(err).ToProto().Err()
	}

	proto := make([]*authv1.Account, len(accounts))
	for i, account := range accounts {
		proto[i] = cast.AccountToProto(&account)
	}

	return connect.NewResponse(&authv1.GetAccountsByUserIDResponse{Accounts: proto}), nil
}

func (s *AccountService) CreateAccount(
	ctx context.Context,
	req *connect.Request[authv1.CreateAccountRequest],
) (*connect.Response[authv1.CreateAccountResponse], error) {
	params := cast.ProtoToCreateAccountParams(req.Msg)

	account, err := db.Queries().CreateOAuthAccount(ctx, params)
	if err != nil {
		return nil, erratic.NewInternalServerError().DataBaseError(err).ToProto().Err()
	}

	return connect.NewResponse(&authv1.CreateAccountResponse{Account: cast.AccountToProto(&account)}), nil
}

func (s *AccountService) GetAccountByID(
	ctx context.Context,
	req *connect.Request[authv1.GetAccountByIDRequest],
) (*connect.Response[authv1.GetAccountByIDResponse], error) {
	id := uuid.MustParse(req.Msg.GetId())

	account, err := db.Queries().GetOAuthAccountByID(ctx, id)
	if err != nil {
		return nil, erratic.NewNotFoundError("entity", "accounts", "id", req.Msg.GetId()).ToProto().Err()
	}

	return connect.NewResponse(&authv1.GetAccountByIDResponse{Account: cast.AccountToProto(&account)}), nil
}

func NewAccountSericeServiceHandler() (string, http.Handler) {
	return authv1connect.NewAccountServiceHandler(&AccountService{})
}
