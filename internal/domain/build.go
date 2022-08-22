package domain

import (
	"encoding/json"
	"time"
)

type Build struct {
	Id        string
	RepoId    string
	Commit    Commit
	Log       Log
	Status    Status
	CreatedAt time.Time
}

func (b Build) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Id        string    `json:"id"`
		Commit    Commit    `json:"commit"`
		Status    string    `json:"status"`
		CreatedAt time.Time `json:"created_at"`
	}{
		Id:        b.Id,
		Commit:    b.Commit,
		Status:    b.Status.String(),
		CreatedAt: b.CreatedAt,
	})
}

type BuildsUsecase interface {
	Create(Build) (id string, err error)
	Update(Build) error
	Delete(id string) error
	GetAllByRepoId(repoId string) ([]Build, error)
	GetById(id string) (Build, error)
}
