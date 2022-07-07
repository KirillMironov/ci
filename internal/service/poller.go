package service

import (
	"context"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/KirillMironov/ci/pkg/logger"
	"time"
)

type Poller struct {
	cloner   cloner
	parser   parser
	executor executor
	logger   logger.Logger
}

type cloner interface {
	Clone(url string) (dir string, err error)
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
		err := p.poll()
		if err != nil {
			p.logger.Error(err)
		}
		timer.Reset(vcs.PollingInterval)
	}
}

func (p Poller) poll() error {
	time.Sleep(time.Second * 10)
	return nil
}
