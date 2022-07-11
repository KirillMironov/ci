package domain

// Pipeline is a collection of steps.
type Pipeline struct {
	Name  string `yaml:"name"`
	Steps []Step `yaml:"steps"`
}
