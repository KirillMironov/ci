package service

import (
	"bytes"
	"context"
	"github.com/KirillMironov/ci/internal/domain"
	"io"
	"os"
)

type Runner struct {
	ciFilename string
	cloner     cloner
	executor   executor
	finder     finder
	parser     parser
	logs       logs
}

// logs stores build logs.
type logs interface {
	Put(domain.Log) (id int, err error)
}

func NewRunner(ciFilename string, cloner cloner, executor executor, finder finder, parser parser, logs logs) *Runner {
	return &Runner{
		ciFilename: ciFilename,
		cloner:     cloner,
		executor:   executor,
		finder:     finder,
		parser:     parser,
		logs:       logs,
	}
}

type (
	// cloner is a service that can clone a repository.
	cloner interface {
		GetLatestCommitHash(url, branch string) (string, error)
		CloneRepository(url, branch, hash string) (archivePath string, removeArchive func(), err error)
	}
	// executor is a service that can execute pipeline steps.
	executor interface {
		ExecuteStep(ctx context.Context, step domain.Step, sourceCodeArchive io.Reader) (logs io.ReadCloser, err error)
	}
	// finder is a service that can find a file in a given archive.
	finder interface {
		FindFile(filename, archivePath string) ([]byte, error)
	}
	// parser is a service that can parse a pipeline.
	parser interface {
		ParsePipeline(b []byte) (domain.Pipeline, error)
	}
)

func (r Runner) Run(ctx context.Context, repo domain.Repository, targetHash string) (logId int, err error) {
	archivePath, removeArchive, err := r.cloner.CloneRepository(repo.URL, repo.Branch, targetHash)
	if err != nil {
		return 0, err
	}
	defer removeArchive()

	yaml, err := r.finder.FindFile(r.ciFilename, archivePath)
	if err != nil {
		return 0, err
	}

	pipeline, err := r.parser.ParsePipeline(yaml)
	if err != nil {
		return 0, err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var buf bytes.Buffer

	for _, step := range pipeline.Steps {
		log, err := func() ([]byte, error) {
			archive, err := os.Open(archivePath)
			if err != nil {
				return nil, err
			}
			defer archive.Close()

			logsReader, err := r.executor.ExecuteStep(ctx, step, archive)
			if err != nil {
				return nil, err
			}
			defer logsReader.Close()

			return io.ReadAll(logsReader)
		}()
		if err != nil {
			return 0, err
		}

		buf.Write(log)
	}

	return r.logs.Put(domain.Log{
		Data: buf.Bytes(),
	})
}
