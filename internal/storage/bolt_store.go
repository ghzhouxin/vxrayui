package storage

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"go.etcd.io/bbolt"

	"zhouxin.learn/go/vxrayui/config"
)

var (
	initOnce sync.Once
	vxrayDb  *bbolt.DB
)

const BucketNameVxray = "vxray"

func Init() {
	initOnce.Do(func() {
		initBoltStore()
	})
}

func initBoltStore() {
	config := config.GetStorage()
	db, err := bbolt.Open(config.Path, 0600, &bbolt.Options{
		Timeout: 1 * time.Second,
	})
	if err != nil {
		log.Fatalf("failed to open bolt db: %v", err)
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(BucketNameVxray))
		return err
	})
	if err != nil {
		log.Fatalf("failed to create bucket: %v", err)
	}

	vxrayDb = db
}

func Set[T any](key string, val T) error {
	return vxrayDb.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(BucketNameVxray))
		data, err := json.Marshal(val)
		if err != nil {
			return err
		}
		return b.Put([]byte(key), data)
	})
}

func Get[T any](key string) (T, error) {
	var val T
	err := vxrayDb.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(BucketNameVxray))
		data := b.Get([]byte(key))
		if data == nil {
			return nil
		}
		return json.Unmarshal(data, &val)
	})

	return val, err
}
