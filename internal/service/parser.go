package service

import (
	"github.com/KirillMironov/ci/internal/domain"
	"gopkg.in/yaml.v2"
)

// Parser is a service that can parse a pipeline.
type Parser struct{}

// ParsePipeline parses a pipeline from a YAML file.
func (Parser) ParsePipeline(b []byte) (pipeline domain.Pipeline, _ error) {
	return pipeline, yaml.Unmarshal(b, &pipeline)
}
