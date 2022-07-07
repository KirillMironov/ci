package main

import (
	"github.com/KirillMironov/ci/internal/service"
	"github.com/KirillMironov/ci/internal/transport"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

func main() {
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
		cloner   = &service.Cloner{}
		parser   = &service.Parser{}
		executor = service.NewDockerExecutor(cli)
		poller   = service.NewPoller(cloner, parser, executor, logger)
		handler  = transport.NewHandler(poller)
	)

	err = handler.InitRoutes().Run(":8080")
	if err != nil {
		logger.Fatal(err)
	}
}
