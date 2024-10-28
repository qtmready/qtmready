// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: ctrlplane/healthz/v1/healthz.proto

package healthzv1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	v1 "go.breu.io/quantm/internal/nomad/proto/ctrlplane/healthz/v1"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
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
	// HealthCheckServiceName is the fully-qualified name of the HealthCheckService service.
	HealthCheckServiceName = "ctrlplane.healthz.v1.HealthCheckService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// HealthCheckServiceStatusProcedure is the fully-qualified name of the HealthCheckService's Status
	// RPC.
	HealthCheckServiceStatusProcedure = "/ctrlplane.healthz.v1.HealthCheckService/Status"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	healthCheckServiceServiceDescriptor      = v1.File_ctrlplane_healthz_v1_healthz_proto.Services().ByName("HealthCheckService")
	healthCheckServiceStatusMethodDescriptor = healthCheckServiceServiceDescriptor.Methods().ByName("Status")
)

// HealthCheckServiceClient is a client for the ctrlplane.healthz.v1.HealthCheckService service.
type HealthCheckServiceClient interface {
	// buf:lint:ignore RPCNamingConventions
	Status(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[v1.StatusResponse], error)
}

// NewHealthCheckServiceClient constructs a client for the ctrlplane.healthz.v1.HealthCheckService
// service. By default, it uses the Connect protocol with the binary Protobuf Codec, asks for
// gzipped responses, and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply
// the connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewHealthCheckServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) HealthCheckServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &healthCheckServiceClient{
		status: connect.NewClient[emptypb.Empty, v1.StatusResponse](
			httpClient,
			baseURL+HealthCheckServiceStatusProcedure,
			connect.WithSchema(healthCheckServiceStatusMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// healthCheckServiceClient implements HealthCheckServiceClient.
type healthCheckServiceClient struct {
	status *connect.Client[emptypb.Empty, v1.StatusResponse]
}

// Status calls ctrlplane.healthz.v1.HealthCheckService.Status.
func (c *healthCheckServiceClient) Status(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[v1.StatusResponse], error) {
	return c.status.CallUnary(ctx, req)
}

// HealthCheckServiceHandler is an implementation of the ctrlplane.healthz.v1.HealthCheckService
// service.
type HealthCheckServiceHandler interface {
	// buf:lint:ignore RPCNamingConventions
	Status(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[v1.StatusResponse], error)
}

// NewHealthCheckServiceHandler builds an HTTP handler from the service implementation. It returns
// the path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewHealthCheckServiceHandler(svc HealthCheckServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	healthCheckServiceStatusHandler := connect.NewUnaryHandler(
		HealthCheckServiceStatusProcedure,
		svc.Status,
		connect.WithSchema(healthCheckServiceStatusMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/ctrlplane.healthz.v1.HealthCheckService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case HealthCheckServiceStatusProcedure:
			healthCheckServiceStatusHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedHealthCheckServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedHealthCheckServiceHandler struct{}

func (UnimplementedHealthCheckServiceHandler) Status(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[v1.StatusResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("ctrlplane.healthz.v1.HealthCheckService.Status is not implemented"))
}
