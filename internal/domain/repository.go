package domain

import "github.com/KirillMironov/ci/pkg/duration"

// Repository represents a source code repository.
type Repository struct {
	Id              string            `json:"id"`
	URL             string            `json:"url"`
	Branch          string            `json:"branch"`
	PollingInterval duration.Duration `json:"polling_interval"`
	Builds          []Build           `json:"builds"`
}

// RepositoryURL used to identify a repository.
type RepositoryURL string
