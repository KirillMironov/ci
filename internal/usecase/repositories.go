package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/rs/xid"
)

type Repositories struct {
	storage repositoriesStorage
	add     chan<- domain.Repository
	remove  chan<- string
}

type repositoriesStorage interface {
	Create(context.Context, domain.Repository) error
	Delete(ctx context.Context, id string) error
	GetAll(context.Context) ([]domain.Repository, error)
	GetById(ctx context.Context, id string) (domain.Repository, error)
	GetByURL(ctx context.Context, url string) (domain.Repository, error)
}

func NewRepositories(storage repositoriesStorage, add chan<- domain.Repository, remove chan<- string) *Repositories {
	return &Repositories{
		storage: storage,
		add:     add,
		remove:  remove,
	}
}

func (r Repositories) Add(ctx context.Context, repo domain.Repository) error {
	_, err := r.storage.GetByURL(ctx, repo.URL)
	if !errors.Is(err, domain.ErrNotFound) {
		return fmt.Errorf("repository with url %q already exists", repo.URL)
	}

	repo.Id = xid.New().String()

	err = r.storage.Create(ctx, repo)
	if err != nil {
		return fmt.Errorf("failed to create repository: %w", err)
	}

	r.add <- repo
	return nil
}

func (r Repositories) Delete(ctx context.Context, id string) error {
	r.remove <- id

	err := r.storage.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete repository: %w", err)
	}

	return nil
}

func (r Repositories) GetAll(ctx context.Context) ([]domain.Repository, error) {
	repos, err := r.storage.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get repositories: %w", err)
	}

	return repos, nil
}

func (r Repositories) GetById(ctx context.Context, id string) (domain.Repository, error) {
	repo, err := r.storage.GetById(ctx, id)
	if err != nil {
		return domain.Repository{}, fmt.Errorf("failed to get repository: %w", err)
	}

	return repo, nil
}

func (r Repositories) GetByURL(ctx context.Context, url string) (domain.Repository, error) {
	repo, err := r.storage.GetByURL(ctx, url)
	if err != nil {
		return domain.Repository{}, fmt.Errorf("failed to get repository: %w", err)
	}

	return repo, nil
}
