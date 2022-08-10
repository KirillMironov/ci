package service

import (
	"github.com/KirillMironov/ci/internal/domain"
	"gopkg.in/yaml.v2"
)

// YAMLParser used to parse YAML files into pipelines.
type YAMLParser struct{}

// ParsePipeline parses a pipeline from a given YAML file.
func (YAMLParser) ParsePipeline(b []byte) (pipeline domain.Pipeline, _ error) {
	return pipeline, yaml.Unmarshal(b, &pipeline)
}
