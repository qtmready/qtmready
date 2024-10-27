// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: ctrlplane/auth/v1/accounts.proto

package authv1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	v1 "go.breu.io/quantm/internal/nomad/proto/ctrlplane/auth/v1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_13_0

const (
	// AccountServiceName is the fully-qualified name of the AccountService service.
	AccountServiceName = "ctrlplane.auth.v1.AccountService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// AccountServiceGetAccountByProviderAccountIDProcedure is the fully-qualified name of the
	// AccountService's GetAccountByProviderAccountID RPC.
	AccountServiceGetAccountByProviderAccountIDProcedure = "/ctrlplane.auth.v1.AccountService/GetAccountByProviderAccountID"
	// AccountServiceGetAccountsByUserIDProcedure is the fully-qualified name of the AccountService's
	// GetAccountsByUserID RPC.
	AccountServiceGetAccountsByUserIDProcedure = "/ctrlplane.auth.v1.AccountService/GetAccountsByUserID"
	// AccountServiceCreateAccountProcedure is the fully-qualified name of the AccountService's
	// CreateAccount RPC.
	AccountServiceCreateAccountProcedure = "/ctrlplane.auth.v1.AccountService/CreateAccount"
	// AccountServiceGetAccountByIDProcedure is the fully-qualified name of the AccountService's
	// GetAccountByID RPC.
	AccountServiceGetAccountByIDProcedure = "/ctrlplane.auth.v1.AccountService/GetAccountByID"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	accountServiceServiceDescriptor                             = v1.File_ctrlplane_auth_v1_accounts_proto.Services().ByName("AccountService")
	accountServiceGetAccountByProviderAccountIDMethodDescriptor = accountServiceServiceDescriptor.Methods().ByName("GetAccountByProviderAccountID")
	accountServiceGetAccountsByUserIDMethodDescriptor           = accountServiceServiceDescriptor.Methods().ByName("GetAccountsByUserID")
	accountServiceCreateAccountMethodDescriptor                 = accountServiceServiceDescriptor.Methods().ByName("CreateAccount")
	accountServiceGetAccountByIDMethodDescriptor                = accountServiceServiceDescriptor.Methods().ByName("GetAccountByID")
)

// AccountServiceClient is a client for the ctrlplane.auth.v1.AccountService service.
type AccountServiceClient interface {
	// Retrieves an account by its provider and identifier.
	GetAccountByProviderAccountID(context.Context, *connect.Request[v1.GetAccountByProviderAccountIDRequest]) (*connect.Response[v1.GetAccountByProviderAccountIDResponse], error)
	// Retrieves accounts associated with a specific user.
	GetAccountsByUserID(context.Context, *connect.Request[v1.GetAccountsByUserIDRequest]) (*connect.Response[v1.GetAccountsByUserIDResponse], error)
	// Creates a new external account.
	CreateAccount(context.Context, *connect.Request[v1.CreateAccountRequest]) (*connect.Response[v1.CreateAccountResponse], error)
	// Retrieves an account by its unique identifier.
	GetAccountByID(context.Context, *connect.Request[v1.GetAccountByIDRequest]) (*connect.Response[v1.GetAccountByIDResponse], error)
}

