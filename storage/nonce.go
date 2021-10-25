package storage

import (
	"bakinbacon/util"
	"encoding/json"

	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"

	"bakinbacon/nonce"
)

func (s *BoltStorage) SaveNonce(cycle int, n nonce.Nonce) error {

	// Nonces are stored within a cycle bucket for easy retrieval
	nonceBytes, err := json.Marshal(n)
	if err != nil {
		return errors.Wrap(err, "Unable to marshal nonce")
	}

	return s.Update(func(tx *bolt.Tx) error {
		cb, err := tx.Bucket([]byte(NONCE_BUCKET)).CreateBucketIfNotExists(util.IToB(cycle))
		if err != nil {
			return errors.Wrap(err, "Unable to create nonce-cycle bucket")
		}
		return cb.Put(util.IToB(n.Level), nonceBytes)
	})
}

func (s *BoltStorage) GetNoncesForCycle(cycle int) ([]nonce.Nonce, error) {

	// Get back all nonces for cycle
	var nonces []nonce.Nonce

	err := s.Update(func(tx *bolt.Tx) error {
		cb, err := tx.Bucket([]byte(NONCE_BUCKET)).CreateBucketIfNotExists(util.IToB(cycle))
		if err != nil {
			return errors.Wrap(err, "Unable to create nonce-cycle bucket")
		}
		c := cb.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			var n nonce.Nonce
			if err := json.Unmarshal(v, &n); err != nil {
				log.WithError(err).Error("Unable to unmarshal nonce")
				continue
			}

			nonces = append(nonces, n)
		}

		return nil
	})

	return nonces, err
}
