package service

import (
	"context"
	"errors"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/KirillMironov/ci/pkg/logger"
	"os"
	"path/filepath"
	"time"
)

// Poller used to poll repositories and run builds.
type Poller struct {
	poll          chan domain.Repository
	run           chan<- RunRequest
	ciFilename    string
	cloner        cloner
	parser        parser
	buildsUsecase domain.BuildsUsecase
	logger        logger.Logger
}

type (
	cloner interface {
		GetLatestCommitHash(domain.Repository) (string, error)
		CloneRepository(repo domain.Repository, targetHash string) (srcCodePath string, err error)
	}
	parser interface {
		ParsePipeline(b []byte) (domain.Pipeline, error)
	}
)

func NewPoller(run chan<- RunRequest, ciFilename string, cloner cloner, parser parser, bu domain.BuildsUsecase,
	logger logger.Logger) *Poller {
	return &Poller{
		run:           run,
		poll:          make(chan domain.Repository),
		ciFilename:    ciFilename,
		cloner:        cloner,
		parser:        parser,
		buildsUsecase: bu,
		logger:        logger,
	}
}

// Start listens on the poll channel and runs a build if the repository contains a new commit.
func (p Poller) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			p.logger.Infof("poller stopped: %v", ctx.Err())
			return
		case repo := <-p.poll:
			latestHash, err := p.cloner.GetLatestCommitHash(repo)
			if err != nil {
				p.logger.Errorf("failed to get latest commit hash: %v", err)
				continue
			}

			builds, err := p.buildsUsecase.GetAllByRepoId(context.Background(), repo.Id)
			if err != nil && !errors.Is(err, domain.ErrNotFound) {
				p.logger.Errorf("failed to get builds: %v", err)
				continue
			}
			if len(builds) > 0 && latestHash == builds[len(builds)-1].Commit.Hash {
				continue
			}

			srcCodePath, err := p.cloner.CloneRepository(repo, latestHash)
			if err != nil {
				p.logger.Errorf("failed to clone repository: %v", err)
				continue
			}

			data, err := os.ReadFile(filepath.Join(srcCodePath, p.ciFilename))
			if err != nil {
				p.logger.Errorf("failed to read ci file: %v", err)
				continue
			}

			pipeline, err := p.parser.ParsePipeline(data)
			if err != nil {
				p.logger.Errorf("failed to parse pipeline: %v", err)
				continue
			}

			p.run <- RunRequest{
				repoId:      repo.Id,
				commit:      domain.Commit{Hash: latestHash},
				pipeline:    pipeline,
				srcCodePath: srcCodePath,
			}
		}
	}
}

// AddRepository sends the repository to the poll channel at regular intervals.
func (p Poller) AddRepository(ctx context.Context, repo domain.Repository) {
	go func() {
		timer := time.NewTimer(repo.PollingInterval.Duration())

		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				p.poll <- repo
				timer.Reset(repo.PollingInterval.Duration())
			}
		}
	}()
}
