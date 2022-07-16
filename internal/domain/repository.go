package domain

import "time"

// Repository represents a source code repository.
type Repository struct {
	URL             string
	Hash            string
	Branch          string
	PollingInterval time.Duration
}
