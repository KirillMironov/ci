package usecase

import (
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
	Create(domain.Repository) error
	Delete(id string) error
	GetAll() ([]domain.Repository, error)
	GetById(id string) (domain.Repository, error)
	GetByURL(url string) (domain.Repository, error)
}

func NewRepositories(storage repositoriesStorage, add chan<- domain.Repository, remove chan<- string) *Repositories {
	return &Repositories{
		storage: storage,
		add:     add,
		remove:  remove,
	}
}

func (r Repositories) Add(repo domain.Repository) error {
	_, err := r.storage.GetByURL(repo.URL)
	if !errors.Is(err, domain.ErrNotFound) {
		return fmt.Errorf("repository with url %q already exists", repo.URL)
	}

	repo.Id = xid.New().String()

	err = r.storage.Create(repo)
	if err != nil {
		return fmt.Errorf("failed to create repository: %w", err)
	}

	r.add <- repo
	return nil
}

func (r Repositories) Delete(id string) error {
	r.remove <- id

	err := r.storage.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete repository: %w", err)
	}

	return nil
}

func (r Repositories) GetAll() ([]domain.Repository, error) {
	repos, err := r.storage.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all repositories: %w", err)
	}

	return repos, nil
}

func (r Repositories) GetById(id string) (domain.Repository, error) {
	repo, err := r.storage.GetById(id)
	if err != nil {
		return domain.Repository{}, fmt.Errorf("failed to get repository by id: %w", err)
	}

	return repo, nil
}

func (r Repositories) GetByURL(url string) (domain.Repository, error) {
	repo, err := r.storage.GetByURL(url)
	if err != nil {
		return domain.Repository{}, fmt.Errorf("failed to get repository by url: %w", err)
	}

	return repo, nil
}
