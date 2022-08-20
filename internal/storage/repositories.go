package storage

import (
	"database/sql"
	"errors"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/jmoiron/sqlx"
	"time"
)

type Repositories struct {
	db *sqlx.DB
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{db: db}
}

func (r Repositories) Create(repo domain.Repository) error {
	var query = "INSERT INTO repositories (id, url, branch, polling_interval, created_at) VALUES ($1, $2, $3, $4, $5)"

	_, err := r.db.Exec(query, repo.Id, repo.URL, repo.Branch, repo.PollingInterval, time.Now())
	return err
}

func (r Repositories) Delete(id string) error {
	var query = "DELETE FROM repositories WHERE id = $1"

	_, err := r.db.Exec(query, id)
	return err
}

func (r Repositories) GetAll() (repos []domain.Repository, err error) {
	var query = "SELECT id, url, branch, polling_interval, created_at FROM repositories"

	rows, err := r.db.Queryx(query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var repo domain.Repository
		err = rows.Scan(&repo.Id, &repo.URL, &repo.Branch, &repo.PollingInterval, &repo.CreatedAt)
		if err != nil {
			return nil, err
		}
		repos = append(repos, repo)
	}

	return repos, rows.Err()
}

func (r Repositories) GetById(id string) (repo domain.Repository, err error) {
	var query = "SELECT id, url, branch, polling_interval, created_at FROM repositories WHERE id = $1"

	row := r.db.QueryRowx(query, id)

	err = row.Scan(&repo.Id, &repo.URL, &repo.Branch, &repo.PollingInterval, &repo.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Repository{}, domain.ErrNotFound
		}
		return domain.Repository{}, err
	}

	return repo, nil
}

func (r Repositories) GetByURL(url string) (repo domain.Repository, err error) {
	var query = "SELECT id, url, branch, polling_interval, created_at FROM repositories WHERE url = $1"

	row := r.db.QueryRowx(query, url)

	err = row.Scan(&repo.Id, &repo.URL, &repo.Branch, &repo.PollingInterval, &repo.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Repository{}, domain.ErrNotFound
		}
		return domain.Repository{}, err
	}

	return repo, nil
}
