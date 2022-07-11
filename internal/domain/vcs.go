package domain

import "time"

// VCS represents a version control system.
type VCS struct {
	URL             string
	PollingInterval time.Duration
}
