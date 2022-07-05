package service

import (
	"github.com/KirillMironov/ci/internal/domain"
	"gopkg.in/yaml.v2"
)

type Parser struct{}

// ParsePipeline parses a pipeline from a string.
func (p *Parser) ParsePipeline(str string) (pipeline domain.Pipeline, _ error) {
	return pipeline, yaml.Unmarshal([]byte(str), &pipeline)
}
