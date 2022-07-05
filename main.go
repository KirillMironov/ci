package main

import (
	"context"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func main() {
	var pipeline domain.Pipeline

	data, err := ioutil.ReadFile("ci.yaml")
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(data, &pipeline)
	if err != nil {
		log.Fatal(err)
	}

	cli, err := client.NewClientWithOpts()
	if err != nil {
		log.Fatal(err)
	}

	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(images)
}
