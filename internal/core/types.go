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
	"go.breu.io/quantm/internal/shared"
)

// Stack Signals.
const (
	StackSignalAssetsRerieved     shared.WorkflowSignal = "stack__assets_retrieved"
	StackSignalInfraProvisioned   shared.WorkflowSignal = "stack__infra_provisioned"
	StackSignalDeploymentComplete shared.WorkflowSignal = "stack__deployment_completed"
	StackSignalManualOverride     shared.WorkflowSignal = "stack__manual_override"
)

// MessageIO payloads.
type (
	MessageIOSendStaleBranchMessagePayload struct {
		TeamID string        `json:"team_id"`
		Stale  *LatestCommit `json:"slate"`
	}

	MessageIOSendNumberOfLinesExceedMessagePayload struct {
		TeamID        string         `json:"team_id"`
		RepoName      string         `json:"repo_name"`
		BranchName    string         `json:"branch_name"`
		Threshold     int            `json:"threshold"`
		BranchChnages *BranchChanges `json:"branch_chnages"`
	}

	MessageIOSendMergeConflictsMessagePayload struct {
		TeamID string        `json:"team_id"`
		Merge  *LatestCommit `json:"merge"`
	}

	MessageIOCompleteOauthResponsePayload struct {
		Code string `json:"code"`
	}
)
