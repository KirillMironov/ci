package domain

import "github.com/KirillMironov/ci/pkg/duration"

type Repository struct {
	Id              string            `json:"id"`
	URL             string            `json:"url"`
	Branch          string            `json:"branch"`
	PollingInterval duration.Duration `json:"polling_interval"`
	Builds          []Build           `json:"builds"`
}

type RepositoriesService interface {
	GetOrCreate(Repository) (Repository, error)
	Update(Repository) error
	Delete(id string) error
	GetAll() ([]Repository, error)
	GetBuilds(id string) ([]Build, error)
}
