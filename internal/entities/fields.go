package entities

import (
	"encoding/json"

	"github.com/gocql/gocql"
)

type AppConfig struct{}

func (config AppConfig) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return json.Marshal(config)
}

func (config *AppConfig) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	return json.Unmarshal(data, config)
}

type BluePrintRegions struct {
	GCP     []string `json:"gcp"`
	AWS     []string `json:"aws"`
	Azure   []string `json:"azure"`
	Default string   `json:"default"`
}

func (regions BluePrintRegions) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return json.Marshal(regions)
}

func (regions *BluePrintRegions) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	return json.Unmarshal(data, regions)
}

type RolloutArtifact struct{}

type RolloutArtifacts map[string]RolloutArtifact

func (artifacts RolloutArtifacts) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return json.Marshal(artifacts)
}

func (artifacts *RolloutArtifacts) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	return json.Unmarshal(data, artifacts)
}
