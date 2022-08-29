package main

import (
	"context"
	"errors"
	"github.com/KirillMironov/ci/config"
	"github.com/KirillMironov/ci/internal/service"
	"github.com/KirillMironov/ci/internal/storage"
	"github.com/KirillMironov/ci/internal/transport"
	"github.com/docker/docker/client"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
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

	// SQLite
	db, err := sqlx.Connect("sqlite3", cfg.SQLite.Path)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(cfg.SQLite.Schema)
	if err != nil {
		logger.Fatal(err)
	}

	// App
	var (
		repositoriesStorage = storage.NewRepositories(db)
		buildsStorage       = storage.NewBuilds(db)
		logsStorage         = storage.NewLogs(db)

		archiver  = &service.TarArchiver{}
		parser    = &service.YAMLParser{}
		cloner    = service.NewCloner(cfg.RepositoriesDir)
		executor  = service.NewDockerExecutor(cli, cfg.ContainerWorkingDir, archiver)
		runner    = service.NewRunner(executor, buildsStorage, logger)
		poller    = service.NewPoller(cfg.CIFilename, cloner, parser, runner, buildsStorage, logger)
		scheduler = service.NewScheduler(poller, repositoriesStorage, logger)

		handler = transport.NewHandler(cfg.StaticRootDir, scheduler, repositoriesStorage, buildsStorage, logsStorage)
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
