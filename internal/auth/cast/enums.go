package cast

import (
	authv1 "go.breu.io/quantm/internal/proto/ctrlplane/auth/v1"
)

const (
	AuthProviderUnknown = "unknown"
	AuthProviderGithub  = "github"
	AuthProviderGoogle  = "google"
)

func AuthProviderToProto(provider string) authv1.AuthProvider {
	switch provider {
	case AuthProviderGithub:
		return authv1.AuthProvider_AUTH_PROVIDER_GITHUB
	case AuthProviderGoogle:
		return authv1.AuthProvider_AUTH_PROVIDER_GOOGLE
	default:
		return authv1.AuthProvider_AUTH_PROVIDER_UNSPECIFIED
	}
}

func ProtoToAuthProvider(proto authv1.AuthProvider) string {
	switch proto {
	case authv1.AuthProvider_AUTH_PROVIDER_GITHUB:
		return AuthProviderGithub
	case authv1.AuthProvider_AUTH_PROVIDER_GOOGLE:
		return AuthProviderGoogle
	case authv1.AuthProvider_AUTH_PROVIDER_UNSPECIFIED:
		return AuthProviderUnknown
	default:
		return AuthProviderUnknown
	}
}
