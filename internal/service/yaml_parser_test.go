package service

import (
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestYAMLParser_ParsePipeline(t *testing.T) {
	var parser YAMLParser
	var yaml = `
name: example

steps:
  - name: version
    image: golang:1.18.3-alpine3.15
    command: ["/bin/sh", "-c"]
    args:
      - go version

  - name: env
    image: busybox:1.35
    env:
      - TEST=true
    command: ["/bin/sh", "-c"]
    args:
      - echo $TEST
      - printenv
`

	pipeline, err := parser.ParsePipeline([]byte(yaml))
	assert.NoError(t, err)
	assert.Equal(t, domain.Pipeline{
		Name: "example",
		Steps: []domain.Step{
			{
				Name:    "version",
				Image:   "golang:1.18.3-alpine3.15",
				Command: []string{"/bin/sh", "-c"},
				Args:    []string{"go version"},
			},
			{
				Name:        "env",
				Image:       "busybox:1.35",
				Environment: []string{"TEST=true"},
				Command:     []string{"/bin/sh", "-c"},
				Args:        []string{"echo $TEST", "printenv"},
			},
		},
	}, pipeline)

	pipeline, err = parser.ParsePipeline([]byte("-"))
	assert.Error(t, err)
	assert.Empty(t, pipeline)
}
