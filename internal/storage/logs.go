package storage

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/KirillMironov/ci/internal/domain"
	"go.etcd.io/bbolt"
)

// Logs used to store domain.Log in a boltdb bucket.
type Logs struct {
	db     *bbolt.DB
	bucket string
}

// NewLogs creates a new bucket for logs with a given name if it doesn't exist.
func NewLogs(db *bbolt.DB, bucket string) (*Logs, error) {
	err := db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		return err
	})
	return &Logs{
		db:     db,
		bucket: bucket,
	}, err
}

func (l Logs) Create(_ context.Context, log domain.Log) error {
	var buf bytes.Buffer
	var encoder = gob.NewEncoder(&buf)

	return l.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(l.bucket))
		if err := encoder.Encode(log); err != nil {
			return err
		}
		return bucket.Put([]byte(log.Id), buf.Bytes())
	})
}

func (l Logs) GetById(_ context.Context, id string) (log domain.Log, err error) {
	err = l.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(l.bucket))
		v := bucket.Get([]byte(id))
		if v == nil {
			return domain.ErrNotFound
		}
		decoder := gob.NewDecoder(bytes.NewReader(v))
		return decoder.Decode(&log)
	})
	return log, err
}
