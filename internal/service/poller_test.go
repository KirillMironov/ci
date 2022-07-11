package service

import (
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPoller_poll(t *testing.T) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "01|02 15:04:05.000",
	})

	cli, err := client.NewClientWithOpts()
	require.NoError(t, err)
	defer cli.Close()

	var (
		cloner   = &Cloner{}
		archiver = &Archiver{}
		parser   = &Parser{}
		executor = NewDockerExecutor(cli)
		poller   = NewPoller(".ci.yaml", cloner, archiver, parser, executor, logger)
	)

	poller.Start(domain.VCS{
		URL:             "",
		PollingInterval: time.Second * 5,
	})
}
