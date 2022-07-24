package storage

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"github.com/KirillMironov/ci/internal/domain"
	"go.etcd.io/bbolt"
)

// Logs used to store domain.Log in boltdb.
type Logs struct {
	db     *bbolt.DB
	bucket string
}

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

func (l Logs) Put(log domain.Log) (id int, err error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	return log.Id, l.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(l.bucket))
		id, _ := b.NextSequence()
		log.Id = int(id)
		if err = encoder.Encode(log); err != nil {
			return err
		}
		return b.Put(intToByte(log.Id), buf.Bytes())
	})
}

func (l Logs) GetById(id int) (log domain.Log, err error) {
	err = l.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(l.bucket))
		v := b.Get(intToByte(id))
		decoder := gob.NewDecoder(bytes.NewReader(v))
		return decoder.Decode(&log)
	})
	return log, err
}

func intToByte(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
