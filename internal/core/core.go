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

	"github.com/gocql/gocql"
	"go.temporal.io/sdk/workflow"

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
	// NOTE: This is not an ideal design, because it only registers providers for the providers. It does not register
	// workflows. We may need to revisit this design in the future.
	Core interface {
		RegisterRepoProvider(RepoProvider, RepoIO)
		ResgisterMessageProvider(MessageProvider, MessageIO)
		RegisterCloudProvider(CloudProvider, CloudIO)
		RegisterCloudResource(provider CloudProvider, driver Driver, resource ResourceConstructor)

		RepoIO(RepoProvider) RepoIO
		CloudProvider(CloudProvider) CloudIO
		ResourceConstructor(CloudProvider, Driver) ResourceConstructor
		MessageProvider(MessageProvider) MessageIO
	}

	Option func(Core)

	// CloudIO is the interface that defines the operations that can be performed on a cloud provider.
	CloudIO interface {
		FillMe()
	}

	// MessageIO is the interface that defines the operations that can be performed on a message provider.
	MessageIO interface {
		SendStaleBranchMessage(ctx context.Context, teamID string, stale *LatestCommit) error
		SendNumberOfLinesExceedMessage(ctx context.Context, teamID, repoName, branchName string, threshold int,
			branchChnages *BranchChanges) error
		SendMergeConflictsMessage(ctx context.Context, teamID string, merge *LatestCommit) error
	}

	// Providers is a struct that holds the different providers that are registered with the core.
	Providers struct {
		repos   map[RepoProvider]RepoIO
		cloud   map[CloudProvider]CloudIO
		message map[MessageProvider]MessageIO
	}

	// CloudResource is the interface that defines the operations that can be performed on a cloud resource.
	CloudResource interface {
		Provision(ctx workflow.Context) (workflow.Future, error)
		DeProvision() error
		Deploy(workflow.Context, []Workload, gocql.UUID) error
		UpdateTraffic(workflow.Context, int32) error
		Marshal() ([]byte, error)
	}

	// ResourceConstructor is the interface that defines the operations that can be performed on a cloud resource constructor.
	ResourceConstructor interface {
		Create(name string, region string, config string, providerConfig string) (CloudResource, error)
		CreateFromJson(data []byte) CloudResource
	}

	core struct {
		providers Providers
		resources map[CloudProvider]map[Driver]ResourceConstructor
		once      sync.Once // Do we really need this?
	}
)

// RegisterCloudResource the cloud resource constructor for against a cloud resource.
//
// All cloud resources workflows can be registered like this but the
// problem with this approach is that the signature of deploy workflow needs to be generic and part of cloud resource interface
// e.g DeployWorkflow(ctx workflow.Context, resource CloudResource)
// but the above won't work without custom data converter as CloudResource is an interface and temporal will not be able to unmarshal it
// so to solve this problem we have to define workflow like this
//
//	DeployWorkflow(ctx workflow.Context, resource []byte)
//
// problem with this approach: we have to serialize and deserialize the data every time especially when the resource is modified by
// the workflow
//
//	r := resource.CreateDummy()
//	wrkr := shared.Temporal().Worker(shared.CoreQueue)
//	wrkr.RegisterWorkflow(r.DeployWorkflow)
//	wrkr.RegisterWorkflow(r.UpdateTraffic)
func (c *core) RegisterCloudResource(provider CloudProvider, driver Driver, resource ResourceConstructor) {
	// TODO: replace this with Once
	if c.resources[provider] == nil {
		c.resources[provider] = make(map[Driver]ResourceConstructor)
	}

	c.resources[provider][driver] = resource
}

func (c *core) RegisterRepoProvider(provider RepoProvider, activities RepoIO) {
	c.providers.repos[provider] = activities
}

func (c *core) RegisterCloudProvider(provider CloudProvider, activities CloudIO) {
	c.providers.cloud[provider] = activities
}

