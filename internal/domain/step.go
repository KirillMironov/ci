package domain

type Step struct {
	Name        string   `yaml:"name"`
	Image       string   `yaml:"image"`
	Environment []string `yaml:"env"`
	Command     []string `yaml:"command"`
}
