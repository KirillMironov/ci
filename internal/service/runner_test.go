package service

import (
	"context"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/KirillMironov/ci/pkg/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRunner(t *testing.T) {
	var expectedLog = "ok"

	tests := map[string]struct {
		executor       executor
		expectedStatus domain.Status
	}{
		"success": {
			executor: mock.Executor{
				HasError: false,
				Log:      expectedLog,
			},
			expectedStatus: domain.Success,
		},
		"failure": {
			executor: mock.Executor{
				HasError: true,
				Log:      expectedLog,
			},
			expectedStatus: domain.Failure,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			var (
				buildsStorage = mock.NewBuilds()
				runner        = NewRunner(tc.executor, buildsStorage, mock.Logger{})
				req           = runRequest{
					ctx:    ctx,
					repoId: "0",
					commit: domain.Commit{Hash: "123"},
					pipeline: domain.Pipeline{
						Name: "test",
						Steps: []domain.Step{
							{},
						},
					},
					srcCodePath: ".",
				}
			)

			go runner.Start(ctx)

			runner.Run(req)

			builds, err := buildsStorage.GetAllByRepoId(req.repoId)
			require.NoError(t, err)
			require.Len(t, builds, 1)

			build := builds[0]
			assert.NotEmpty(t, build.Id)
			assert.Equal(t, req.repoId, build.RepoId)
			assert.Equal(t, req.commit, build.Commit)
			assert.Equal(t, expectedLog, build.Log.Data)
			assert.Equal(t, tc.expectedStatus, build.Status)
			assert.True(t, time.Now().After(build.CreatedAt))
		})
	}
}
