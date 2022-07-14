package main

import (
	"github.com/KirillMironov/ci/internal/repository"
	"github.com/KirillMironov/ci/internal/service"
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

	// Docker Client
	cli, err := client.NewClientWithOpts()
	if err != nil {
		logger.Fatal(err)
	}
	defer cli.Close()

	// App
	db, err := bolt.Open("ci.db", 0600, &bolt.Options{Timeout: time.Second})
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	vcsRepository, err := repository.NewVCS(db)
	if err != nil {
		logger.Fatal(err)
	}

	var (
		cloner   = &service.Cloner{}
		archiver = &service.Archiver{}
		parser   = &service.Parser{}
		executor = service.NewDockerExecutor(cli)
		poller   = service.NewPoller(".ci.yaml", cloner, archiver, parser, executor, vcsRepository, logger)
		handler  = transport.NewHandler(poller)
	)

	err = poller.Recover()
	if err != nil {
		logger.Fatal(err)
	}

	// HTTP Server
	err = handler.InitRoutes().Run(":8080")
	if err != nil {
		logger.Fatal(err)
	}
}
