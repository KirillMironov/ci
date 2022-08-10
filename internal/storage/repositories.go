package storage

import (
	"bytes"
	"encoding/gob"
	"github.com/KirillMironov/ci/internal/domain"
	"go.etcd.io/bbolt"
)

// Repositories used to store domain.Repository in a BoltDB bucket.
type Repositories struct {
	db     *bbolt.DB
	bucket string
}

// NewRepositories creates a new bucket for repositories with a given name if it doesn't exist.
func NewRepositories(db *bbolt.DB, bucket string) (*Repositories, error) {
	err := db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		return err
	})
	return &Repositories{
		db:     db,
		bucket: bucket,
	}, err
}

func (r Repositories) Create(repo domain.Repository) error {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(repo); err != nil {
		return err
	}

	return r.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(r.bucket))
		return b.Put([]byte(repo.Id), buf.Bytes())
	})
}

func (r Repositories) Delete(id string) error {
	return r.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(r.bucket))
		return b.Delete([]byte(id))
	})
}

func (r Repositories) GetAll() (repos []domain.Repository, err error) {
	err = r.db.View(func(tx *bbolt.Tx) error {
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

func (r Repositories) GetById(id string) (repo domain.Repository, err error) {
	err = r.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(r.bucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var tempRepo domain.Repository
			decoder := gob.NewDecoder(bytes.NewReader(v))
			if err = decoder.Decode(&tempRepo); err != nil {
				return err
			}
			if tempRepo.Id == id {
				repo = tempRepo
				return nil
			}
		}
		return domain.ErrNotFound
	})
	return repo, err
}
