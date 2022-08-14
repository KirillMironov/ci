package domain

import "context"

type Log struct {
	Id   string
	Data []byte
}

type LogsUsecase interface {
	Create(context.Context, Log) (id string, err error)
	GetById(ctx context.Context, id string) (Log, error)
}
