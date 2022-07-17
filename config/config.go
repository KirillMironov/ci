package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port                string `default:"8080" envconfig:"PORT"`
	CIFilename          string `default:".ci.yaml" envconfig:"CI_FILENAME"`
	BoltDBPath          string `default:"./ci.db" envconfig:"BOLT_DB_PATH"`
	RepositoriesDir     string `default:"./.cache/git/" envconfig:"REPOSITORIES_DIR"`
	ContainerWorkingDir string `default:"/ci" envconfig:"CONTAINER_WORKING_DIR"`
}

func Load() (cfg Config, _ error) {
	return cfg, envconfig.Process("", &cfg)
}
