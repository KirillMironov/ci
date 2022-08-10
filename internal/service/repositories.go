package service

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"github.com/KirillMironov/ci/internal/domain"
)

// Repositories used to work with repositories storage.
type Repositories struct {
	storage repositoriesStorage
}

type repositoriesStorage interface {
	Create(domain.Repository) error
	Delete(id string) error
	GetAll() ([]domain.Repository, error)
	GetById(id string) (domain.Repository, error)
}

func NewRepositories(storage repositoriesStorage) *Repositories {
	return &Repositories{storage: storage}
}

func (r Repositories) GetOrCreate(repo domain.Repository) (savedRepo domain.Repository, err error) {
	savedRepo, err = r.storage.GetById(repo.Id)
	if errors.Is(err, domain.ErrNotFound) {
		repo.Id = generateIdFromURL(repo.URL)
		return repo, r.storage.Create(repo)
	}
	return savedRepo, err
}

func (r Repositories) Update(repo domain.Repository) error {
	return r.storage.Create(repo)
}

func (r Repositories) Delete(id string) error {
	return r.storage.Delete(id)
}

func (r Repositories) GetAll() ([]domain.Repository, error) {
	return r.storage.GetAll()
}

func (r Repositories) GetById(id string) (domain.Repository, error) {
	return r.storage.GetById(id)
}

func (r Repositories) GetBuilds(id string) ([]domain.Build, error) {
	repo, err := r.storage.GetById(id)
	if err != nil {
		return nil, err
	}
	return repo.Builds, nil
}

func generateIdFromURL(url string) string {
	h := sha1.New()
	h.Write([]byte(url))
	return hex.EncodeToString(h.Sum(nil))
}
