package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"sync"

	"go.etcd.io/bbolt"

	"example.com/nursery-cms/domain"
)

var (
	bucketRecords     = []byte("records")
	bucketEvents      = []byte("audit_events")
	bucketWorkflows   = []byte("workflows")
	bucketAttachments = []byte("attachments")
)

type Store struct {
	db   *bbolt.DB
	path string
	mu   sync.RWMutex
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("store path is required")
	}
	db, err := bbolt.Open(filepath.Clean(path), 0o600, nil)
	if err != nil {
		return nil, err
	}
	s := &Store{db: db, path: path}
	if err := s.initialize(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) initialize() error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		for _, bucket := range [][]byte{bucketRecords, bucketEvents, bucketWorkflows, bucketAttachments} {
			if _, err := tx.CreateBucketIfNotExists(bucket); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}

func (s *Store) Path() string { return s.path }

func encode(value any) ([]byte, error) { return json.Marshal(value) }

func decode(data []byte, target any) error {
	if len(data) == 0 {
		return errors.New("empty entity")
	}
	return json.Unmarshal(data, target)
}

func put(tx *bbolt.Tx, bucket []byte, key string, value any) error {
	b := tx.Bucket(bucket)
	if b == nil {
		return fmt.Errorf("bucket %s missing", bucket)
	}
	encoded, err := encode(value)
	if err != nil {
		return err
	}
	return b.Put([]byte(key), encoded)
}

func get(tx *bbolt.Tx, bucket []byte, key string, target any) error {
	b := tx.Bucket(bucket)
	if b == nil {
		return fmt.Errorf("bucket %s missing", bucket)
	}
	value := b.Get([]byte(key))
	if value == nil {
		return domain.ErrNotFound
	}
	return decode(value, target)
}

func list(tx *bbolt.Tx, bucket []byte, target func([]byte) error) error {
	b := tx.Bucket(bucket)
	if b == nil {
		return fmt.Errorf("bucket %s missing", bucket)
	}
	return b.ForEach(func(_, value []byte) error {
		if value == nil {
			return nil
		}
		return target(value)
	})
}

func (s *Store) update(fn func(*bbolt.Tx) error) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store is closed")
	}
	return s.db.Update(fn)
}

func (s *Store) view(fn func(*bbolt.Tx) error) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store is closed")
	}
	return s.db.View(fn)
}
