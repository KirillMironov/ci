package main

import (
	"github.com/KirillMironov/ci/internal/domain"
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

	log.Println(pipeline)
}
