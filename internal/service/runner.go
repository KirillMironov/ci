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
	buildsUsecase domain.BuildsUsecase
	logsUsecase   domain.LogsUsecase
	logger        logger.Logger
}

type (
	RunRequest struct {
		repoId      string
		commit      domain.Commit
		pipeline    domain.Pipeline
		srcCodePath string
	}
	executor interface {
		ExecuteStep(ctx context.Context, step domain.Step, srcCodePath string) (logs io.ReadCloser, err error)
	}
)

func NewRunner(run <-chan RunRequest, executor executor, bu domain.BuildsUsecase, lu domain.LogsUsecase,
	logger logger.Logger) *Runner {
	return &Runner{
		run:           run,
		executor:      executor,
		buildsUsecase: bu,
		logsUsecase:   lu,
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
			var build = domain.Build{Commit: req.commit}
			var logsBuf bytes.Buffer

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

			logId, err := r.logsUsecase.Create(ctx, domain.Log{Data: logsBuf.Bytes()})
			if err != nil {
				r.logger.Errorf("failed to create log: %v", err)
				return
			}

			build.LogId = logId

			err = r.buildsUsecase.Create(ctx, build, req.repoId)
			if err != nil {
				r.logger.Errorf("failed to create build: %v", err)
			}
		}
	}
}
