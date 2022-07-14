package service

import (
	"context"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/docker/docker/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDockerExecutor_ExecuteStep(t *testing.T) {
	cli, err := client.NewClientWithOpts()
	require.NoError(t, err)

	var executor = NewDockerExecutor(cli, "/ci")

	logs, err := executor.ExecuteStep(context.Background(), domain.Step{
		Name:        "ls",
		Image:       "busybox:1.35",
		Environment: []string{"FOO=BAR"},
		Command:     []string{"ls", "-la"},
	}, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, logs)
}
