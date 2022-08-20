package service

import (
	"context"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/KirillMironov/ci/pkg/logger"
	"sync"
)

// Scheduler used to schedule repositories polling.
type Scheduler struct {
	add                 chan domain.Repository
	remove              <-chan string
	activePolling       map[string]context.CancelFunc
	once                sync.Once
	poller              poller
	repositoriesUsecase repositoriesUsecase
	logger              logger.Logger
}

type (
	poller interface {
		AddRepository(context.Context, domain.Repository)
	}
	repositoriesUsecase interface {
		GetAll() ([]domain.Repository, error)
	}
)

func NewScheduler(add chan domain.Repository, remove <-chan string, poller poller, ru repositoriesUsecase,
	logger logger.Logger) *Scheduler {
	return &Scheduler{
		add:                 add,
		remove:              remove,
		activePolling:       make(map[string]context.CancelFunc),
		poller:              poller,
		repositoriesUsecase: ru,
		logger:              logger,
	}
}

// Start listens for repositories additions and deletions and starts polling.
func (s *Scheduler) Start(ctx context.Context) {
	s.once.Do(func() {
		go func() {
			repos, err := s.repositoriesUsecase.GetAll()
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
			pollCtx, cancel := context.WithCancel(ctx)
			s.activePolling[repo.Id] = cancel
			s.poller.AddRepository(pollCtx, repo)
		case id := <-s.remove:
			if cancel, ok := s.activePolling[id]; ok {
				cancel()
				delete(s.activePolling, id)
			}
		}
	}
}
