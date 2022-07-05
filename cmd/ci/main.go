package main

import (
	"context"
	"github.com/KirillMironov/ci/internal/service"
	"github.com/docker/docker/client"
	"io/ioutil"
	"log"
)

func main() {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		log.Fatal(err)
	}

	var (
		parser   = service.Parser{}
		executor = service.NewExecutor(cli)
	)

	data, err := ioutil.ReadFile("ci.yaml")
	if err != nil {
		log.Fatal(err)
	}

	pipeline, err := parser.ParsePipeline(string(data))
	if err != nil {
		log.Fatal(err)
	}

	for _, step := range pipeline.Steps {
		err = executor.Execute(context.Background(), step)
		if err != nil {
			log.Fatal(err)
		}
	}
}
