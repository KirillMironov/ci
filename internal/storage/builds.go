package storage

import (
	"database/sql"
	"errors"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/jmoiron/sqlx"
	"time"
)

type Builds struct {
	db *sqlx.DB
}

func NewBuilds(db *sqlx.DB) *Builds {
	return &Builds{db: db}
}

func (b Builds) Create(build domain.Build) error {
	var (
		buildQuery  = "INSERT INTO builds (id, repo_id, status, created_at) VALUES ($1, $2, $3, $4)"
		commitQuery = "INSERT INTO commits (build_id, hash) VALUES ($1, $2)"
		logQuery    = "INSERT INTO logs (build_id, data) VALUES ($1, $2)"
	)

	tx, err := b.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(buildQuery, build.Id, build.RepoId, build.Status, time.Now())
	if err != nil {
		return err
	}

	_, err = tx.Exec(commitQuery, build.Id, build.Commit.Hash)
	if err != nil {
		return err
	}

	_, err = tx.Exec(logQuery, build.Id, build.Log.Data)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (b Builds) Delete(id string) error {
	var query = "DELETE FROM builds WHERE id = $1"

	_, err := b.db.Exec(query, id)
	return err
}

func (b Builds) GetAllByRepoId(repoId string) (builds []domain.Build, err error) {
	var query = `SELECT id, repo_id, status, created_at, c.hash FROM builds b 
    	JOIN commits c ON b.id = c.build_id WHERE b.repo_id = $1`

	rows, err := b.db.Queryx(query, repoId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var build domain.Build
		err = rows.Scan(&build.Id, &build.RepoId, &build.Status, &build.CreatedAt, &build.Commit.Hash)
		if err != nil {
			return nil, err
		}
		builds = append(builds, build)
	}

	return builds, rows.Err()
}

func (b Builds) GetById(id string) (build domain.Build, err error) {
	var query = `SELECT id, repo_id, status, created_at, c.hash FROM builds b 
    	JOIN commits c ON b.id = c.build_id WHERE b.id = $1`

	row := b.db.QueryRowx(query, id)

	err = row.Scan(&build.Id, &build.RepoId, &build.Status, &build.CreatedAt, &build.Commit.Hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Build{}, domain.ErrNotFound
		}
		return domain.Build{}, err
	}

	return build, nil
}
