package usecase

import (
	"fmt"
	"github.com/KirillMironov/ci/internal/domain"
)

type Logs struct {
	storage logsStorage
}

type logsStorage interface {
	GetByBuildId(id string) (domain.Log, error)
}

func NewLogs(storage logsStorage) *Logs {
	return &Logs{storage: storage}
}

func (l Logs) GetByBuildId(buildId string) (domain.Log, error) {
	log, err := l.storage.GetByBuildId(buildId)
	if err != nil {
		return domain.Log{}, fmt.Errorf("failed to get log by build id: %w", err)
	}

	return log, nil
}
