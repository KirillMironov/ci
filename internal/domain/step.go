package domain

type Step struct {
	Name        string              `yaml:"name"`
	Image       string              `yaml:"image"`
	Environment []map[string]string `yaml:"env"`
	Command     []string            `yaml:"command"`
}
