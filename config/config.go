package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port       string `default:"8080" envconfig:"PORT"`
	BoltDBPath string `default:"ci.db" envconfig:"BOLT_DB_PATH"`
	CIFilename string `default:".ci.yaml" envconfig:"CI_FILENAME"`
}

func Load() (cfg Config, _ error) {
	return cfg, envconfig.Process("", &cfg)
}