func (c *core) RepoIO(name RepoProvider) RepoIO {
	if p, ok := c.providers.repos[name]; ok {
		return p
	}

	panic(NewProviderNotFoundError(name.String()))
}

func (c *core) ResourceConstructor(provider CloudProvider, driver Driver) ResourceConstructor {
	p, ok := c.resources[provider]
	if !ok {
		panic(NewProviderNotFoundError(provider.String()))
	}

	if r, ok := p[driver]; ok {
		return r
	}

	panic(NewResourceNotFoundError(driver.String(), provider.String()))
}

func (c *core) CloudProvider(name CloudProvider) CloudIO {
	if p, ok := c.providers.cloud[name]; ok {
		return p
	}

	panic(NewProviderNotFoundError(name.String()))
}

func (c *core) ResgisterMessageProvider(provider MessageProvider, activities MessageIO) {
	c.providers.message[provider] = activities
}

func (c *core) MessageProvider(name MessageProvider) MessageIO {
	if p, ok := c.providers.message[name]; ok {
		return p
	}

	panic(NewProviderNotFoundError(name.String()))
}

// WithMessageProvider registers a repo provider with the core.
func WithMessageProvider(provider MessageProvider, io MessageIO) Option {
	return func(c Core) {
		shared.Logger().Info("core: registering message provider", "name", provider.String())
		c.ResgisterMessageProvider(provider, io)
	}
}

// WithRepoProvider registers a repo provider with the core.
func WithRepoProvider(provider RepoProvider, io RepoIO) Option {
	return func(c Core) {
		shared.Logger().Info("core: registering repo provider", "name", provider.String())
		c.RegisterRepoProvider(provider, io)
	}
}

// WithCloudProvider registers a cloud provider with the core.
func WithCloudProvider(provider CloudProvider, io CloudIO) Option {
	return func(c Core) {
		shared.Logger().Info("core: registering cloud provider", "name", provider.String())
		c.RegisterCloudProvider(provider, io)
	}
}

func WithCloudResource(provider CloudProvider, driver Driver, res ResourceConstructor) Option {
	return func(c Core) {
		shared.Logger().Info("core: registering cloud resource", "name", driver.String())
		c.RegisterCloudResource(provider, driver, res)
	}
}

// Instance returns a singleton instance of the core. It is best to call this function in the main() function to
// register the providers available to the service. This is because the core uses workflow and providers implementations
// to access the providers.
func Instance(opts ...Option) Core {
	if instance == nil {
		shared.Logger().Info("core: instance not initialized, initializing now ...")
		once.Do(func() {
			instance = &core{
				providers: Providers{
					repos:   make(map[RepoProvider]RepoIO),
					cloud:   make(map[CloudProvider]CloudIO),
					message: make(map[MessageProvider]MessageIO),
				},

				resources: make(map[CloudProvider]map[Driver]ResourceConstructor),
			}

			for _, opt := range opts {
				opt(instance)
			}
		})
	}

	return instance
}

func getRegion(provider CloudProvider, blueprint *Blueprint) string {
	switch provider {
	case CloudProviderAWS:
		return blueprint.Regions.Aws[0]
	case CloudProviderGCP:
		return blueprint.Regions.Gcp[0]
	case CloudProviderAzure:
		return blueprint.Regions.Azure[0]
	}

	return ""
}

// getProviderConfig gets a specific provider config from blueprint

// TODO: the provider config only has GCP config, this was hardcoded for demo,
// need to make it generic so we can get any provider config based on its name.
func getProviderConfig(provider CloudProvider, blueprint *Blueprint) string {
	switch provider {
	case CloudProviderAWS:
		return "blueprint.ProviderConfig.Aws"
	case CloudProviderGCP:
		return blueprint.ProviderConfig
	case CloudProviderAzure:
		return "blueprint.ProviderConfig.Azure"
	}

	return ""
}
