package storage

import (
	"database/sql"
	"errors"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/jmoiron/sqlx"
)

type Logs struct {
	db *sqlx.DB
}

func NewLogs(db *sqlx.DB) *Logs {
	return &Logs{db: db}
}

func (l Logs) GetByBuildId(buildId string) (log domain.Log, err error) {
	var query = "SELECT data FROM logs WHERE build_id = $1"

	row := l.db.QueryRow(query, buildId)

	err = row.Scan(&log.Data)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Log{}, domain.ErrNotFound
		}
		return domain.Log{}, err
	}

	return log, nil
}
