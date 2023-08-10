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
	"errors"
	"strings"

	artifactregistry "cloud.google.com/go/artifactregistry/apiv1"
	"cloud.google.com/go/artifactregistry/apiv1/artifactregistrypb"
	"github.com/gocql/gocql"
	"go.temporal.io/sdk/activity"

	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/shared"
)

var (
	activities *Activities
)

type (
	Activities            struct{}
	ArtifactRegistryImage struct {
		location   string
		project    string
		pkg        string //image name
		repository string
	}
)

// GetResources gets resources from DB against a stack.
func (a *Activities) GetResources(ctx context.Context, stackID string) (*SlicedResult[Resource], error) {
	log := activity.GetLogger(ctx)
	resources := make([]Resource, 0)
	err := db.Filter(&Resource{}, &resources, db.QueryParams{"stack_id": stackID})

	if err != nil {
		log.Error("GetResources Error", "error", err)
	}

	return &SlicedResult[Resource]{Data: resources}, err
}

// GetWorkloads gets workloads from DB against a stack.
func (a *Activities) GetWorkloads(ctx context.Context, stackID string) (*SlicedResult[Workload], error) {
	log := activity.GetLogger(ctx)
	workloads := make([]Workload, 0)
	err := db.Filter(&Workload{}, &workloads, db.QueryParams{"stack_id": stackID})

	if err != nil {
		log.Error("GetWorkloads Error", "error", err)
	}

	return &SlicedResult[Workload]{Data: workloads}, err
}

// GetWorkloads gets workloads from DB against a stack.
func (a *Activities) GetRepos(ctx context.Context, stackID string) (*SlicedResult[Repo], error) {
	log := activity.GetLogger(ctx)
	repos := make([]Repo, 0)
	err := db.Filter(&Repo{}, &repos, db.QueryParams{"stack_id": stackID})

	if err != nil {
		log.Error("GetRepos Error", "error", err)
	}

	return &SlicedResult[Repo]{Data: repos}, err
}

// GetBluePrint gets blueprint from DB against a stack.
func (a *Activities) GetBluePrint(ctx context.Context, stackID string) (*Blueprint, error) {
	log := activity.GetLogger(ctx)
	blueprint := &Blueprint{}
	params := db.QueryParams{"stack_id": stackID}

	if err := db.Get(blueprint, params); err != nil {
		log.Error("GetBlueprint Error", "stack", stackID, "error", err)
		return blueprint, err
	}

	return blueprint, nil
}

// CreateChangeset create changeset entity with provided ID.
func (a *Activities) CreateChangeset(ctx context.Context, changeSet *ChangeSet, id gocql.UUID) error {
	err := db.CreateWithID(changeSet, id)
	return err
}

// TagGcpImage creates a new tag on a docker image in GCP artifact registry
func (a *Activities) TagGcpImage(ctx context.Context, image string, digest string, tag string) error {

	logger := activity.GetLogger(ctx)
	c, err := artifactregistry.NewRESTClient(ctx)
	if err != nil {
		logger.Error("Could not create REST client for artifact registry", "Error", err)
		return err
	}
	defer c.Close()

	imageparts, err := ParseArtifactRegistryImage(image)
	if err != nil {
		logger.Error("Error in parsing artifact registry image", "Error", err)
		return err
	}

	parent := "projects/" + imageparts.project + "/locations/" + imageparts.location + "/repositories/" + imageparts.repository + "/packages/" + imageparts.pkg
	newtag := &artifactregistrypb.Tag{
		Name:    parent + "/tags/" + tag,
		Version: parent + "/versions/" + digest,
	}
	req := &artifactregistrypb.UpdateTagRequest{Tag: newtag}
	_, err = c.UpdateTag(ctx, req)

	if err != nil {
		logger.Error("Error in updating tag", "Error", err)
		return err
	}

	logger.Info("Tag updated")
	return nil
}

func ParseArtifactRegistryImage(image string) (*ArtifactRegistryImage, error) {
	arImage := new(ArtifactRegistryImage)
	// get location
	image = strings.Replace(image, "-docker.pkg.dev", "", 1)
	result := strings.Split(image, "/")

	// assuming here that the image name will have no slashes except to separate location, repo, project and package
	// sample image: asia-southeast1-docker.pkg.dev/breu-dev/ctrlplane/helloworld
	// sample image: <location>-docker.pkg.dev/<project>/<repository>/<package>
	if len(result) > 4 {
		shared.Logger().Error("Unexpected image string, can't parse", "Image", image)
		return nil, errors.New("Unexpected image string, can't parse")
	}

	arImage.location = result[0]
	arImage.project = result[1]
	arImage.repository = result[2]
	arImage.pkg = result[3]
	return arImage, nil
}
