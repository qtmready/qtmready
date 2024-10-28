package convert

import (
	commonv1 "go.breu.io/quantm/internal/nomad/proto/ctrlplane/common/v1"
)

const (
	AuthProviderUnknown = "unknown"
	AuthProviderGithub  = "github"
	AuthProviderGoogle  = "google"
)

func AuthProviderToProto(provider string) commonv1.AuthProvider {
	switch provider {
	case AuthProviderGithub:
		return commonv1.AuthProvider_AUTH_PROVIDER_GITHUB
	case AuthProviderGoogle:
		return commonv1.AuthProvider_AUTH_PROVIDER_GOOGLE
	default:
		return commonv1.AuthProvider_AUTH_PROVIDER_UNSPECIFIED
	}
}

func ProtoToAuthProvider(proto commonv1.AuthProvider) string {
	switch proto {
	case commonv1.AuthProvider_AUTH_PROVIDER_GITHUB:
		return AuthProviderGithub
	case commonv1.AuthProvider_AUTH_PROVIDER_GOOGLE:
		return AuthProviderGoogle
	case commonv1.AuthProvider_AUTH_PROVIDER_UNSPECIFIED:
		return AuthProviderUnknown
	default:
		return AuthProviderUnknown
	}
}
