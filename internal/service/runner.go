package service

import (
	"bytes"
	"context"
	"github.com/KirillMironov/ci/internal/domain"
	"io"
	"os"
)

type Runner struct {
	executor executor
}

type executor interface {
	ExecuteStep(ctx context.Context, step domain.Step, srcCodeArchive io.Reader) (logs io.ReadCloser, err error)
}

func NewRunner(executor executor) *Runner {
	return &Runner{executor: executor}
}

func (r Runner) Run(ctx context.Context, pipeline domain.Pipeline, srcCodeArchivePath string) (logs []byte, err error) {
	var buf bytes.Buffer

	for _, step := range pipeline.Steps {
		err = func() error {
			archive, err := os.Open(srcCodeArchivePath)
			if err != nil {
				return err
			}
			defer archive.Close()

			stepLogs, err := r.executor.ExecuteStep(ctx, step, archive)
			if err != nil {
				return err
			}
			defer stepLogs.Close()

			_, err = io.Copy(&buf, stepLogs)
			return err
		}()
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}
