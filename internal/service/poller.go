package service

import (
	"bytes"
	"context"
	"github.com/KirillMironov/ci/internal/domain"
	"io"
	"os"
	"time"
)

// Poller is a service that can poll a source code repository.
type Poller struct {
	ciFilename string
	cloner     cloner
	executor   executor
	finder     finder
	parser     parser
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

// NewPoller creates a new Poller.
func NewPoller(ciFilename string, cloner cloner, executor executor, finder finder, parser parser) *Poller {
	return &Poller{
		ciFilename: ciFilename,
		cloner:     cloner,
		executor:   executor,
		finder:     finder,
		parser:     parser,
	}
}

// Poll starts repository polling with a given interval.
func (p Poller) Poll(ctx context.Context, repo domain.Repository, prevHash string) (latestHash string, logs []byte,
	err error) {
	timer := time.NewTimer(repo.PollingInterval)

	for {
		select {
		case <-ctx.Done():
			return "", nil, ctx.Err()
		case <-timer.C:
			latestHash, err = p.cloner.GetLatestCommitHash(repo.URL, repo.Branch)
			if err != nil {
				return "", nil, err
			}
			if latestHash == prevHash {
				timer.Reset(repo.PollingInterval)
				continue
			}

			logs, err = p.poll(ctx, repo, latestHash)
			return latestHash, logs, err
		}
	}
}

// poll clones a repository and executes a pipeline.
func (p Poller) poll(ctx context.Context, repo domain.Repository, hash string) ([]byte, error) {
	archivePath, removeArchive, err := p.cloner.CloneRepository(repo.URL, repo.Branch, hash)
	if err != nil {
		return nil, err
	}
	defer removeArchive()

	yaml, err := p.finder.FindFile(p.ciFilename, archivePath)
	if err != nil {
		return nil, err
	}

	pipeline, err := p.parser.ParsePipeline(yaml)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var buf bytes.Buffer

	for _, step := range pipeline.Steps {
		log, err := p.executeStep(ctx, step, archivePath)
		if err != nil {
			return nil, err
		}
		buf.Write(log)
	}

	return buf.Bytes(), nil
}

// executeStep executes a pipeline step.
func (p Poller) executeStep(ctx context.Context, step domain.Step, archivePath string) ([]byte, error) {
	archive, err := os.Open(archivePath)
	if err != nil {
		return nil, err
	}
	defer archive.Close()

	logsReader, err := p.executor.ExecuteStep(ctx, step, archive)
	if err != nil {
		return nil, err
	}
	defer logsReader.Close()

	return io.ReadAll(logsReader)
}
