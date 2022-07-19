package service

import (
	"context"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/KirillMironov/ci/pkg/logger"
	"sync"
)

type Scheduler struct {
	add          chan domain.Repository
	once         sync.Once
	poller       poller
	repositories repositories
	logger       logger.Logger
}

type (
	// poller is a service that can poll a source code repository.
	poller interface {
		Poll(repo domain.Repository, prevHash string) (newHash string, err error)
	}
	// repositories stores information about source code repositories.
	repositories interface {
		Put(domain.Repository) error
		GetAll() ([]domain.Repository, error)
	}
)

// NewScheduler creates a new Scheduler.
func NewScheduler(add chan domain.Repository, poller poller, repositories repositories,
	logger logger.Logger) *Scheduler {
	return &Scheduler{
		add:          add,
		poller:       poller,
		repositories: repositories,
		logger:       logger,
	}
}

// Start listens for new repositories.
func (s *Scheduler) Start(ctx context.Context) {
	s.once.Do(func() {
		go func() {
			repos, err := s.repositories.GetAll()
			if err != nil {
				s.logger.Errorf("failed to get repositories: %v", err)
				return
			}

			for _, repo := range repos {
				s.add <- repo
			}
		}()
	})

	for {
		select {
		case <-ctx.Done():
			s.logger.Infof("scheduler stopped: %v", ctx.Err())
			return
		case repo := <-s.add:
			s.logger.Infof("starting polling %s", repo.URL)
			go s.poll(repo)
		}
	}
}

// poll starts polling a source code repository.
func (s *Scheduler) poll(repo domain.Repository) {
	err := s.repositories.Put(repo)
	if err != nil {
		s.logger.Errorf("failed to save repository %s: %v", repo.URL, err)
		return
	}

	repo.Hash, err = s.poller.Poll(repo, repo.Hash)
	if err != nil {
		s.logger.Errorf("failed to poll %s: %v", repo.URL, err)
		return
	}

	s.add <- repo
}
