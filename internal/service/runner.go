package service

import (
	"bytes"
	"context"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/KirillMironov/ci/pkg/logger"
	"github.com/rs/xid"
	"io"
)

// Runner used to execute pipeline.
type Runner struct {
	run           chan runRequest
	executor      executor
	buildsStorage domain.BuildsStorage
	logger        logger.Logger
}

type (
	runRequest struct {
		ctx         context.Context
		repoId      string
		commit      domain.Commit
		pipeline    domain.Pipeline
		srcCodePath string
	}
	executor interface {
		ExecuteStep(ctx context.Context, step domain.Step, srcCodePath string) (logs io.ReadCloser, err error)
	}
)

func NewRunner(executor executor, bs domain.BuildsStorage, logger logger.Logger) *Runner {
	return &Runner{
		run:           make(chan runRequest),
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
			var (
				build = domain.Build{
					Id:     xid.New().String(),
					RepoId: req.repoId,
					Commit: req.commit,
					Status: domain.InProgress,
				}
				logsBuf bytes.Buffer
			)

			err := r.buildsStorage.Create(build)
			if err != nil {
				r.logger.Error(err)
				continue
			}

			build.Status = domain.Success

			for _, step := range req.pipeline.Steps {
				err = func() error {
					stepLogs, err := r.executor.ExecuteStep(req.ctx, step, req.srcCodePath)
					_, _ = io.Copy(&logsBuf, stepLogs)
					stepLogs.Close()
					return err
				}()
				if err != nil {
					build.Status = domain.Failure
					break
				}
			}

			build.Log = domain.Log{Data: logsBuf.String()}

			err = r.buildsStorage.Update(build)
			if err != nil {
				r.logger.Error(err)
			}
		}
	}
}

func (r Runner) Run(req runRequest) {
	r.run <- req
}
