package service

import (
	"context"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/KirillMironov/ci/pkg/logger"
	"sync"
)

type Scheduler struct {
	put    chan domain.Repository
	delete chan domain.RepositoryURL
	// polling used to cancel polling of a repository, if it's already running.
	polling      map[domain.RepositoryURL]context.CancelFunc
	once         sync.Once
	poller       poller
	repositories repositories
	logger       logger.Logger
}

type poller interface {
	AddRepository(context.Context, domain.Repository)
}

// NewScheduler creates a new Scheduler.
func NewScheduler(poller poller, repositories repositories, logger logger.Logger) *Scheduler {
	return &Scheduler{
		put:          make(chan domain.Repository),
		delete:       make(chan domain.RepositoryURL),
		polling:      make(map[domain.RepositoryURL]context.CancelFunc),
		poller:       poller,
		repositories: repositories,
		logger:       logger,
	}
}

func (s *Scheduler) Put(repo domain.Repository) {
	s.put <- repo
}

func (s *Scheduler) Delete(repoURL domain.RepositoryURL) {
	s.delete <- repoURL
}

// Start listens for new repositories.
func (s *Scheduler) Start(ctx context.Context) {
	s.once.Do(func() {
		go func() {
			repos, err := s.repositories.GetAll()
			if err != nil {
				s.logger.Errorf("failed to get saved repositories: %v", err)
				return
			}

			for _, repo := range repos {
				s.put <- repo
			}
		}()
	})

	for {
		select {
		case <-ctx.Done():
			s.logger.Infof("scheduler stopped: %v", ctx.Err())
			return
		case repo := <-s.put:
			s.cancelPolling(domain.RepositoryURL(repo.URL))

			pollCtx, cancel := context.WithCancel(ctx)
			s.polling[domain.RepositoryURL(repo.URL)] = cancel

			s.poller.AddRepository(pollCtx, repo)
		case repoURL := <-s.delete:
			s.cancelPolling(repoURL)

			err := s.repositories.Delete(repoURL)
			if err != nil {
				s.logger.Errorf("failed to delete repository %s: %v", repoURL, err)
			}
		}
	}
}

func (s *Scheduler) cancelPolling(repoURL domain.RepositoryURL) {
	if cancel, ok := s.polling[repoURL]; ok {
		cancel()
		delete(s.polling, repoURL)
	}
}
