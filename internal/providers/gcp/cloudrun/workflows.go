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

// UpdateTraffic updates the traffic to given percentage
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
