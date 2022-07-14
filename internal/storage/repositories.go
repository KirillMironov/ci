package storage

import (
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/boltdb/bolt"
	"time"
)

// Repositories is a boltdb-based repositories storage.
type Repositories struct {
	db     *bolt.DB
	bucket string
}

// NewRepositories creates a new bucket for repositories with a given name if it doesn't exist.
func NewRepositories(db *bolt.DB, bucket string) (*Repositories, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		return err
	})
	return &Repositories{
		db:     db,
		bucket: bucket,
	}, err
}

// Put adds or updates a repository.
func (r Repositories) Put(repo domain.Repository) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(r.bucket))
		return b.Put([]byte(repo.URL), []byte(repo.PollingInterval.String()))
	})
}

// GetAll returns all repositories.
func (r Repositories) GetAll() (repos []domain.Repository, err error) {
	err = r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(r.bucket))
		return b.ForEach(func(k, v []byte) error {
			interval, err := time.ParseDuration(string(v))
			if err != nil {
				return err
			}
			repos = append(repos, domain.Repository{
				URL:             string(k),
				PollingInterval: interval,
			})
			return nil
		})
	})
	return repos, err
}
