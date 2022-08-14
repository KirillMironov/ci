package storage

import (
	"bytes"
	"context"
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

func (r Repositories) Create(_ context.Context, repo domain.Repository) error {
	var buf bytes.Buffer
	var encoder = gob.NewEncoder(&buf)

	return r.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(r.bucket))
		err := encoder.Encode(repo)
		if err != nil {
			return err
		}
		return bucket.Put([]byte(repo.Id), buf.Bytes())
	})
}

func (r Repositories) Delete(_ context.Context, id string) error {
	return r.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(r.bucket))
		return bucket.Delete([]byte(id))
	})
}

func (r Repositories) GetAll(_ context.Context) (repos []domain.Repository, err error) {
	err = r.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(r.bucket))
		return bucket.ForEach(func(k, v []byte) error {
			var repo domain.Repository
			var decoder = gob.NewDecoder(bytes.NewReader(v))
			if err = decoder.Decode(&repo); err != nil {
				return err
			}
			repos = append(repos, repo)
			return nil
		})
	})
	return repos, err
}

func (r Repositories) GetById(_ context.Context, id string) (repo domain.Repository, err error) {
	err = r.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(r.bucket))
		v := bucket.Get([]byte(id))
		if v == nil {
			return domain.ErrNotFound
		}
		decoder := gob.NewDecoder(bytes.NewReader(v))
		return decoder.Decode(&repo)
	})
	return repo, err
}

func (r Repositories) GetByURL(_ context.Context, url string) (repo domain.Repository, err error) {
	err = r.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(r.bucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var tempRepo domain.Repository
			decoder := gob.NewDecoder(bytes.NewReader(v))
			if err = decoder.Decode(&tempRepo); err != nil {
				return err
			}
			if tempRepo.URL == url {
				repo = tempRepo
				return nil
			}
		}
		return domain.ErrNotFound
	})
	return repo, err
}
