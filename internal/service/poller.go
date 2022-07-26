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
	cloner       cloner
	runner       runner
	repositories repositories
	logger       logger.Logger
}

type (
	runner interface {
		Run(ctx context.Context, repo domain.Repository, targetHash string) (logId int, err error)
	}
	// repositories stores information about source code repositories.
	repositories interface {
		Save(domain.Repository) error
		Delete(domain.RepositoryURL) error
		GetAll() ([]domain.Repository, error)
		GetByURL(url string) (domain.Repository, error)
	}
)

func NewPoller(cloner cloner, runner runner, repositories repositories, logger logger.Logger) *Poller {
	return &Poller{
		poll:         make(chan domain.Repository),
		cloner:       cloner,
		runner:       runner,
		repositories: repositories,
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

			logId, err := p.runner.Run(ctx, repo, latestHash)
			if err != nil {
				p.logger.Errorf("failed to run: %v", err)
				continue
			}

			repo.Builds = append(repo.Builds, domain.Build{
				Commit: domain.Commit{Hash: latestHash},
				LogId:  logId,
			})

			err = p.repositories.Save(repo)
			if err != nil {
				p.logger.Errorf("failed to save repository: %v", err)
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
