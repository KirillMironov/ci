package mock

import (
	"context"
	"github.com/KirillMironov/ci/internal/domain"
)

type Poller struct{}

func (Poller) AddRepository(context.Context, domain.Repository) {}
