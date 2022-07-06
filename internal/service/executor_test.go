package service

import (
	"context"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/docker/docker/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestExecutor_Execute(t *testing.T) {
	cli, err := client.NewClientWithOpts()
	require.NoError(t, err)

	var executor = NewExecutor(cli)

	err = executor.Execute(context.Background(), domain.Step{
		Name:        "test",
		Image:       "busybox:1.35",
		Environment: []string{"FOO=BAR"},
		Command:     []string{"printenv"},
	})
	assert.NoError(t, err)
}
