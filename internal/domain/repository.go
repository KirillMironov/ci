package domain

import (
	"github.com/KirillMironov/ci/pkg/duration"
	"time"
)

type Repository struct {
	Id              string            `json:"id"`
	URL             string            `json:"url"`
	Branch          string            `json:"branch"`
	PollingInterval duration.Duration `json:"polling_interval"`
	CreatedAt       time.Time         `json:"created_at"`
}

type RepositoriesUsecase interface {
	Add(Repository) error
	Delete(id string) error
	GetAll() ([]Repository, error)
	GetById(id string) (Repository, error)
	GetByURL(url string) (Repository, error)
}
