package storage

import (
	"bytes"
	"encoding/gob"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/boltdb/bolt"
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
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(repo); err != nil {
		return err
	}

	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(r.bucket))
		return b.Put([]byte(repo.URL), buf.Bytes())
	})
}

// Delete deletes a repository by its URL.
func (r Repositories) Delete(repoURL domain.RepositoryURL) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(r.bucket))
		return b.Delete([]byte(repoURL))
	})
}

// GetAll returns all repositories.
func (r Repositories) GetAll() (repos []domain.Repository, err error) {
	err = r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(r.bucket))
		return b.ForEach(func(k, v []byte) error {
			var repo domain.Repository
			decoder := gob.NewDecoder(bytes.NewReader(v))
			if err = decoder.Decode(&repo); err != nil {
				return err
			}
			repos = append(repos, repo)
			return nil
		})
	})
	return repos, err
}
