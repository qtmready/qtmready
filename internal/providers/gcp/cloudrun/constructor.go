package cloudrun

import (
	"encoding/json"
	"sync"

	"go.breu.io/quantm/internal/core"
	"go.breu.io/quantm/internal/shared"
)

type (
	Constructor struct{}
)

var (
	registerOnce sync.Once
)

// Create creates cloud run resource.
func (c *Constructor) Create(name string, region string, config string, providerConfig string) (core.CloudResource, error) {
	cr := &Resource{Name: name, Region: region}
	json.Unmarshal([]byte(config), &cr.Config)
	cr.AllowUnauthenticatedAccess = true
	cr.Cpu = "2000m"
	cr.Memory = "1024Mi"
	cr.MinInstances = 0
	cr.MaxInstances = 5
	cr.Generation = 2

	cr.Port = 8000
	cr.CpuIdle = true

	// get gcp project from configuration
	pconfig := new(GCPConfig)
	err := json.Unmarshal([]byte(providerConfig), pconfig)

	if err != nil {
		shared.Logger().Error("Unable to parse provider config for cloudrun")
		return nil, err
	}

	cr.Project = pconfig.Project

	shared.Logger().Info("cloud run", "object", providerConfig, "umarshaled", pconfig, "project", cr.Project)

	w := &workflows{}

	registerOnce.Do(func() {
		coreWrkr := shared.Temporal().Worker(shared.CoreQueue)
		coreWrkr.RegisterWorkflow(w.DeployWorkflow)
		coreWrkr.RegisterWorkflow(w.UpdateTraffic)
	})

	return cr, nil
}

// CreateFromJson creates a Resource object from JSON.
func (c *Constructor) CreateFromJson(data []byte) core.CloudResource {
	cr := &Resource{}
	_ = json.Unmarshal(data, cr)

	return cr
}
