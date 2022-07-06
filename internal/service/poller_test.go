package service

import (
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/KirillMironov/ci/pkg/mock"
	"testing"
	"time"
)

func TestPoller_Poll(t *testing.T) {
	var poller = NewPoller(mock.Parser{}, mock.Executor{}, mock.Logger{})

	poller.Poll(domain.VCS{
		URL:             "https://github.com/KirillMironov/ci",
		PollingInterval: time.Second * 2,
	})
}
