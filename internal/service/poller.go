package service

import (
	"archive/tar"
	"bytes"
	"context"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/KirillMironov/ci/pkg/logger"
	"io"
	"os"
	"time"
)

// Poller is a service that can poll a VCS and execute a pipeline.
type Poller struct {
	cloner   cloner
	parser   parser
	executor executor
	logger   logger.Logger
}

// cloner is a service that can clone a repository.
type cloner interface {
	CloneRepository(url string) (sourceCodeArchivePath string, removeArchive func() error, err error)
}

// parser is a service that can parse a pipeline.
type parser interface {
	ParsePipeline(b []byte) (domain.Pipeline, error)
}

// executor is a service that can execute pipeline steps.
type executor interface {
	Execute(ctx context.Context, step domain.Step, sourceCodeArchive io.Reader) error
}

// NewPoller creates a new Poller.
func NewPoller(cloner cloner, parser parser, executor executor, logger logger.Logger) *Poller {
	return &Poller{
		cloner:   cloner,
		parser:   parser,
		executor: executor,
		logger:   logger,
	}
}

// Start starts VCS polling with a given interval.
func (p Poller) Start(vcs domain.VCS) {
	timer := time.NewTimer(vcs.PollingInterval)

	for range timer.C {
		err := p.poll(vcs)
		if err != nil {
			p.logger.Error(err)
		}
		timer.Reset(vcs.PollingInterval)
	}
}

// poll polls fetches a pipeline from a VCS and executes it.
func (p Poller) poll(vcs domain.VCS) error {
	sourceCodePath, remove, err := p.cloner.CloneRepository(vcs.URL)
	if err != nil {
		return err
	}
	defer func() {
		if err = remove(); err != nil {
			p.logger.Error(err)
		}
	}()

	yaml, err := p.findPipeline(sourceCodePath)
	if err != nil {
		return err
	}

	pipeline, err := p.parser.ParsePipeline(yaml)
	if err != nil {
		return err
	}

	for _, step := range pipeline.Steps {
		file, err := os.Open(sourceCodePath)
		if err != nil {
			return err
		}

		err = p.executor.Execute(context.Background(), step, file)
		if err != nil {
			p.logger.Error(err)
			file.Close()
			return err
		}

		file.Close()
	}

	return nil
}

// findPipeline finds a pipeline in a source code archive.
func (Poller) findPipeline(archivePath string) ([]byte, error) {
	const yamlFilename = ".ci.yaml"

	file, err := os.Open(archivePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	tarReader := tar.NewReader(file)

	buf := bytes.NewBuffer(nil)

	for {
		header, err := tarReader.Next()
		if err != nil {
			return nil, err
		}
		if header.Name == yamlFilename {
			_, err = io.Copy(buf, tarReader)
			if err != nil {
				return nil, err
			}
			return buf.Bytes(), nil
		}
	}
}
