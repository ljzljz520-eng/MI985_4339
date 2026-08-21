package store

import (
	"sort"

	"example.com/nursery-cms/domain"
	"go.etcd.io/bbolt"
)

func (s *Store) SaveRecord(record domain.Record) error {
	if err := record.Validate(); err != nil {
		return err
	}
	return s.update(func(tx *bbolt.Tx) error { return put(tx, bucketRecords, record.ID, record) })
}

func (s *Store) GetRecord(id string) (domain.Record, error) {
	var record domain.Record
	err := s.view(func(tx *bbolt.Tx) error { return get(tx, bucketRecords, id, &record) })
	return record, err
}

func (s *Store) ListRecords(batchID string) ([]domain.Record, error) {
	result := make([]domain.Record, 0)
	err := s.view(func(tx *bbolt.Tx) error {
		return list(tx, bucketRecords, func(data []byte) error {
			var record domain.Record
			if err := decode(data, &record); err != nil {
				return err
			}
			if batchID == "" || record.BatchID == batchID {
				result = append(result, record)
			}
			return nil
		})
	})
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, err
}

func (s *Store) DeleteRecord(id string) error {
	return s.update(func(tx *bbolt.Tx) error { return tx.Bucket(bucketRecords).Delete([]byte(id)) })
}

func (s *Store) SaveAttachment(attachment domain.Attachment) error {
	if attachment.ID == "" || attachment.RecordID == "" || attachment.Name == "" {
		return domain.ErrNotFound
	}
	return s.update(func(tx *bbolt.Tx) error { return put(tx, bucketAttachments, attachment.ID, attachment) })
}

func (s *Store) ListAttachments(recordID string) ([]domain.Attachment, error) {
	result := make([]domain.Attachment, 0)
	err := s.view(func(tx *bbolt.Tx) error {
		return list(tx, bucketAttachments, func(data []byte) error {
			var item domain.Attachment
			if err := decode(data, &item); err != nil {
				return err
			}
			if recordID == "" || item.RecordID == recordID {
				result = append(result, item)
			}
			return nil
		})
	})
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, err
}
