package domain

import (
	"context"
	"encoding/json"
)

type Build struct {
	Id     string
	LogId  string
	Commit Commit
	Status Status
}

func (b Build) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Id     string `json:"id"`
		LogId  string `json:"log_id"`
		Commit Commit `json:"commit"`
		Status string `json:"status"`
	}{
		Id:     b.Id,
		LogId:  b.LogId,
		Commit: b.Commit,
		Status: b.Status.String(),
	})
}

type BuildsUsecase interface {
	Create(ctx context.Context, build Build, repoId string) error
	Delete(ctx context.Context, id, repoId string) error
	GetAllByRepoId(ctx context.Context, repoId string) ([]Build, error)
	GetById(ctx context.Context, id, repoId string) (Build, error)
}
