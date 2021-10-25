package storage

import (
	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
)

func (s *BoltStorage) GetNotifiersConfig(notifier string) ([]byte, error) {

	var config []byte

	err := s.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(CONFIG_BUCKET)).Bucket([]byte(NOTIFICATIONS_BUCKET))
		if b == nil {
			return errors.New("Unable to locate notifications bucket")
		}

		config = b.Get([]byte(notifier))

		return nil
	})

	return config, err
}

func (s *BoltStorage) SaveNotifiersConfig(notifier string, config []byte) error {

	return s.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(CONFIG_BUCKET)).Bucket([]byte(NOTIFICATIONS_BUCKET))
		if b == nil {
			return errors.New("Unable to locate notifications bucket")
		}

		return b.Put([]byte(notifier), config)
	})
}
