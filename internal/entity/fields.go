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

package entity

import (
	"encoding/json"
	"errors"

	"github.com/gocql/gocql"
)

type (
	// StackConfig holds the configuration for an application.
	StackConfig struct{}

	// BluePrintRegions sets the cloud regions where a blueprint can be deployed.
	BluePrintRegions struct {
		GCP     []string `json:"gcp"`
		AWS     []string `json:"aws"`
		Azure   []string `json:"azure"`
		Default string   `json:"default"`
	}

	RepoProvider        string
	RepoProviderMapType map[string]RepoProvider

	// RolloutState is the state of a rollout.
	RolloutState        string
	RolloutStateMapType map[string]RolloutState

	ChangeSetRepoMarker struct {
		Provider   string `json:"provider"`
		CommitID   string `json:"commit_id"`
		HasChanged bool   `json:"changed"`
	}

	ChangeSetRepoMarkers []ChangeSetRepoMarker
)

const (
	RolloutStateQueued     RolloutState = "queued"
	RolloutStateInProgress RolloutState = "in_progress"
	RolloutStateCompleted  RolloutState = "completed"
	RolloutStateRejected   RolloutState = "rejected"
)

var (
	RolloutStateMap = RolloutStateMapType{
		RolloutStateQueued.String():     RolloutStateQueued,
		RolloutStateInProgress.String(): RolloutStateInProgress,
		RolloutStateCompleted.String():  RolloutStateCompleted,
		RolloutStateRejected.String():   RolloutStateRejected,
	}
)

const (
	RepoProviderGithub    RepoProvider = "github"
	RepoProviderGitlab    RepoProvider = "gitlab"
	RepoProviderBitbucket RepoProvider = "bitbucket"
)

var (
	RepoProviderMap = RepoProviderMapType{
		RepoProviderGithub.String():    RepoProviderGithub,
		RepoProviderGitlab.String():    RepoProviderGitlab,
		RepoProviderBitbucket.String(): RepoProviderBitbucket,
	}
)

func (config StackConfig) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return json.Marshal(config)
}

func (config *StackConfig) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	return json.Unmarshal(data, config)
}

func (regions BluePrintRegions) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return json.Marshal(regions)
}

func (regions *BluePrintRegions) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	return json.Unmarshal(data, regions)
}

func (marker ChangeSetRepoMarkers) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return json.Marshal(marker)
}

func (marker *ChangeSetRepoMarkers) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	return json.Unmarshal(data, marker)
}

func (rs RolloutState) String() string {
	return string(rs)
}

func (rs RolloutState) MarshalJSON() ([]byte, error) {
	return json.Marshal(rs.String())
}

func (rs *RolloutState) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	val, ok := RolloutStateMap[s]
	if !ok {
		return errors.New("invalid rollout state")
	}

	*rs = val

	return nil
}

func (rs RolloutState) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return json.Marshal(rs)
}

func (rs *RolloutState) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	return json.Unmarshal(data, rs)
}

func (rp RepoProvider) String() string {
	return string(rp)
}

func (rp RepoProvider) MarshalJSON() ([]byte, error) {
	return json.Marshal(rp.String())
}

func (rp *RepoProvider) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	val, ok := RepoProviderMap[s]
	if !ok {
		return errors.New("invalid repo provider")
	}

	*rp = val

	return nil
}

func (rp RepoProvider) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return json.Marshal(rp)
}

func (rp *RepoProvider) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	return json.Unmarshal(data, rp)
}
