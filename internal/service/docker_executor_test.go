//go:build integration

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

	var executor = NewDockerExecutor(cli, "/ci", &TarArchiver{})

	tests := map[string]struct {
		step          domain.Step
		expectedLogs  string
		expectedError error
	}{
		"success": {
			step: domain.Step{
				Name:        "echo",
				Image:       "busybox:1.35",
				Environment: []string{"FOO=BAR"},
				Command:     []string{"/bin/sh", "-c"},
				Args:        []string{"echo $FOO"},
			},
			expectedLogs:  "BAR\r\n",
			expectedError: nil,
		},
		"exit error code": {
			step: domain.Step{
				Name:    "exit",
				Image:   "busybox:1.35",
				Command: []string{"/bin/sh", "-c"},
				Args:    []string{"echo hello; exit 1"},
			},
			expectedLogs:  "hello\r\n",
			expectedError: domain.ExitError{Code: 1},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			logs, err := executor.ExecuteStep(context.Background(), tc.step, t.TempDir())
			assert.ErrorIs(t, err, tc.expectedError)

			data, _ := io.ReadAll(logs)
			assert.Equal(t, tc.expectedLogs, string(data))
		})
	}
}
