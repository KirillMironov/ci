package service

import (
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/KirillMironov/ci/pkg/mock"
	"testing"
	"time"
)

func TestPoller_Poll(t *testing.T) {
	var poller = NewPoller(mock.Cloner{}, mock.Parser{}, mock.Executor{}, mock.Logger{})

	poller.Start(domain.VCS{
		URL:             "https://github.com/KirillMironov/ci",
		PollingInterval: time.Second * 2,
	})
}
