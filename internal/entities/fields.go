// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved. 
//
// This software is made available by Breu, Inc., under the terms of the 
// Breu Community License Agreement ("BCL Agreement"), version 1.0, found at  
// https://www.breu.io/license/community. By installating, downloading, 
// accessing, using or distrubting any of the software, you agree to the  
// terms of the license agreement. 
//
// The above copyright notice and the subsequent license agreement shall be 
// included in all copies or substantial portions of the software. 
//
// Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, 
// IMPLIED, STATUTORY, OR OTHERWISE, AND SPECIFICALLY DISCLAIMS ANY WARRANTY OF 
// MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE 
// SOFTWARE. 
//
// Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT 
// LIMITED TO, LOST PROFITS OR ANY CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, 
// OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, ARISING 
// OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY  
// APPLICABLE LAW. 

package entities

import (
	"encoding/json"

	"github.com/gocql/gocql"
)

type (
	AppConfig struct{}

	BluePrintRegions struct {
		GCP     []string `json:"gcp"`
		AWS     []string `json:"aws"`
		Azure   []string `json:"azure"`
		Default string   `json:"default"`
	}

	RolloutArtifact struct{}

	RolloutArtifacts map[string]RolloutArtifact
)

func (config AppConfig) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return json.Marshal(config)
}

func (config *AppConfig) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	return json.Unmarshal(data, config)
}

func (regions BluePrintRegions) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return json.Marshal(regions)
}

func (regions *BluePrintRegions) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	return json.Unmarshal(data, regions)
}

func (artifacts RolloutArtifacts) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return json.Marshal(artifacts)
}

func (artifacts *RolloutArtifacts) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	return json.Unmarshal(data, artifacts)
}
