package mock

import (
	"context"
	"github.com/KirillMironov/ci/internal/domain"
)

type Executor struct{}

func (Executor) Execute(context.Context, domain.Step, string) error {
	return nil
}
