package service

import (
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const repositoriesPath = "./.repositories"

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

	err = os.Mkdir(repositoriesPath, 0750)
	if err != nil {
		logger.Fatal(err)
	}
	abs, err := filepath.Abs(repositoriesPath)
	if err != nil {
		logger.Fatal(err)
	}

	var (
		cloner   = NewCloner(abs)
		parser   = &Parser{}
		executor = NewDockerExecutor(cli)
		poller   = NewPoller(cloner, parser, executor, logger)
	)

	poller.Start(domain.VCS{
		URL:             "",
		PollingInterval: time.Second * 5,
	})
}
