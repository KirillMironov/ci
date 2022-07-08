package mock

import (
	"context"
	"github.com/KirillMironov/ci/internal/domain"
	"io"
)

type Executor struct{}

func (Executor) Execute(context.Context, domain.Step, io.Reader) error {
	return nil
}
