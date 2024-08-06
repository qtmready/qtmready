package defs

import (
	"fmt"
)

type (
	providerNotFoundError struct {
		name string
	}

	resourceNotFoundError struct {
		name     string
		provider string
	}

	queueError struct {
		pr   *RepoIOPullRequest
		repo *Repo
		code int
	}
)

func (e *providerNotFoundError) Error() string {
	return fmt.Sprintf("provider %s not found. please register your providers first.", e.name)
}

func (e *resourceNotFoundError) Error() string {
	return fmt.Sprintf("resource %s not found. please register your resource with the provider %s first.", e.name, e.provider)
}

func (e *queueError) Error() string {
	msg := ""

	switch e.code {
	case 10400:
		msg = fmt.Sprintf("unable to schedule pr %d in repo %s", e.pr.Number, e.repo.Name)
	case 10409:
		msg = fmt.Sprintf("pr %d in repo %s is already scheduled", e.pr.Number, e.repo.Name)
	default:
		msg = fmt.Sprintf("unknown error for pr %d in repo %s", e.pr.Number, e.repo.Name)
	}

	return msg
}

// NewProviderNotFoundError creates an error for when a provider is not found.
//
// It takes the name of the provider that wasn't found and returns an error
// that can be used to inform the user about the missing provider.
func NewProviderNotFoundError(name string) error {
	return &providerNotFoundError{name}
}

// NewResourceNotFoundError creates an error for when a resource is not found.
//
// It takes the name of the resource and the provider it should be associated with.
// The returned error can be used to inform the user about the missing resource
// and which provider it should be registered with.
func NewResourceNotFoundError(name string, provider string) error {
	return &resourceNotFoundError{name, provider}
}

// NewQueueSchedulingError creates an error for when a pull request cannot be scheduled.
//
// It takes a RepoIOPullRequest and a Repo, returning an error that indicates
// the pull request could not be scheduled for the given repository.
func NewQueueSchedulingError(pr *RepoIOPullRequest, repo *Repo) error {
	return &queueError{pr, repo, 10400}
}

// NewQueueDuplicatedError creates an error for when a pull request is already scheduled.
//
// It takes a RepoIOPullRequest and a Repo, returning an error that indicates
// the pull request is already scheduled for the given repository.
func NewQueueDuplicatedError(pr *RepoIOPullRequest, repo *Repo) error {
	return &queueError{pr, repo, 10409}
}
