package main

import (
	"context"
	"github.com/KirillMironov/ci/config"
	"github.com/KirillMironov/ci/internal/domain"
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
		archiver  = &service.TarArchiver{}
		parser    = &service.YAMLParser{}
		cloner    = service.NewCloner(cfg.RepositoriesDir, archiver)
		executor  = service.NewDockerExecutor(cli, cfg.ContainerWorkingDir)
		poller    = service.NewPoller(cfg.CIFilename, cloner, executor, archiver, parser)
		add       = make(chan domain.Repository)
		scheduler = service.NewScheduler(add, poller, repositories, logger)
		handler   = transport.NewHandler(add)
	)

	ctx, cancel := context.WithCancel(context.Background())
	if err != nil {
		logger.Fatal(err)
	}
	defer cancel()

	// Scheduler
	go scheduler.Start(ctx)

	// HTTP Server
	err = handler.InitRoutes().Run(":" + cfg.Port)
	if err != nil {
		logger.Fatal(err)
	}
}
