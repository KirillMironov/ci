package config

import (
	_ "embed"
	"github.com/kelseyhightower/envconfig"
)

//go:embed schema.sql
var schema string

type Config struct {
	Port string `default:"8080" envconfig:"PORT"`

	CIFilename          string `default:".ci.yaml" envconfig:"CI_FILENAME"`
	StaticRootDir       string `default:"./web/" envconfig:"STATIC_ROOT_DIR"`
	RepositoriesDir     string `default:"./.cache/git/" envconfig:"REPOSITORIES_DIR"`
	ContainerWorkingDir string `default:"/ci" envconfig:"CONTAINER_WORKING_DIR"`

	SQLite struct {
		Path   string `default:"./sqlite.db" envconfig:"SQLITE_PATH"`
		Schema string `envconfig:"SQLITE_SCHEMA"`
	}
}

func Load() (cfg Config, _ error) {
	cfg.SQLite.Schema = schema
	return cfg, envconfig.Process("", &cfg)
}
