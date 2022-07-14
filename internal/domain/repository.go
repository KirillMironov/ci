package domain

import "time"

// Repository represents a source code repository.
type Repository struct {
	URL             string
	PollingInterval time.Duration
}
