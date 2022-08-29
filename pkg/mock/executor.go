package mock

import (
	"context"
	"github.com/KirillMironov/ci/internal/domain"
	"io"
	"strings"
)

type Executor struct {
	HasError bool
	Log      string
}

func (e Executor) ExecuteStep(context.Context, domain.Step, string) (io.ReadCloser, error) {
	var logs = io.NopCloser(strings.NewReader(e.Log))

	if e.HasError {
		return logs, domain.ExitError{Code: 1}
	}
	return logs, nil
}
