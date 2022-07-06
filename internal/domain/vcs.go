package domain

import "time"

type VCS struct {
	URL             string
	PollingInterval time.Duration
}
