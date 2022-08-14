package storage

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/KirillMironov/ci/internal/domain"
	"go.etcd.io/bbolt"
)

// Builds used to store domain.Build in a BoltDB bucket.
type Builds struct {
	db     *bbolt.DB
	bucket string
}

// NewBuilds creates a new bucket for builds with a given name if it doesn't exist.
func NewBuilds(db *bbolt.DB, bucket string) (*Builds, error) {
	err := db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		return err
	})
	return &Builds{
		db:     db,
		bucket: bucket,
	}, err
}

func (b Builds) Create(_ context.Context, build domain.Build, repoId string) error {
	var buf bytes.Buffer
	var encoder = gob.NewEncoder(&buf)

	return b.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(b.bucket))
		nestedBucket, err := bucket.CreateBucketIfNotExists([]byte(repoId))
		if err != nil {
			return err
		}
		if err = encoder.Encode(build); err != nil {
			return err
		}
		return nestedBucket.Put([]byte(build.Id), buf.Bytes())
	})
}

func (b Builds) Delete(_ context.Context, id, repoId string) error {
	return b.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(b.bucket))
		nestedBucket := bucket.Bucket([]byte(repoId))
		if nestedBucket == nil {
			return domain.ErrNotFound
		}
		return nestedBucket.Delete([]byte(id))
	})
}

func (b Builds) GetAllByRepoId(_ context.Context, repoId string) (builds []domain.Build, err error) {
	err = b.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(b.bucket))
		nestedBucket := bucket.Bucket([]byte(repoId))
		if nestedBucket == nil {
			return domain.ErrNotFound
		}
		return nestedBucket.ForEach(func(k, v []byte) error {
			var build domain.Build
			var decoder = gob.NewDecoder(bytes.NewReader(v))
			if err = decoder.Decode(&build); err != nil {
				return err
			}
			builds = append(builds, build)
			return nil
		})
	})
	return builds, err
}

func (b Builds) GetById(_ context.Context, id, repoId string) (build domain.Build, err error) {
	err = b.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(b.bucket))
		nestedBucket := bucket.Bucket([]byte(repoId))
		if nestedBucket == nil {
			return domain.ErrNotFound
		}
		v := nestedBucket.Get([]byte(id))
		if v == nil {
			return domain.ErrNotFound
		}
		decoder := gob.NewDecoder(bytes.NewReader(v))
		return decoder.Decode(&build)
	})
	return build, err
}
