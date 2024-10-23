package handler

import (
	"context"

	emptypb "google.golang.org/protobuf/types/known/emptypb"

	healthzv1 "go.breu.io/quantm/internal/nomad/proto/ctrlplane/healthz/v1"
)

type (
	HealthCheckService struct {
		healthzv1.UnimplementedHealthCheckServiceServer
	}
)

func (s *HealthCheckService) Status(ctx context.Context, _ *emptypb.Empty) (*healthzv1.HealthCheckResponse, error) {
	return &healthzv1.HealthCheckResponse{
		Database: true,
		Temporal: true,
	}, nil
}
