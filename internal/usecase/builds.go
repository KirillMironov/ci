package usecase

import (
	"context"
	"fmt"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/rs/xid"
)

type Builds struct {
	storage buildsStorage
}

type buildsStorage interface {
	Create(ctx context.Context, build domain.Build, repoId string) error
	Delete(ctx context.Context, id, repoId string) error
	GetAllByRepoId(ctx context.Context, repoId string) ([]domain.Build, error)
	GetById(ctx context.Context, id, repoId string) (domain.Build, error)
}

func NewBuilds(storage buildsStorage) *Builds {
	return &Builds{storage: storage}
}

func (b Builds) Create(ctx context.Context, build domain.Build, repoId string) error {
	build.Id = xid.New().String()

	err := b.storage.Create(ctx, build, repoId)
	if err != nil {
		return fmt.Errorf("failed to create build: %w", err)
	}

	return nil
}

func (b Builds) Delete(ctx context.Context, id, repoId string) error {
	err := b.storage.Delete(ctx, id, repoId)
	if err != nil {
		return fmt.Errorf("failed to delete build: %w", err)
	}

	return nil
}

func (b Builds) GetAllByRepoId(ctx context.Context, repoId string) ([]domain.Build, error) {
	builds, err := b.storage.GetAllByRepoId(ctx, repoId)
	if err != nil {
		return nil, fmt.Errorf("failed to get builds by repoId: %w", err)
	}

	return builds, nil
}

func (b Builds) GetById(ctx context.Context, id, repoId string) (domain.Build, error) {
	build, err := b.storage.GetById(ctx, id, repoId)
	if err != nil {
		return domain.Build{}, fmt.Errorf("failed to get build by id: %w", err)
	}

	return build, nil
}
