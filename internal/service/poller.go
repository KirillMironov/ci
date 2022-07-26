package service

import (
	"context"
	"errors"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/KirillMironov/ci/pkg/logger"
	"time"
)

// Poller is a service that can poll a source code repository.
type Poller struct {
	poll         chan domain.Repository
	ciFilename   string
	runner       runner
	cloner       cloner
	finder       finder
	parser       parser
	repositories repositories
	logs         logs
	logger       logger.Logger
}

type (
	runner interface {
		Run(ctx context.Context, pipeline domain.Pipeline, srcCodeArchivePath string) (logs []byte, err error)
	}
	cloner interface {
		GetLatestCommitHash(url, branch string) (string, error)
		CloneRepository(url, branch, hash string) (archivePath string, removeArchive func(), err error)
	}
	finder interface {
		FindFile(filename, archivePath string) ([]byte, error)
	}
	parser interface {
		ParsePipeline(b []byte) (domain.Pipeline, error)
	}
	repositories interface {
		Save(domain.Repository) error
		Delete(domain.RepositoryURL) error
		GetAll() ([]domain.Repository, error)
		GetByURL(url string) (domain.Repository, error)
	}
	logs interface {
		Save(domain.Log) (id int, err error)
	}
)

func NewPoller(ciFilename string, runner runner, cloner cloner, finder finder, parser parser, repositories repositories,
	logs logs, logger logger.Logger) *Poller {
	return &Poller{
		poll:         make(chan domain.Repository),
		ciFilename:   ciFilename,
		runner:       runner,
		cloner:       cloner,
		finder:       finder,
		parser:       parser,
		repositories: repositories,
		logs:         logs,
		logger:       logger,
	}
}

func (p Poller) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			p.logger.Infof("poller stopped: %v", ctx.Err())
			return
		case repo := <-p.poll:
			latestHash, err := p.cloner.GetLatestCommitHash(repo.URL, repo.Branch)
			if err != nil {
				p.logger.Errorf("failed to get latest commit hash: %v", err)
				continue
			}
			if repo.Builds != nil && latestHash == repo.Builds[len(repo.Builds)-1].Commit.Hash {
				continue
			}

			err = p.build(ctx, repo, latestHash)
			if err != nil {
				p.logger.Errorf("failed to build repository: %q; %v", repo.URL, err)
			}
		}
	}
}

func (p Poller) AddRepository(ctx context.Context, repo domain.Repository) {
	go func() {
		timer := time.NewTimer(repo.PollingInterval)

		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				savedRepo, err := p.repositories.GetByURL(repo.URL)
				if err != nil && !errors.Is(err, domain.ErrRepoNotFound) {
					p.logger.Errorf("failed to get saved repository: %v", err)
				}
				repo.Builds = savedRepo.Builds

				p.poll <- repo
				timer.Reset(repo.PollingInterval)
			}
		}
	}()
}

func (p Poller) build(ctx context.Context, repo domain.Repository, targetHash string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	archivePath, removeArchive, err := p.cloner.CloneRepository(repo.URL, repo.Branch, targetHash)
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

	logId, err := p.logs.Save(domain.Log{Data: pipelineLogs})
	if err != nil {
		return err
	}

	repo.Builds = append(repo.Builds, domain.Build{
		Commit: domain.Commit{Hash: targetHash},
		LogId:  logId,
	})

	return p.repositories.Save(repo)
}
