package main

import (
	"context"
	"errors"
	"github.com/KirillMironov/ci/config"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/KirillMironov/ci/internal/service"
	"github.com/KirillMironov/ci/internal/storage"
	"github.com/KirillMironov/ci/internal/transport"
	"github.com/KirillMironov/ci/internal/usecase"
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
	repositoriesStorage, err := storage.NewRepositories(db, "repositories")
	if err != nil {
		logger.Fatal(err)
	}

	buildsStorage, err := storage.NewBuilds(db, "builds")
	if err != nil {
		logger.Fatal(err)
	}

	logsStorage, err := storage.NewLogs(db, "logs")
	if err != nil {
		logger.Fatal(err)
	}

	var (
		add    = make(chan domain.Repository)
		remove = make(chan string)
		run    = make(chan service.RunRequest)

		repositoriesUsecase = usecase.NewRepositories(repositoriesStorage, add, remove)
		buildsUsecase       = usecase.NewBuilds(buildsStorage)
		logsUsecase         = usecase.NewLogs(logsStorage)

		archiver  = &service.TarArchiver{}
		parser    = &service.YAMLParser{}
		cloner    = service.NewCloner(cfg.RepositoriesDir)
		executor  = service.NewDockerExecutor(cli, cfg.ContainerWorkingDir, archiver)
		runner    = service.NewRunner(run, executor, buildsUsecase, logsUsecase, logger)
		poller    = service.NewPoller(run, cfg.CIFilename, cloner, parser, buildsUsecase, logger)
		scheduler = service.NewScheduler(add, remove, poller, repositoriesUsecase, logger)

		handler = transport.NewHandler(cfg.StaticRootPath, repositoriesUsecase, buildsUsecase, logsUsecase)
	)

	// Scheduler & Poller & Runner
	ctx, cancel := context.WithCancel(context.Background())
	if err != nil {
		logger.Fatal(err)
	}
	defer cancel()

	go scheduler.Start(ctx)
	go poller.Start(ctx)
	go runner.Start(ctx)

	// HTTP Server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: handler.Routes(),
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
