package main

import (
	"context"
	"errors"
	"github.com/KirillMironov/ci/config"
	"github.com/KirillMironov/ci/internal/service"
	"github.com/KirillMironov/ci/internal/storage"
	"github.com/KirillMironov/ci/internal/transport"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"go.etcd.io/bbolt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	db, err := bbolt.Open(cfg.BoltDBPath, 0600, &bbolt.Options{Timeout: time.Second})
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	// App
	repositories, err := storage.NewRepositories(db, "repositories")
	if err != nil {
		logger.Fatal(err)
	}

	logs, err := storage.NewLogs(db, "logs")
	if err != nil {
		logger.Fatal(err)
	}

	var (
		archiver  = &service.TarArchiver{}
		parser    = &service.YAMLParser{}
		cloner    = service.NewCloner(cfg.RepositoriesDir, archiver)
		executor  = service.NewDockerExecutor(cli, cfg.ContainerWorkingDir)
		runner    = service.NewRunner(executor)
		poller    = service.NewPoller(cfg.CIFilename, runner, cloner, archiver, parser, repositories, logs, logger)
		scheduler = service.NewScheduler(poller, repositories, logger)
		handler   = transport.NewHandler(scheduler)
	)

	// Scheduler & Poller
	ctx, cancel := context.WithCancel(context.Background())
	if err != nil {
		logger.Fatal(err)
	}
	defer cancel()

	go scheduler.Start(ctx)
	go poller.Start(ctx)

	// HTTP Server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: handler.InitRoutes(),
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	srvCtx, srvCancel := context.WithTimeout(context.Background(), time.Second*5)
	defer srvCancel()

	err = srv.Shutdown(srvCtx)
	if err != nil {
		logger.Fatal(err)
	}
}
