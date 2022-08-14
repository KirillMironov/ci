package usecase

import (
	"context"
	"fmt"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/rs/xid"
)

type Logs struct {
	storage logsStorage
}

type logsStorage interface {
	Create(context.Context, domain.Log) error
	GetById(ctx context.Context, id string) (domain.Log, error)
}

func NewLogs(storage logsStorage) *Logs {
	return &Logs{storage: storage}
}

func (l Logs) Create(ctx context.Context, log domain.Log) (id string, err error) {
	log.Id = xid.New().String()

	err = l.storage.Create(ctx, log)
	if err != nil {
		return "", fmt.Errorf("failed to create log: %w", err)
	}

	return log.Id, nil
}

func (l Logs) GetById(ctx context.Context, id string) (domain.Log, error) {
	log, err := l.storage.GetById(ctx, id)
	if err != nil {
		return domain.Log{}, fmt.Errorf("failed to get log: %w", err)
	}

	return log, nil
}
