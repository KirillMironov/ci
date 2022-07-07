package mock

import "github.com/KirillMironov/ci/internal/domain"

type Parser struct{}

func (Parser) ParsePipeline(string) (domain.Pipeline, error) {
	return domain.Pipeline{}, nil
}
