package service

import (
	"github.com/KirillMironov/ci/internal/domain"
	"gopkg.in/yaml.v2"
)

// YAMLParser is a service that can parse a pipeline from a YAML file.
type YAMLParser struct{}

// ParsePipeline parses a pipeline from a YAML file.
func (YAMLParser) ParsePipeline(b []byte) (pipeline domain.Pipeline, _ error) {
	return pipeline, yaml.Unmarshal(b, &pipeline)
}
