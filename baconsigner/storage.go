package baconsigner

import (
	"bakinbacon/storage"
	"bakinbacon/util"
	bolt "go.etcd.io/bbolt"
)

const (
	SIGNER_SK       = "signersk"
	PUBLIC_KEY_HASH = "pkh"
	CONFIG_BUCKET   = "config"
	SIGNER_TYPE     = "signertype"
)

// SignerStorage - Generic signer storage interface
type SignerStorage interface {
	GetDelegate() (string, string, error)
	SetDelegate(sk, pkh string) error
	GetSignerType() (int, error)
	SetSignerType(signerType int) error
	GetSignerSk() (string, error)
}

// Storage - Specific implementation of signer storage for bolt db
// other implementations can be considered for testing, etc.
type Storage struct {
	db *storage.BoltStorage
}

func (s *Storage) GetDelegate() (string, string, error) {

	var sk, pkh string

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(CONFIG_BUCKET))
		sk = string(b.Get([]byte(SIGNER_SK)))
		pkh = string(b.Get([]byte(PUBLIC_KEY_HASH)))
		return nil
	})

	return sk, pkh, err
}

func (s *Storage) SetDelegate(sk, pkh string) error {

	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(CONFIG_BUCKET))
		if err := b.Put([]byte(SIGNER_SK), []byte(sk)); err != nil {
			return err
		}
		if err := b.Put([]byte(PUBLIC_KEY_HASH), []byte(pkh)); err != nil {
			return err
		}
		return nil
	})
}

func (s *Storage) GetSignerType() (int, error) {

	var signerType int = 0

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(CONFIG_BUCKET))
		signerTypeBytes := b.Get([]byte(SIGNER_TYPE))
		if signerTypeBytes != nil {
			signerType = util.BToI(signerTypeBytes)
		}
		return nil
	})

	return signerType, err
}

func (s *Storage) SetSignerType(signerType int) error {

	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(CONFIG_BUCKET))
		return b.Put([]byte(SIGNER_TYPE), util.IToB(signerType))
	})
}

func (s *Storage) GetSignerSk() (string, error) {

	var sk string

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(CONFIG_BUCKET))
		sk = string(b.Get([]byte(SIGNER_SK)))
		return nil
	})

	return sk, err
}

func (s *Storage) SetSignerSk(sk string) error {

	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(CONFIG_BUCKET))
		return b.Put([]byte(SIGNER_SK), []byte(sk))
	})
}
