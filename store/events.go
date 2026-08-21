package store

import (
	"sort"

	"example.com/nursery-cms/domain"
	"go.etcd.io/bbolt"
)

func (s *Store) SaveAuditEvent(event domain.AuditEvent) error {
	if event.ID == "" || event.Action == "" {
		return domain.ErrNotFound
	}
	return s.update(func(tx *bbolt.Tx) error { return put(tx, bucketEvents, event.ID, event) })
}

func (s *Store) ListAuditEvents(batchID string) ([]domain.AuditEvent, error) {
	result := make([]domain.AuditEvent, 0)
	err := s.view(func(tx *bbolt.Tx) error {
		return list(tx, bucketEvents, func(data []byte) error {
			var event domain.AuditEvent
			if err := decode(data, &event); err != nil {
				return err
			}
			if batchID == "" || event.BatchID == batchID {
				result = append(result, event)
			}
			return nil
		})
	})
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, err
}

func (s *Store) SaveWorkflow(workflow domain.Workflow) error {
	if workflow.ID == "" || workflow.BatchID == "" {
		return domain.ErrNotFound
	}
	return s.update(func(tx *bbolt.Tx) error { return put(tx, bucketWorkflows, workflow.ID, workflow) })
}

func (s *Store) GetWorkflow(id string) (domain.Workflow, error) {
	var workflow domain.Workflow
	err := s.view(func(tx *bbolt.Tx) error { return get(tx, bucketWorkflows, id, &workflow) })
	return workflow, err
}

func (s *Store) ListWorkflows(batchID string) ([]domain.Workflow, error) {
	result := make([]domain.Workflow, 0)
	err := s.view(func(tx *bbolt.Tx) error {
		return list(tx, bucketWorkflows, func(data []byte) error {
			var item domain.Workflow
			if err := decode(data, &item); err != nil {
				return err
			}
			if batchID == "" || item.BatchID == batchID {
				result = append(result, item)
			}
			return nil
		})
	})
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, err
}
