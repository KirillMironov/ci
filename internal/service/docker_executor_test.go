package service

import (
	"context"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/docker/docker/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
)

func TestDockerExecutor_ExecuteStep(t *testing.T) {
	cli, err := client.NewClientWithOpts()
	require.NoError(t, err)

	var executor = NewDockerExecutor(cli, "/ci")

	logs, err := executor.ExecuteStep(context.Background(), domain.Step{
		Name:        "echo",
		Image:       "busybox:1.35",
		Environment: []string{"FOO=BAR"},
		Command:     []string{"/bin/sh", "-c"},
		Args:        []string{"echo $FOO"},
	}, nil)
	assert.NoError(t, err)
	defer logs.Close()

	data, err := io.ReadAll(logs)
	assert.NoError(t, err)
	assert.Equal(t, "BAR\r\n", string(data))
}

func TestDockerExecutor_ExecuteStepError(t *testing.T) {
	cli, err := client.NewClientWithOpts()
	require.NoError(t, err)

	var executor = NewDockerExecutor(cli, "/ci")
	logs, err := executor.ExecuteStep(context.Background(), domain.Step{
		Name:    "exit",
		Image:   "busybox:1.35",
		Command: []string{"/bin/sh", "-c"},
		Args:    []string{"echo hello; exit 1"},
	}, nil)
	assert.Error(t, err)
	defer logs.Close()

	data, err := io.ReadAll(logs)
	assert.NoError(t, err)
	assert.Equal(t, "hello\r\n", string(data))
}
