package domain

import "time"

// Repository represents a source code repository.
type Repository struct {
	URL             string
	Hash            string
	Branch          string
	PollingInterval time.Duration
}

// RepositoryURL used to identify a repository.
type RepositoryURL string
