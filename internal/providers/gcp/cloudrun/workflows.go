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
	"time"

	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/shared"
)

// TODO: ysf: i don't think we need workflows here, or should we have workflows here?

type (
	workflows struct{}
)

func (w *workflows) DeployWorkflow(ctx workflow.Context, r *Resource, wl *Workload) (*Resource, error) {
	r.ServiceName = wl.Name
	activityOpts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	actx := workflow.WithActivityOptions(ctx, activityOpts)

	err := workflow.ExecuteActivity(actx, activities.GetNextRevision, r).Get(actx, r)
	if err != nil {
		shared.Logger().Error("Error in Executing activity: GetNextRevision", "error", err)
		return r, err
	}

	err = workflow.ExecuteActivity(actx, activities.DeployRevision, r, wl).Get(actx, nil)
	if err != nil {
		shared.Logger().Error("Error in Executing activity: DeployDummy", "error", err)
		return r, err
	}

	return r, nil
}

// UpdateTraffic updates the traffic to given percentage.
func (w *workflows) UpdateTraffic(ctx workflow.Context, r *Resource, percent int32) error {
	shared.Logger().Info("Distributing traffic between revisions", r.Revision, r.LastRevision)

	activityOpts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	actx := workflow.WithActivityOptions(ctx, activityOpts)

	err := workflow.ExecuteActivity(actx, activities.UpdateTrafficActivity, r, percent).Get(ctx, r)
	if err != nil {
		shared.Logger().Error("Error in Executing activity: UpdateTrafficActivity", "error", err)
		return err
	}

	return nil
}
