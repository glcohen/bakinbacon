package storage

import (
	"bakinbacon/util"
	bolt "go.etcd.io/bbolt"
)

func (s *BoltStorage) GetBakingWatermark() (int, error) {
	return s.getWatermark(BAKING_BUCKET)
}

func (s *BoltStorage) GetEndorsingWatermark() (int, error) {
	return s.getWatermark(ENDORSING_BUCKET)
}

func (s *BoltStorage) getWatermark(wBucket string) (int, error) {

	var watermark uint64

	err := s.View(func(tx *bolt.Tx) error {
		watermark = tx.Bucket([]byte(wBucket)).Sequence()
		return nil
	})

	return int(watermark), err
}

func (s *BoltStorage) RecordBakedBlock(level int, blockHash string) error {
	return s.recordOperation(BAKING_BUCKET, level, blockHash)
}

func (s *BoltStorage) RecordEndorsement(level int, endorsementHash string) error {
	return s.recordOperation(ENDORSING_BUCKET, level, endorsementHash)
}

func (s *BoltStorage) recordOperation(opBucket string, level int, opHash string) error {
	return s.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(opBucket))
		if err := b.SetSequence(uint64(level)); err != nil { // Record our watermark
			return err
		}
		return b.Put(util.IToB(level), []byte(opHash)) // Save the level:opHash
	})
}
