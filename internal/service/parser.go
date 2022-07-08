package service

import (
	"github.com/KirillMironov/ci/internal/domain"
	"gopkg.in/yaml.v2"
)

type Parser struct{}

func (Parser) ParsePipeline(b []byte) (pipeline domain.Pipeline, _ error) {
	return pipeline, yaml.Unmarshal(b, &pipeline)
}
