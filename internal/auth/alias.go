package auth

import (
	"go.breu.io/quantm/internal/auth/config"
	"go.breu.io/quantm/internal/auth/nomad"
)

var (
	Secret    = config.Secret
	SetSecret = config.SetSecret
)

var (
	NomadAuthContext           = nomad.GetAuthContext
	NomadInterceptor           = nomad.AuthInterceptor
	NomadAccountServiceHandler = nomad.NewAccountSericeServiceHandler
	NomadOrgServiceHandler     = nomad.NewOrgServiceServiceHandler
	NomadUserServiceHandler    = nomad.NewUserSericeServiceHandler
)
