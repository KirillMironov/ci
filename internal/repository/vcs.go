package repository

import (
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/boltdb/bolt"
	"time"
)

type VCS struct {
	db     *bolt.DB
	bucket string
}

func NewVCS(db *bolt.DB) (*VCS, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("vcs"))
		return err
	})
	return &VCS{
		db:     db,
		bucket: "vcs",
	}, err
}

func (v VCS) Put(vcs domain.VCS) error {
	return v.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(v.bucket))
		return b.Put([]byte(vcs.URL), []byte(vcs.PollingInterval.String()))
	})
}

func (v VCS) GetAll() (arr []domain.VCS, err error) {
	err = v.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(v.bucket))
		return b.ForEach(func(k, v []byte) error {
			interval, err := time.ParseDuration(string(v))
			if err != nil {
				return err
			}
			arr = append(arr, domain.VCS{
				URL:             string(k),
				PollingInterval: interval,
			})
			return nil
		})
	})
	return arr, err
}
