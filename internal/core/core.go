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
	"context"
	"sync"

	"go.breu.io/quantm/internal/shared"
)

var (
	instance Core
	once     sync.Once
)

type (
	// Core is the interface that defines the core of the application. It is the main entry point for the application.
	// It is responsible for registering different providers and exposing them to the rest of the application.
	//
	// NOTE: This is not an ideal design, because it only registers activities for the providers. It does not register
	// workflows. We may need to revisit this design in the future.
	Core interface {
		RegisterRepoProvider(RepoProvider, RepoProviderActivities)
		RegisterCloudProvider(CloudProvider, CloudProviderActivities)

		RepoProvider(RepoProvider) RepoProviderActivities
		CloudProvider(CloudProvider) CloudProviderActivities
	}

	CoreOption func(Core)

	RepoProviderActivities interface {
		GetLatestCommit(context.Context, string, string) (string, error)
	}

	CloudProviderActivities interface {
		FillMe()
	}

	ProviderActivities struct {
		repos map[RepoProvider]RepoProviderActivities
		cloud map[CloudProvider]CloudProviderActivities
	}

	core struct {
		activity ProviderActivities
	}
)

func (c *core) RegisterRepoProvider(provider RepoProvider, activities RepoProviderActivities) {
	c.activity.repos[provider] = activities
}

func (c *core) RegisterCloudProvider(provider CloudProvider, activities CloudProviderActivities) {
	c.activity.cloud[provider] = activities
}

func (c *core) RepoProvider(name RepoProvider) RepoProviderActivities {
	if p, ok := c.activity.repos[name]; ok {
		return p
	}

	panic(NewProviderNotFoundError(name.String()))
}

func (c *core) CloudProvider(name CloudProvider) CloudProviderActivities {
	if p, ok := c.activity.cloud[name]; ok {
		return p
	}

	panic(NewProviderNotFoundError(name.String()))
}

// WithRepoProvider registers a repo provider with the core.
func WithRepoProvider(name RepoProvider, provider RepoProviderActivities) CoreOption {
	return func(c Core) {
		shared.Logger().Info("core: registering repo provider", "name", name.String())
		c.RegisterRepoProvider(name, provider)
	}
}

// WithCloudProvider registers a cloud provider with the core.
func WithCloudProvider(name CloudProvider, provider CloudProviderActivities) CoreOption {
	return func(c Core) {
		shared.Logger().Info("core: registering cloud provider", "name", name.String())
		c.RegisterCloudProvider(name, provider)
	}
}

// Instance returns a singleton instance of the core. It is best to call this function in the main() function to
// register the providers available to the service. This is because the core uses workflow and activity implementations
// to access the providers.
func Instance(opts ...CoreOption) Core {
	if instance == nil {
		shared.Logger().Info("core: instance not initialized, initializing now ...")
		once.Do(func() {
			instance = &core{
				activity: ProviderActivities{
					repos: make(map[RepoProvider]RepoProviderActivities),
					cloud: make(map[CloudProvider]CloudProviderActivities),
				},
			}

			for _, opt := range opts {
				opt(instance)
			}
		})
	}

	return instance
}
