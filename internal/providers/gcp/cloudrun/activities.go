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

package cloudrun

import (
	"context"
	"strconv"
	"strings"

	run "cloud.google.com/go/run/apiv2"
	"cloud.google.com/go/run/apiv2/runpb"
	"go.temporal.io/sdk/activity"
)

type (
	Activities struct{}
)

// DeployRevision deploys a new revision on Resource if the service is already created.
// If no service is running, then it will create a new service and deploy first revision.
func (a *Activities) DeployRevision(ctx context.Context, r *Resource, wl *Workload) error {
	logger := activity.GetLogger(ctx)

	client, err := run.NewServicesRESTClient(context.Background())
	if err != nil {
		logger.Error("could not create service client", "error", err)
	}

	defer func() { _ = client.Close() }()

	// Create service if this is the first revision
	if r.Revision == r.GetFirstRevision() {
		service := r.GetServiceTemplate(ctx, wl)
		logger.Info("deploying service", "data", service, "parent", r.GetParent(), "ID", wl.Name)

		csr := &runpb.CreateServiceRequest{Parent: r.GetParent(), Service: service, ServiceId: wl.Name}

		op, err := client.CreateService(ctx, csr)
		if err != nil {
			logger.Error("Could not create service", "Error", err)
			return err
		}

		logger.Info("waiting for service creation")

		_, _ = op.Wait(ctx)
	} else { // otherwise create a new revision and route 50% traffic to it
		// Get the already deployed service on cloud run.
		// TODO: We should be able to construct the service template of currently deployed service by caching the data in quantm
		req := &runpb.GetServiceRequest{Name: r.GetParent() + "/services/" + wl.Name}
		service, err := client.GetService(ctx, req)

		if err != nil {
			logger.Error("could not get service", "Error", err)
			return err
		}

		// assuming there is no side container on cloud run
		service.Template.Containers[0].Image = wl.Image

		logger.Info("50 percent traffic to latest", "revision", r.Revision)

		tt := &runpb.TrafficTarget{Type: runpb.TrafficTargetAllocationType_TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST, Percent: 50}

		tt1 := &runpb.TrafficTarget{
			Type: runpb.TrafficTargetAllocationType_TRAFFIC_TARGET_ALLOCATION_TYPE_REVISION, Revision: r.LastRevision, Percent: 50,
		}

		service.Traffic = []*runpb.TrafficTarget{tt, tt1}
		service.Template.Revision = r.Revision

		usr := &runpb.UpdateServiceRequest{Service: service}
		op, err := client.UpdateService(ctx, usr)

		if err != nil {
			logger.Error("could not update service", "Error", err)
			return err
		}

		logger.Info("waiting for service revision update")

		_, _ = op.Wait(ctx)
	}

	// Allow access to all users
	if r.AllowUnauthenticatedAccess {
		_ = r.AllowAccessToAll(ctx)
	}

	return nil
}

// UpdateTrafficActivity updates traffic percentage on a cloud run resource
// This cannot be done in the workflow because of the blocking updateservice call.
func (a *Activities) UpdateTrafficActivity(ctx context.Context, r *Resource, percent int32) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Update traffic", "revision", r.Revision, "percentage", percent)

	service := r.GetService(ctx)
	svctx := context.Background()

	serviceClient, err := run.NewServicesRESTClient(svctx)
	if err != nil {
		logger.Error("New service rest client", "Error", err)
		return nil
	}

	defer func() { _ = serviceClient.Close() }()

	if r.Revision == r.GetFirstRevision() {
		ttc := &runpb.TrafficTarget{Type: runpb.TrafficTargetAllocationType_TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST, Percent: 100}
		service.Traffic = []*runpb.TrafficTarget{ttc}
	} else {
		ttc := &runpb.TrafficTarget{Type: runpb.TrafficTargetAllocationType_TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST, Percent: percent}
		ttp := &runpb.TrafficTarget{
			Type:     runpb.TrafficTargetAllocationType_TRAFFIC_TARGET_ALLOCATION_TYPE_REVISION,
			Revision: r.LastRevision,
			Percent:  100 - percent,
		}
		service.Traffic = []*runpb.TrafficTarget{ttc, ttp}
	}

	req := &runpb.UpdateServiceRequest{Service: service}
	lro, err := serviceClient.UpdateService(svctx, req)

	if err != nil {
		logger.Error("Update service", "Error", err)
		return err
	} else {
		logger.Info("waiting for service update")

		_, _ = lro.Wait(svctx)
	}

	return nil
}

// GetNextRevision Gets next revision Name to be deployed
// TODO: save the active resource's data on each deployment and on next deployment trigger get the associated data from the
// saved deployment.
func (a *Activities) GetNextRevision(ctx context.Context, r *Resource) (*Resource, error) {
	revision := r.GetFirstRevision()
	r.LastRevision = ""

	// get the deployed service, if not found then it will be first revision
	//TODO: we should get the revision from the saved cache is quantm. We should not have to Get cloud run service for it
	svc := r.GetService(ctx)

	if svc != nil {
		rev := svc.Template.Revision
		r.LastRevision = rev

		// FIXME: revision name would be <service name>-<revision number> e.g first revision for helloworld service would be helloworld-0,
		// second will be helloworld-1. it should be helloworld-${changesetID}
		ss := strings.Split(rev, r.ServiceName+"-")
		revVersion, _ := strconv.Atoi(ss[1])
		revVersion++
		revision = r.ServiceName + "-" + strconv.Itoa(revVersion)
	}

	r.Revision = revision
	activity.GetLogger(ctx).Info("Next revision", "name", revision)

	return r, nil
}
