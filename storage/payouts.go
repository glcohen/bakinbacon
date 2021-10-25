package storage

import (
	"bakinbacon/util"
	"encoding/json"

	"github.com/pkg/errors"

	bolt "go.etcd.io/bbolt"

	"bakinbacon/payouts"
)

const (
	METADATA = "metadata"
)

// GetPayoutsMetadata returns a byte-slice of raw JSON from DB
func (s *BoltStorage) GetPayoutsMetadata() (map[int]json.RawMessage, error) {

	payoutsMetadata := make(map[int]json.RawMessage)

	err := s.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(PAYOUTS_BUCKET))
		if b == nil {
			return errors.New("Unable to locate cycle payouts bucket")
		}

		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {

			// keys are cycle numbers, which are buckets of data
			cycleBucket := b.Bucket(k)
			cycle := util.BToI(k)

			// Get metadata key from bucket
			metadataBytes := cycleBucket.Get([]byte(METADATA))

			// Unmarshal ...
			var tmpMetadata json.RawMessage
			if err := json.Unmarshal(metadataBytes, &tmpMetadata); err != nil {
				return errors.Wrap(err, "Unable to fetch metadata")
			}

			// ... and add to map
			payoutsMetadata[cycle] = tmpMetadata
		}

		return nil
	})

	return payoutsMetadata, err
}

func (s *BoltStorage) GetCyclePayouts(cycle int) (map[string]json.RawMessage, error) {

	// key is delegator address
	cyclePayouts := make(map[string]json.RawMessage)

	err := s.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(PAYOUTS_BUCKET)).Bucket(util.IToB(cycle))
		if b == nil {
			return errors.New("Unable to locate cycle payouts bucket")
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if string(k) == METADATA {
				continue
			}

			// Add the value (JSON object) to the map
			cyclePayouts[string(k)] = v
		}

		return nil
	})

	return cyclePayouts, err
}

func (s *BoltStorage) SaveCycleRewardMetadata(rewardCycle int, metadata *payouts.CycleRewardMetadata) error {

	// Marshal metadata to JSON and store in 'metadata' key
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return errors.Wrap(err, "Unable to marshal rewards metadata")
	}

	return s.Update(func(tx *bolt.Tx) error {
		b, err := tx.Bucket([]byte(PAYOUTS_BUCKET)).CreateBucketIfNotExists(util.IToB(rewardCycle))
		if err != nil {
			return errors.New("Unable to create cycle payouts bucket")
		}

		return b.Put([]byte(METADATA), metadataBytes)
	})
}

func (s *BoltStorage) SaveDelegatorReward(rewardCycle int, rewardRecord *payouts.DelegatorReward) error {

	// Marshal reward
	rewardRecordBytes, err := json.Marshal(rewardRecord)
	if err != nil {
		return errors.Wrap(err, "Unable to marshal reward record")
	}

	return s.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(PAYOUTS_BUCKET)).Bucket(util.IToB(rewardCycle))
		if b == nil {
			return errors.New("Unable to locate cycle payouts bucket")
		}

		// Store the record as the value of the record address (key)
		// This will allow for easier scanning/searching for a payment record
		return b.Put([]byte(rewardRecord.Delegator), rewardRecordBytes)
	})
}

