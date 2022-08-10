package service

import "github.com/KirillMironov/ci/internal/domain"

// Logs used to work with logs storage.
type Logs struct {
	storage logsStorage
}

type logsStorage interface {
	Create(domain.Log) (id int, err error)
	GetById(id int) (domain.Log, error)
}

func NewLogs(storage logsStorage) *Logs {
	return &Logs{storage: storage}
}

func (l Logs) Create(log domain.Log) (id int, err error) {
	return l.storage.Create(log)
}

func (l Logs) GetById(id int) (domain.Log, error) {
	return l.storage.GetById(id)
}
