package service

import (
	"context"
	"errors"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/KirillMironov/ci/pkg/logger"
	"github.com/rs/xid"
	"sync"
)

// Scheduler used to schedule repositories polling.
type Scheduler struct {
	add                 chan domain.Repository
	remove              chan string
	activePolling       map[string]context.CancelFunc
	once                sync.Once
	poller              poller
	repositoriesStorage domain.RepositoriesStorage
	logger              logger.Logger
}

type poller interface {
	AddRepository(context.Context, domain.Repository)
}

func NewScheduler(poller poller, rs domain.RepositoriesStorage, logger logger.Logger) *Scheduler {
	return &Scheduler{
		add:                 make(chan domain.Repository),
		remove:              make(chan string),
		activePolling:       make(map[string]context.CancelFunc),
		poller:              poller,
		repositoriesStorage: rs,
		logger:              logger,
	}
}

// Start listens for repositories additions and deletions and starts polling.
func (s *Scheduler) Start(ctx context.Context) {
	s.once.Do(func() {
		go func() {
			repos, err := s.repositoriesStorage.GetAll()
			if err != nil {
				s.logger.Error(err)
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
			_, err := s.repositoriesStorage.GetByURL(repo.URL)
			if !errors.Is(err, domain.ErrNotFound) {
				s.logger.Errorf("repository with url %q already exists", repo.URL)
				continue
			}

			repo.Id = xid.New().String()

			err = s.repositoriesStorage.Create(repo)
			if err != nil {
				s.logger.Errorf("failed to create repository: %v", err)
				continue
			}

			pollCtx, cancel := context.WithCancel(ctx)
			s.activePolling[repo.Id] = cancel
			s.poller.AddRepository(pollCtx, repo)
		case id := <-s.remove:
			if cancel, ok := s.activePolling[id]; ok {
				cancel()
				delete(s.activePolling, id)
			}

			err := s.repositoriesStorage.Delete(id)
			if err != nil {
				s.logger.Errorf("failed to delete repository: %v", err)
			}
		}
	}
}

func (s *Scheduler) Add(repo domain.Repository) {
	s.add <- repo
}

func (s *Scheduler) Remove(id string) {
	s.remove <- id
}
