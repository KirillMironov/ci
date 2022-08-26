package service

import (
	"bytes"
	"context"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/KirillMironov/ci/pkg/logger"
	"io"
)

// Runner used to execute pipeline.
type Runner struct {
	run           <-chan RunRequest
	executor      executor
	buildsStorage domain.BuildsStorage
	logger        logger.Logger
}

type (
	RunRequest struct {
		build       domain.Build
		pipeline    domain.Pipeline
		srcCodePath string
	}
	executor interface {
		ExecuteStep(ctx context.Context, step domain.Step, srcCodePath string) (logs io.ReadCloser, err error)
	}
)

func NewRunner(run <-chan RunRequest, executor executor, bs domain.BuildsStorage, logger logger.Logger) *Runner {
	return &Runner{
		run:           run,
		executor:      executor,
		buildsStorage: bs,
		logger:        logger,
	}
}

// Start listens on run channel and executes pipeline steps.
func (r Runner) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			r.logger.Infof("runner stopped: %v", ctx.Err())
			return
		case req := <-r.run:
			var build = req.build
			var logsBuf bytes.Buffer

			build.Status = domain.Success

			for _, step := range req.pipeline.Steps {
				err := func() error {
					stepLogs, err := r.executor.ExecuteStep(ctx, step, req.srcCodePath)
					if err != nil {
						return err
					}
					defer stepLogs.Close()

					_, err = io.Copy(&logsBuf, stepLogs)
					return err
				}()
				if err != nil {
					build.Status = domain.Failure
					break
				}
			}

			build.Log = domain.Log{Data: logsBuf.String()}

			err := r.buildsStorage.Update(build)
			if err != nil {
				r.logger.Error(err)
			}
		}
	}
}
