package service

import (
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParser_ParsePipeline(t *testing.T) {
	var (
		parser Parser
		yaml   = `
name: example

steps:
  - name: version
    image: golang:1.18.3-alpine3.15
    command:
      - go
      - version

  - name: env
    image: busybox:1.35
    env:
      - TEST=true
    command:
      - printenv
`
	)

	pipeline, err := parser.ParsePipeline(yaml)
	assert.NoError(t, err)
	assert.Equal(t, domain.Pipeline{
		Name: "example",
		Steps: []domain.Step{
			{
				Name:    "version",
				Image:   "golang:1.18.3-alpine3.15",
				Command: []string{"go", "version"},
			},
			{
				Name:        "env",
				Image:       "busybox:1.35",
				Environment: []string{"TEST=true"},
				Command:     []string{"printenv"},
			},
		},
	}, pipeline)
}
