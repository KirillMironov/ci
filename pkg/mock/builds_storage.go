package mock

import "github.com/KirillMironov/ci/internal/domain"

type builds struct {
	storage map[string]domain.Build
}

func NewBuilds() *builds {
	return &builds{
		storage: make(map[string]domain.Build),
	}
}

func (b builds) Create(build domain.Build) error {
	b.storage[build.Id] = build
	return nil
}

func (b builds) Update(build domain.Build) error {
	b.storage[build.Id] = build
	return nil
}

func (b builds) Delete(id string) error {
	delete(b.storage, id)
	return nil
}

func (b builds) GetAllByRepoId(repoId string) (builds []domain.Build, _ error) {
	for _, build := range b.storage {
		if build.RepoId == repoId {
			builds = append(builds, build)
		}
	}
	if len(builds) == 0 {
		return nil, domain.ErrNotFound
	}
	return builds, nil
}

func (b builds) GetById(id string) (domain.Build, error) {
	build, ok := b.storage[id]
	if !ok {
		return domain.Build{}, domain.ErrNotFound
	}
	return build, nil
}
