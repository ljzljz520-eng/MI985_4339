package service

import (
	"context"
	"fmt"
	"sort"

	"example.com/nursery-cms/domain"
)

type BatchQueue struct {
	batches map[string][]domain.ImportedRecord
}

func NewBatchQueue() *BatchQueue {
	return &BatchQueue{batches: make(map[string][]domain.ImportedRecord)}
}

func (q *BatchQueue) Enqueue(batchID string, item domain.ImportedRecord) error {
	if q == nil {
		return fmt.Errorf("queue is nil")
	}
	if err := domain.ValidateBatchID(batchID); err != nil {
		return err
	}
	if item.BatchID != batchID {
		return fmt.Errorf("item batch mismatch")
	}
	q.batches[batchID] = append(q.batches[batchID], item)
	return nil
}

func (q *BatchQueue) EnqueueMany(batchID string, items []domain.ImportedRecord) error {
	for _, item := range items {
		if err := q.Enqueue(batchID, item); err != nil {
			return err
		}
	}
	return nil
}

func (q *BatchQueue) Peek(batchID string) (domain.ImportedRecord, bool) {
	if q == nil {
		return domain.ImportedRecord{}, false
	}
	items := q.batches[batchID]
	if len(items) == 0 {
		return domain.ImportedRecord{}, false
	}
	return items[0], true
}

func (q *BatchQueue) Drain(batchID string) []domain.ImportedRecord {
	if q == nil {
		return nil
	}
	items := append([]domain.ImportedRecord(nil), q.batches[batchID]...)
	delete(q.batches, batchID)
	return items
}

func (q *BatchQueue) Size(batchID string) int {
	if q == nil {
		return 0
	}
	return len(q.batches[batchID])
}

func (q *BatchQueue) Batches() []string {
	if q == nil {
		return nil
	}
	result := make([]string, 0, len(q.batches))
	for batchID := range q.batches {
		result = append(result, batchID)
	}
	sort.Strings(result)
	return result
}

func (s *Service) ImportQueued(ctx context.Context, queue *BatchQueue, batchID string) (domain.ImportResult, error) {
	if queue == nil {
		return domain.ImportResult{BatchID: batchID}, fmt.Errorf("queue is nil")
	}
	items := queue.Drain(batchID)
	if len(items) == 0 {
		return domain.ImportResult{BatchID: batchID}, fmt.Errorf("batch queue is empty")
	}
	return s.ImportBatch(ctx, batchID, items)
}

func (s *Service) QueueAndImport(ctx context.Context, batchID string, items []domain.ImportedRecord) (domain.ImportResult, error) {
	queue := NewBatchQueue()
	if err := queue.EnqueueMany(batchID, items); err != nil {
		return domain.ImportResult{BatchID: batchID}, err
	}
	return s.ImportQueued(ctx, queue, batchID)
}

func (q *BatchQueue) Validate() error {
	if q == nil {
		return fmt.Errorf("queue is nil")
	}
	for batchID, items := range q.batches {
		if err := domain.ValidateBatchID(batchID); err != nil {
			return err
		}
		for _, item := range items {
			if item.BatchID != batchID {
				return fmt.Errorf("queued item mismatch")
			}
		}
	}
	return nil
}

func (q *BatchQueue) Clone() *BatchQueue {
	clone := NewBatchQueue()
	if q == nil {
		return clone
	}
	for batchID, items := range q.batches {
		clone.batches[batchID] = append([]domain.ImportedRecord(nil), items...)
	}
	return clone
}
