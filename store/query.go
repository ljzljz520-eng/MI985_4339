package store

import (
	"strings"

	"example.com/nursery-cms/domain"
	"go.etcd.io/bbolt"
)

type RecordFilter struct {
	BatchID string
	Owner   string
	Text    string
	Status  []domain.RecordStatus
}

func (s *Store) FilterRecords(filter RecordFilter) ([]domain.Record, error) {
	result := make([]domain.Record, 0)
	err := s.view(func(tx *bbolt.Tx) error {
		return list(tx, bucketRecords, func(data []byte) error {
			var record domain.Record
			if err := decode(data, &record); err != nil {
				return err
			}
			if !matchRecord(record, filter) {
				return nil
			}
			result = append(result, record)
			return nil
		})
	})
	return result, err
}

func matchRecord(record domain.Record, filter RecordFilter) bool {
	if filter.BatchID != "" && record.BatchID != filter.BatchID {
		return false
	}
	if filter.Owner != "" && record.Owner != filter.Owner {
		return false
	}
	if len(filter.Status) > 0 && !statusIn(record.Status, filter.Status) {
		return false
	}
	if filter.Text != "" && !strings.Contains(record.Title, filter.Text) && !strings.Contains(record.Content, filter.Text) {
		return false
	}
	return true
}

func statusIn(value domain.RecordStatus, statuses []domain.RecordStatus) bool {
	for _, status := range statuses {
		if status == value {
			return true
		}
	}
	return false
}

func (s *Store) CountRecords(filter RecordFilter) (int, error) {
	items, err := s.FilterRecords(filter)
	return len(items), err
}

func (s *Store) SaveRecords(records []domain.Record) error {
	return s.update(func(tx *bbolt.Tx) error {
		for _, record := range records {
			if err := record.Validate(); err != nil {
				return err
			}
			if err := put(tx, bucketRecords, record.ID, record); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) DeleteBatch(batchID string) error {
	return s.update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(bucketRecords)
		keys := make([][]byte, 0)
		if err := bucket.ForEach(func(key, value []byte) error {
			if value == nil {
				return nil
			}
			var record domain.Record
			if err := decode(value, &record); err != nil {
				return err
			}
			if record.BatchID == batchID {
				keys = append(keys, append([]byte(nil), key...))
			}
			return nil
		}); err != nil {
			return err
		}
		for _, key := range keys {
			if err := bucket.Delete(key); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) SaveAllEntities(record domain.Record, event domain.AuditEvent, workflow domain.Workflow, attachment domain.Attachment) error {
	return s.update(func(tx *bbolt.Tx) error {
		if err := put(tx, bucketRecords, record.ID, record); err != nil {
			return err
		}
		if err := put(tx, bucketEvents, event.ID, event); err != nil {
			return err
		}
		if err := put(tx, bucketWorkflows, workflow.ID, workflow); err != nil {
			return err
		}
		return put(tx, bucketAttachments, attachment.ID, attachment)
	})
}