// NewAccountServiceClient constructs a client for the ctrlplane.auth.v1.AccountService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewAccountServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) AccountServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &accountServiceClient{
		getAccountByProviderAccountID: connect.NewClient[v1.GetAccountByProviderAccountIDRequest, v1.GetAccountByProviderAccountIDResponse](
			httpClient,
			baseURL+AccountServiceGetAccountByProviderAccountIDProcedure,
			connect.WithSchema(accountServiceGetAccountByProviderAccountIDMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		getAccountsByUserID: connect.NewClient[v1.GetAccountsByUserIDRequest, v1.GetAccountsByUserIDResponse](
			httpClient,
			baseURL+AccountServiceGetAccountsByUserIDProcedure,
			connect.WithSchema(accountServiceGetAccountsByUserIDMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		createAccount: connect.NewClient[v1.CreateAccountRequest, v1.CreateAccountResponse](
			httpClient,
			baseURL+AccountServiceCreateAccountProcedure,
			connect.WithSchema(accountServiceCreateAccountMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		getAccountByID: connect.NewClient[v1.GetAccountByIDRequest, v1.GetAccountByIDResponse](
			httpClient,
			baseURL+AccountServiceGetAccountByIDProcedure,
			connect.WithSchema(accountServiceGetAccountByIDMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// accountServiceClient implements AccountServiceClient.
type accountServiceClient struct {
	getAccountByProviderAccountID *connect.Client[v1.GetAccountByProviderAccountIDRequest, v1.GetAccountByProviderAccountIDResponse]
	getAccountsByUserID           *connect.Client[v1.GetAccountsByUserIDRequest, v1.GetAccountsByUserIDResponse]
	createAccount                 *connect.Client[v1.CreateAccountRequest, v1.CreateAccountResponse]
	getAccountByID                *connect.Client[v1.GetAccountByIDRequest, v1.GetAccountByIDResponse]
}

// GetAccountByProviderAccountID calls
// ctrlplane.auth.v1.AccountService.GetAccountByProviderAccountID.
func (c *accountServiceClient) GetAccountByProviderAccountID(ctx context.Context, req *connect.Request[v1.GetAccountByProviderAccountIDRequest]) (*connect.Response[v1.GetAccountByProviderAccountIDResponse], error) {
	return c.getAccountByProviderAccountID.CallUnary(ctx, req)
}

// GetAccountsByUserID calls ctrlplane.auth.v1.AccountService.GetAccountsByUserID.
func (c *accountServiceClient) GetAccountsByUserID(ctx context.Context, req *connect.Request[v1.GetAccountsByUserIDRequest]) (*connect.Response[v1.GetAccountsByUserIDResponse], error) {
	return c.getAccountsByUserID.CallUnary(ctx, req)
}

// CreateAccount calls ctrlplane.auth.v1.AccountService.CreateAccount.
func (c *accountServiceClient) CreateAccount(ctx context.Context, req *connect.Request[v1.CreateAccountRequest]) (*connect.Response[v1.CreateAccountResponse], error) {
	return c.createAccount.CallUnary(ctx, req)
}

// GetAccountByID calls ctrlplane.auth.v1.AccountService.GetAccountByID.
func (c *accountServiceClient) GetAccountByID(ctx context.Context, req *connect.Request[v1.GetAccountByIDRequest]) (*connect.Response[v1.GetAccountByIDResponse], error) {
	return c.getAccountByID.CallUnary(ctx, req)
}

// AccountServiceHandler is an implementation of the ctrlplane.auth.v1.AccountService service.
type AccountServiceHandler interface {
	// Retrieves an account by its provider and identifier.
	GetAccountByProviderAccountID(context.Context, *connect.Request[v1.GetAccountByProviderAccountIDRequest]) (*connect.Response[v1.GetAccountByProviderAccountIDResponse], error)
	// Retrieves accounts associated with a specific user.
	GetAccountsByUserID(context.Context, *connect.Request[v1.GetAccountsByUserIDRequest]) (*connect.Response[v1.GetAccountsByUserIDResponse], error)
	// Creates a new external account.
	CreateAccount(context.Context, *connect.Request[v1.CreateAccountRequest]) (*connect.Response[v1.CreateAccountResponse], error)
	// Retrieves an account by its unique identifier.
	GetAccountByID(context.Context, *connect.Request[v1.GetAccountByIDRequest]) (*connect.Response[v1.GetAccountByIDResponse], error)
}

// NewAccountServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewAccountServiceHandler(svc AccountServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	accountServiceGetAccountByProviderAccountIDHandler := connect.NewUnaryHandler(
		AccountServiceGetAccountByProviderAccountIDProcedure,
		svc.GetAccountByProviderAccountID,
		connect.WithSchema(accountServiceGetAccountByProviderAccountIDMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	accountServiceGetAccountsByUserIDHandler := connect.NewUnaryHandler(
		AccountServiceGetAccountsByUserIDProcedure,
		svc.GetAccountsByUserID,
		connect.WithSchema(accountServiceGetAccountsByUserIDMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	accountServiceCreateAccountHandler := connect.NewUnaryHandler(
		AccountServiceCreateAccountProcedure,
		svc.CreateAccount,
		connect.WithSchema(accountServiceCreateAccountMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	accountServiceGetAccountByIDHandler := connect.NewUnaryHandler(
		AccountServiceGetAccountByIDProcedure,
		svc.GetAccountByID,
		connect.WithSchema(accountServiceGetAccountByIDMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/ctrlplane.auth.v1.AccountService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case AccountServiceGetAccountByProviderAccountIDProcedure:
			accountServiceGetAccountByProviderAccountIDHandler.ServeHTTP(w, r)
		case AccountServiceGetAccountsByUserIDProcedure:
			accountServiceGetAccountsByUserIDHandler.ServeHTTP(w, r)
		case AccountServiceCreateAccountProcedure:
			accountServiceCreateAccountHandler.ServeHTTP(w, r)
		case AccountServiceGetAccountByIDProcedure:
			accountServiceGetAccountByIDHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedAccountServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedAccountServiceHandler struct{}

func (UnimplementedAccountServiceHandler) GetAccountByProviderAccountID(context.Context, *connect.Request[v1.GetAccountByProviderAccountIDRequest]) (*connect.Response[v1.GetAccountByProviderAccountIDResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("ctrlplane.auth.v1.AccountService.GetAccountByProviderAccountID is not implemented"))
}

func (UnimplementedAccountServiceHandler) GetAccountsByUserID(context.Context, *connect.Request[v1.GetAccountsByUserIDRequest]) (*connect.Response[v1.GetAccountsByUserIDResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("ctrlplane.auth.v1.AccountService.GetAccountsByUserID is not implemented"))
}

func (UnimplementedAccountServiceHandler) CreateAccount(context.Context, *connect.Request[v1.CreateAccountRequest]) (*connect.Response[v1.CreateAccountResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("ctrlplane.auth.v1.AccountService.CreateAccount is not implemented"))
}

func (UnimplementedAccountServiceHandler) GetAccountByID(context.Context, *connect.Request[v1.GetAccountByIDRequest]) (*connect.Response[v1.GetAccountByIDResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("ctrlplane.auth.v1.AccountService.GetAccountByID is not implemented"))
}
