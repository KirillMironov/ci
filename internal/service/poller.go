package service

import (
	"context"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/KirillMironov/ci/pkg/logger"
	"io/ioutil"
	"os"
	"time"
)

type Poller struct {
	cloner   cloner
	parser   parser
	executor executor
	logger   logger.Logger
}

type cloner interface {
	CloneRepository(url string) (dir string, err error)
}

type parser interface {
	ParsePipeline(str string) (domain.Pipeline, error)
}

type executor interface {
	Execute(ctx context.Context, step domain.Step, sourceCodePath string) error
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
	const yamlFilename = "/ci.yaml"

	sourceCodePath, err := p.cloner.CloneRepository(vcs.URL)
	if err != nil {
		return err
	}
	defer os.RemoveAll(sourceCodePath)

	yaml, err := ioutil.ReadFile(sourceCodePath + yamlFilename)
	if err != nil {
		return err
	}

	pipeline, err := p.parser.ParsePipeline(string(yaml))
	if err != nil {
		return err
	}

	for _, step := range pipeline.Steps {
		err = p.executor.Execute(context.Background(), step, sourceCodePath)
		if err != nil {
			return err
		}
	}

	return nil
}
