package domain

import (
	"context"
	"github.com/KirillMironov/ci/pkg/duration"
)

type Repository struct {
	Id              string            `json:"id"`
	URL             string            `json:"url"`
	Branch          string            `json:"branch"`
	PollingInterval duration.Duration `json:"polling_interval"`
}

type RepositoriesUsecase interface {
	Add(context.Context, Repository) error
	Delete(ctx context.Context, id string) error
	GetAll(context.Context) ([]Repository, error)
	GetById(ctx context.Context, id string) (Repository, error)
	GetByURL(ctx context.Context, url string) (Repository, error)
}
