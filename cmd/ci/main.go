package main

import (
	"github.com/KirillMironov/ci/internal/service"
	"github.com/KirillMironov/ci/internal/transport"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

const repositoriesPath = "./.repositories"

func main() {
	// Logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "01|02 15:04:05.000",
	})

	// Docker Client
	cli, err := client.NewClientWithOpts()
	if err != nil {
		logger.Fatal(err)
	}
	defer cli.Close()

	// Repositories path
	err = os.Mkdir(repositoriesPath, 0750)
	if err != nil {
		logger.Fatal(err)
	}
	abs, err := filepath.Abs(repositoriesPath)
	if err != nil {
		logger.Fatal(err)
	}

	// App
	var (
		cloner   = service.NewCloner(abs)
		parser   = &service.Parser{}
		executor = service.NewDockerExecutor(cli)
		poller   = service.NewPoller(cloner, parser, executor, logger)
		handler  = transport.NewHandler(poller)
	)

	// HTTP Server
	err = handler.InitRoutes().Run(":8080")
	if err != nil {
		logger.Fatal(err)
	}
}
