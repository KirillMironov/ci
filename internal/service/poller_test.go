package service

import (
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestPoller_Poll(t *testing.T) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "01|02 15:04:05.000",
	})

	cli, err := client.NewClientWithOpts()
	if err != nil {
		logger.Fatal(err)
	}
	defer cli.Close()

	var (
		cloner   = &Cloner{}
		parser   = &Parser{}
		executor = NewDockerExecutor(cli)
		poller   = NewPoller(cloner, parser, executor, logger)
	)

	poller.Start(domain.VCS{
		URL:             "",
		PollingInterval: time.Second * 5,
	})
}
