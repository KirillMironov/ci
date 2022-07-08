package service

import (
	"archive/tar"
	"bytes"
	"context"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/KirillMironov/ci/pkg/logger"
	"io"
	"time"
)

type Poller struct {
	cloner   cloner
	parser   parser
	executor executor
	logger   logger.Logger
}

type cloner interface {
	CloneRepository(url string) (sourceCode io.ReadCloser, err error)
}

type parser interface {
	ParsePipeline(b []byte) (domain.Pipeline, error)
}

type executor interface {
	Execute(ctx context.Context, step domain.Step, sourceCode io.Reader) error
}

func NewPoller(cloner cloner, parser parser, executor executor, logger logger.Logger) *Poller {
	return &Poller{
		cloner:   cloner,
		parser:   parser,
		executor: executor,
		logger:   logger,
	}
}

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

func (p Poller) poll(vcs domain.VCS) error {
	sourceCode, err := p.cloner.CloneRepository(vcs.URL)
	if err != nil {
		return err
	}
	defer sourceCode.Close()

	yaml, err := p.findPipeline(sourceCode)
	if err != nil {
		return err
	}

	pipeline, err := p.parser.ParsePipeline(yaml)
	if err != nil {
		return err
	}

	for _, step := range pipeline.Steps {
		err = p.executor.Execute(context.Background(), step, sourceCode)
		if err != nil {
			return err
		}
	}

	return nil
}

func (Poller) findPipeline(archive io.Reader) ([]byte, error) {
	const yamlFilename = "ci.yaml"

	tarReader := tar.NewReader(archive)

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
