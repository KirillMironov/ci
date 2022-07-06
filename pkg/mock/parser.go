package mock

import "github.com/KirillMironov/ci/internal/domain"

type Parser struct{}

func (p Parser) ParsePipeline(string) (domain.Pipeline, error) {
	return domain.Pipeline{}, nil
}
