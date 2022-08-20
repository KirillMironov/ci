package usecase

import (
	"fmt"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/rs/xid"
)

type Builds struct {
	storage buildsStorage
}

type buildsStorage interface {
	Create(domain.Build) error
	Delete(id string) error
	GetAllByRepoId(repoId string) ([]domain.Build, error)
	GetById(id string) (domain.Build, error)
}

func NewBuilds(storage buildsStorage) *Builds {
	return &Builds{storage: storage}
}

func (b Builds) Create(build domain.Build) error {
	build.Id = xid.New().String()

	err := b.storage.Create(build)
	if err != nil {
		return fmt.Errorf("failed to create build: %w", err)
	}

	return nil
}

func (b Builds) Delete(id string) error {
	err := b.storage.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete build: %w", err)
	}

	return nil
}

func (b Builds) GetAllByRepoId(repoId string) ([]domain.Build, error) {
	builds, err := b.storage.GetAllByRepoId(repoId)
	if err != nil {
		return nil, fmt.Errorf("failed to get builds by repo id: %w", err)
	}

	return builds, nil
}

func (b Builds) GetById(id string) (domain.Build, error) {
	build, err := b.storage.GetById(id)
	if err != nil {
		return domain.Build{}, fmt.Errorf("failed to get build by id: %w", err)
	}

	return build, nil
}
