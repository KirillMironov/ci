package service

import (
	"context"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/KirillMironov/ci/pkg/logger"
	"time"
)

type Poller struct {
	Interval time.Duration
	parser   parser
	executor executor
	logger   logger.Logger
}

type parser interface {
	ParsePipeline(str string) (domain.Pipeline, error)
}

type executor interface {
	Execute(ctx context.Context, step domain.Step) error
}

func NewPoller(interval time.Duration, parser parser, executor executor, logger logger.Logger) *Poller {
	return &Poller{
		Interval: interval,
		parser:   parser,
		executor: executor,
		logger:   logger,
	}
}

func (p Poller) Poll() {

}
