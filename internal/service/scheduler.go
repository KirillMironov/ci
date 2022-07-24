package service

import (
	"context"
	"errors"
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
	logs         logs
	logger       logger.Logger
}

type (
	// poller is a service that can poll a source code repository.
	poller interface {
		Poll(ctx context.Context, repo domain.Repository, prevHash string) (newHash string, logs []byte, err error)
	}
	// repositories stores information about source code repositories.
	repositories interface {
		Put(domain.Repository) error
		Delete(domain.RepositoryURL) error
		GetAll() ([]domain.Repository, error)
		GetByURL(url string) (domain.Repository, error)
	}
	logs interface {
		Put(log domain.Log) (id int, err error)
		GetById(id int) (domain.Log, error)
	}
)

// NewScheduler creates a new Scheduler.
func NewScheduler(poller poller, repositories repositories, logs logs, logger logger.Logger) *Scheduler {
	return &Scheduler{
		put:          make(chan domain.Repository),
		delete:       make(chan domain.RepositoryURL),
		polling:      make(map[domain.RepositoryURL]context.CancelFunc),
		poller:       poller,
		repositories: repositories,
		logs:         logs,
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
			s.logger.Infof("starting polling %s", repo.URL)
			if cancel, ok := s.polling[domain.RepositoryURL(repo.URL)]; ok {
				cancel()
			}
			pollCtx, cancel := context.WithCancel(ctx)
			s.polling[domain.RepositoryURL(repo.URL)] = cancel
			go s.poll(pollCtx, repo)
		case repoURL := <-s.delete:
			s.logger.Infof("removing repository %s", repoURL)
			if cancel, ok := s.polling[repoURL]; ok {
				cancel()
			}
			err := s.repositories.Delete(repoURL)
			if err != nil {
				s.logger.Errorf("failed to remove repository %s: %v", repoURL, err)
			}
		}
	}
}

func (s *Scheduler) Put(repo domain.Repository) {
	s.put <- repo
}

func (s *Scheduler) Delete(repoURL domain.RepositoryURL) {
	s.delete <- repoURL
}

// poll starts polling a source code repository.
func (s *Scheduler) poll(ctx context.Context, repo domain.Repository) {
	savedRepo, err := s.repositories.GetByURL(repo.URL)
	if err != nil && !errors.Is(err, domain.ErrRepoNotFound) {
		s.logger.Errorf("failed to get repository %s: %v", repo.URL, err)
		return
	}
	repo.Builds = savedRepo.Builds

	var prevHash string
	if repo.Builds != nil {
		prevHash = repo.Builds[len(repo.Builds)-1].Commit.Hash
	}

	hash, buildLogs, err := s.poller.Poll(ctx, repo, prevHash)
	if err != nil {
		s.logger.Errorf("failed to poll %s: %v", repo.URL, err)
		return
	}

	logId, err := s.logs.Put(domain.Log{Data: buildLogs})
	if err != nil {
		s.logger.Errorf("failed to put log: %v", err)
		return
	}

	repo.Builds = append(repo.Builds, domain.Build{
		Commit: domain.Commit{Hash: hash},
		LogId:  logId,
	})

	err = s.repositories.Put(repo)
	if err != nil {
		s.logger.Errorf("failed to put repository %s: %v", repo.URL, err)
		return
	}

	s.put <- repo
}
