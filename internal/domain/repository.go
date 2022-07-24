package domain

import "time"

// Repository represents a source code repository.
type Repository struct {
	URL             string
	Branch          string
	PollingInterval time.Duration
	Builds          []Build
}

// RepositoryURL used to identify a repository.
type RepositoryURL string
