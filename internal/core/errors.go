// Copyright © 2023, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRIBUTING ANY OF
// THE SOFTWARE, YOU AGREE TO THE TERMS OF THE LICENSE AGREEMENT.
//
// The above copyright notice and the subsequent license agreement shall be included in all copies or substantial
// portions of the software.
//
// Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, IMPLIED, STATUTORY, OR OTHERWISE, AND
// SPECIFICALLY DISCLAIMS ANY WARRANTY OF MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE
// SOFTWARE.
//
// Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT LIMITED TO, LOST PROFITS OR ANY
// CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// ARISING OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY APPLICABLE LAW.

package core

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

func NewProviderNotFoundError(name string) error {
	return &providerNotFoundError{name}
}

func NewResourceNotFoundError(name string, provider string) error {
	return &resourceNotFoundError{name, provider}
}

func NewQueueError(pr *RepoIOPullRequest, repo *Repo, code int) error {
	return &queueError{pr, repo, code}
}

func NewQueueSchedulingError(pr *RepoIOPullRequest, repo *Repo) error {
	return &queueError{pr, repo, 10400}
}

func NewQueueDuplicatedError(pr *RepoIOPullRequest, repo *Repo) error {
	return &queueError{pr, repo, 10409}
}
