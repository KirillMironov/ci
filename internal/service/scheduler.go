package service

import (
	"context"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/KirillMironov/ci/pkg/logger"
	"sync"
)

// Scheduler used to schedule repositories polling.
type Scheduler struct {
	put    chan domain.Repository
	delete chan string
	// activePolling used to cancel a repository polling if it's already running.
	activePolling       map[string]context.CancelFunc
	once                sync.Once
	poller              poller
	repositoriesService domain.RepositoriesService
	logger              logger.Logger
}

type poller interface {
	AddRepository(context.Context, domain.Repository)
}

func NewScheduler(poller poller, repositoriesService domain.RepositoriesService, logger logger.Logger) *Scheduler {
	return &Scheduler{
		put:                 make(chan domain.Repository),
		delete:              make(chan string),
		activePolling:       make(map[string]context.CancelFunc),
		poller:              poller,
		repositoriesService: repositoriesService,
		logger:              logger,
	}
}

func (s *Scheduler) Put(repo domain.Repository) {
	s.put <- repo
}

func (s *Scheduler) Delete(id string) {
	s.delete <- id
}

// Start listens for repositories additions and deletions and starts polling.
func (s *Scheduler) Start(ctx context.Context) {
	s.once.Do(func() {
		go func() {
			repos, err := s.repositoriesService.GetAll()
			if err != nil {
				s.logger.Errorf("failed to get all saved repositories: %v", err)
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
			s.cancelPolling(repo.Id)

			repo, err := s.repositoriesService.GetOrCreate(repo)
			if err != nil {
				s.logger.Errorf("failed to create repository: %v", err)
				continue
			}

			pollCtx, cancel := context.WithCancel(ctx)
			s.activePolling[repo.Id] = cancel
			s.poller.AddRepository(pollCtx, repo)
		case repoURL := <-s.delete:
			s.cancelPolling(repoURL)

			err := s.repositoriesService.Delete(repoURL)
			if err != nil {
				s.logger.Errorf("failed to delete repository %s: %v", repoURL, err)
			}
		}
	}
}

func (s *Scheduler) cancelPolling(id string) {
	if cancel, ok := s.activePolling[id]; ok {
		cancel()
		delete(s.activePolling, id)
	}
}
