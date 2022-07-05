package domain

type Pipeline struct {
	Name  string `yaml:"name"`
	Steps []Step `yaml:"steps"`
}
