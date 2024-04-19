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
		RegisterRepoProvider(RepoProvider, RepoProviderActivities)
		RegisterCloudProvider(CloudProvider, CloudProviderActivities)
		RegisterCloudResource(provider CloudProvider, driver Driver, resource ResourceConstructor)
		ResgisterMessageProvider(MessageProvider, MessageProviderActivities)

		RepoProvider(RepoProvider) RepoProviderActivities
		CloudProvider(CloudProvider) CloudProviderActivities
		ResourceConstructor(CloudProvider, Driver) ResourceConstructor
		MessageProvider(MessageProvider) MessageProviderActivities
	}

	Option func(Core)

	RepoProviderActivities interface {
		GetLatestCommit(context.Context, string, string) (string, error)
		DeployChangeset(ctx context.Context, repoID string, changesetID *gocql.UUID) error
		TagCommit(ctx context.Context, repoID string, commitSHA string, tagName string, tagMessage string) error
		CreateBranch(ctx context.Context, installationID int64, repoID int64, repoName string, repoOwner string, targetCommit string,
			newBranchName string) error
		DeleteBranch(ctx context.Context, installationID int64, repoName string, repoOwner string, branchName string) error
		MergeBranch(ctx context.Context, installationID int64, repoName string, repoOwner string, baseBranch string,
			targetBranch string) error
		ChangesInBranch(ctx context.Context, installationID int64, repoName string, repoOwner string, defaultBranch string,
			targetBranch string) (*shared.BranchChanges, error)
		GetAllBranches(ctx context.Context, installationID int64, repoName string, repoOwner string) ([]string, error)
		TriggerCIAction(ctx context.Context, installationID int64, repoOwner string, repoName string, targetBranch string) error
		GetRepoTeamID(ctx context.Context, repoID string) (string, error)
	}

	CloudProviderActivities interface {
		FillMe()
	}

	MessageProviderActivities interface {
		SendChannelMessage(ctx context.Context, teamID, msg string) error // TODO: figure out the signature
	}

	Providers struct {
		repos   map[RepoProvider]RepoProviderActivities
		cloud   map[CloudProvider]CloudProviderActivities
		message map[MessageProvider]MessageProviderActivities
	}

	CloudResource interface {
		Provision(ctx workflow.Context) (workflow.Future, error)
		DeProvision() error
		Deploy(workflow.Context, []Workload, gocql.UUID) error
		UpdateTraffic(workflow.Context, int32) error
		Marshal() ([]byte, error)
	}

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

func (c *core) RegisterRepoProvider(provider RepoProvider, activities RepoProviderActivities) {
	c.providers.repos[provider] = activities
}

func (c *core) RegisterCloudProvider(provider CloudProvider, activities CloudProviderActivities) {
	c.providers.cloud[provider] = activities
}

func (c *core) RepoProvider(name RepoProvider) RepoProviderActivities {
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

func (c *core) CloudProvider(name CloudProvider) CloudProviderActivities {
	if p, ok := c.providers.cloud[name]; ok {
		return p
	}

	panic(NewProviderNotFoundError(name.String()))
}

func (c *core) ResgisterMessageProvider(provider MessageProvider, activities MessageProviderActivities) {
	c.providers.message[provider] = activities
}

func (c *core) MessageProvider(name MessageProvider) MessageProviderActivities {
	if p, ok := c.providers.message[name]; ok {
		return p
	}

	panic(NewProviderNotFoundError(name.String()))
}

// WithMessageProvider registers a repo provider with the core.
func WithMessageProvider(name MessageProvider, provider MessageProviderActivities) Option {
	return func(c Core) {
		shared.Logger().Info("core: registering message provider", "name", name.String())
		c.ResgisterMessageProvider(name, provider)
	}
}

// WithRepoProvider registers a repo provider with the core.
func WithRepoProvider(name RepoProvider, provider RepoProviderActivities) Option {
	return func(c Core) {
		shared.Logger().Info("core: registering repo provider", "name", name.String())
		c.RegisterRepoProvider(name, provider)
	}
}

// WithCloudProvider registers a cloud provider with the core.
func WithCloudProvider(name CloudProvider, provider CloudProviderActivities) Option {
	return func(c Core) {
		shared.Logger().Info("core: registering cloud provider", "name", name.String())
		c.RegisterCloudProvider(name, provider)
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
					repos:   make(map[RepoProvider]RepoProviderActivities),
					cloud:   make(map[CloudProvider]CloudProviderActivities),
					message: make(map[MessageProvider]MessageProviderActivities),
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
