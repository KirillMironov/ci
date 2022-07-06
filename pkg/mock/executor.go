package mock

import (
	"context"
	"github.com/KirillMironov/ci/internal/domain"
)

type Executor struct{}

func (e Executor) Execute(context.Context, domain.Step) error {
	return nil
}
