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
	_ = json.Unmarshal([]byte(config), &cr.Config)

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
