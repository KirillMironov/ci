package service

import (
	"context"
	"errors"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/KirillMironov/ci/pkg/logger"
	"time"
)

// Poller used to poll repositories and run builds.
type Poller struct {
	poll                chan domain.Repository
	ciFilename          string
	runner              runner
	cloner              cloner
	finder              finder
	parser              parser
	repositoriesService domain.RepositoriesService
	logsService         logsService
	logger              logger.Logger
}

type (
	runner interface {
		Run(ctx context.Context, pipeline domain.Pipeline, srcCodeArchivePath string) (logs []byte, err error)
	}
	cloner interface {
		GetLatestCommitHash(domain.Repository) (string, error)
		CloneRepository(domain.Repository, string) (archivePath string, removeArchive func(), err error)
	}
	finder interface {
		FindFile(filename, archivePath string) ([]byte, error)
	}
	parser interface {
		ParsePipeline(b []byte) (domain.Pipeline, error)
	}
	logsService interface {
		Create(domain.Log) (id int, err error)
	}
)

func NewPoller(ciFilename string, runner runner, cloner cloner, finder finder, parser parser,
	repositoriesService domain.RepositoriesService, logsService logsService, logger logger.Logger) *Poller {
	return &Poller{
		poll:                make(chan domain.Repository),
		ciFilename:          ciFilename,
		runner:              runner,
		cloner:              cloner,
		finder:              finder,
		parser:              parser,
		repositoriesService: repositoriesService,
		logsService:         logsService,
		logger:              logger,
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
			if repo.Builds != nil && latestHash == repo.Builds[len(repo.Builds)-1].Commit.Hash {
				continue
			}

			err = p.build(ctx, repo, latestHash)
			if err != nil {
				p.logger.Errorf("failed to build: %q; %v", repo.URL, err)
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
				builds, err := p.repositoriesService.GetBuilds(repo.Id)
				if err != nil && !errors.Is(err, domain.ErrNotFound) {
					p.logger.Errorf("failed to get builds: %v", err)
					timer.Reset(repo.PollingInterval.Duration())
					continue
				}

				repo.Builds = builds
				p.poll <- repo

				timer.Reset(repo.PollingInterval.Duration())
			}
		}
	}()
}

func (p Poller) build(ctx context.Context, repo domain.Repository, targetHash string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	archivePath, removeArchive, err := p.cloner.CloneRepository(repo, targetHash)
	if err != nil {
		return err
	}
	defer removeArchive()

	yaml, err := p.finder.FindFile(p.ciFilename, archivePath)
	if err != nil {
		return err
	}

	pipeline, err := p.parser.ParsePipeline(yaml)
	if err != nil {
		return err
	}

	pipelineLogs, err := p.runner.Run(ctx, pipeline, archivePath)
	if err != nil {
		return err
	}

	logId, err := p.logsService.Create(domain.Log{Data: pipelineLogs})
	if err != nil {
		return err
	}

	repo.Builds = append(repo.Builds, domain.Build{
		Commit: domain.Commit{Hash: targetHash},
		LogId:  logId,
	})

	return p.repositoriesService.Update(repo)
}
