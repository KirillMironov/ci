package main

import (
	"github.com/KirillMironov/ci/config"
	"github.com/KirillMironov/ci/internal/service"
	"github.com/KirillMironov/ci/internal/storage"
	"github.com/KirillMironov/ci/internal/transport"
	"github.com/boltdb/bolt"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"time"
)

func main() {
	// Logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "01|02 15:04:05.000",
	})

	// Config
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal(err)
	}

	// Docker Client
	cli, err := client.NewClientWithOpts()
	if err != nil {
		logger.Fatal(err)
	}
	defer cli.Close()

	// BoltDB
	db, err := bolt.Open(cfg.BoltDBPath, 0600, &bolt.Options{Timeout: time.Second})
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	// App
	repositories, err := storage.NewRepositories(db, "repositories")
	if err != nil {
		logger.Fatal(err)
	}

	var (
		cloner   = &service.Cloner{}
		archiver = &service.Archiver{}
		parser   = &service.Parser{}
		executor = service.NewDockerExecutor(cli)
		poller   = service.NewPoller(cfg.CIFilename, cloner, archiver, parser, executor, repositories, logger)
		handler  = transport.NewHandler(poller)
	)

	err = poller.Recover()
	if err != nil {
		logger.Fatal(err)
	}

	// HTTP Server
	err = handler.InitRoutes().Run(":" + cfg.Port)
	if err != nil {
		logger.Fatal(err)
	}
}
