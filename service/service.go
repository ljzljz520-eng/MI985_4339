package service

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"example.com/nursery-cms/domain"
	"example.com/nursery-cms/reporting"
	"example.com/nursery-cms/store"
	"example.com/nursery-cms/validation"
)

type Service struct {
	Store     *store.Store
	Validator validation.BatchValidator
}

func New(st *store.Store) *Service {
	return &Service{Store: st, Validator: validation.NewBatchValidator()}
}

func (s *Service) CreateRecord(ctx context.Context, record domain.Record) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	if err := record.Validate(); err != nil {
		return err
	}
	if err := s.Store.SaveRecord(record); err != nil {
		return err
	}
	return s.appendEvent(record.BatchID, record.ID, "created", record.Owner, record.Summary())
}

func (s *Service) SubmitForReview(ctx context.Context, id, actor string) (domain.Record, error) {
	if err := contextErr(ctx); err != nil {
		return domain.Record{}, err
	}
	record, err := s.Store.GetRecord(id)
	if err != nil {
		return record, err
	}
	if err := record.Transition(domain.StatusReview); err != nil {
		return record, err
	}
	if err := s.Store.SaveRecord(record); err != nil {
		return record, err
	}
	return record, s.appendEvent(record.BatchID, record.ID, "submitted", actor, "review requested")
}

func (s *Service) Confirm(ctx context.Context, id, actor string) (domain.Record, error) {
	if err := contextErr(ctx); err != nil {
		return domain.Record{}, err
	}
	record, err := s.Store.GetRecord(id)
	if err != nil {
		return record, err
	}
	if err := record.Transition(domain.StatusApproved); err != nil {
		return record, err
	}
	if err := s.Store.SaveRecord(record); err != nil {
		return record, err
	}
	return record, s.appendEvent(record.BatchID, record.ID, "confirmed", actor, "review approved")
}

func (s *Service) Publish(ctx context.Context, id, actor string) (domain.Record, error) {
	if err := contextErr(ctx); err != nil {
		return domain.Record{}, err
	}
	record, err := s.Store.GetRecord(id)
	if err != nil {
		return record, err
	}
	if err := record.Transition(domain.StatusPublished); err != nil {
		return record, err
	}
	if err := s.Store.SaveRecord(record); err != nil {
		return record, err
	}
	return record, s.appendEvent(record.BatchID, record.ID, "published", actor, "content published")
}

func (s *Service) Archive(ctx context.Context, id, actor string) (domain.Record, error) {
	if err := contextErr(ctx); err != nil {
		return domain.Record{}, err
	}
	record, err := s.Store.GetRecord(id)
	if err != nil {
		return record, err
	}
	if err := record.Transition(domain.StatusArchived); err != nil {
		return record, err
	}
	if err := s.Store.SaveRecord(record); err != nil {
		return record, err
	}
	return record, s.appendEvent(record.BatchID, record.ID, "archived", actor, "content archived")
}

func (s *Service) appendEvent(batchID, recordID, action, actor, detail string) error {
	events, err := s.Store.ListAuditEvents(batchID)
	if err != nil {
		return err
	}
	event := domain.AuditEvent{ID: domain.EventID(batchID, len(events)+1), RecordID: recordID, BatchID: batchID, Action: action, Actor: actor, Detail: detail, CreatedAt: "fixed"}
	return s.Store.SaveAuditEvent(event)
}

func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("context is nil")
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (s *Service) StartWorkflow(ctx context.Context, batchID, kind, owner string) (domain.Workflow, error) {
	if err := contextErr(ctx); err != nil {
		return domain.Workflow{}, err
	}
	if err := domain.ValidateBatchID(batchID); err != nil {
		return domain.Workflow{}, err
	}
	w := domain.NewWorkflow(domain.WorkflowID(batchID), batchID, kind, owner)
	if err := s.Store.SaveWorkflow(w); err != nil {
		return w, err
	}
	return w, s.appendEvent(batchID, "", "workflow_started", owner, kind)
}

func (s *Service) CloseWorkflow(ctx context.Context, batchID, state string) (domain.Workflow, error) {
	if err := contextErr(ctx); err != nil {
		return domain.Workflow{}, err
	}
	w, err := s.Store.GetWorkflow(domain.WorkflowID(batchID))
	if err != nil {
		return w, err
	}
	if err := w.Close(state); err != nil {
		return w, err
	}
	if err := s.Store.SaveWorkflow(w); err != nil {
		return w, err
	}
	return w, s.appendEvent(batchID, "", "workflow_"+state, w.Owner, "workflow closed")
}

func (s *Service) Search(ctx context.Context, batchID string, status domain.RecordStatus) ([]domain.Record, error) {
	if err := contextErr(ctx); err != nil {
		return nil, err
	}
	items, err := s.Store.ListRecords(batchID)
	if err != nil {
		return nil, err
	}
	if status == "" {
		return items, nil
	}
	filtered := items[:0]
	for _, item := range items {
		if item.Status == status {
			filtered = append(filtered, item)
		}
	}
	return filtered, nil
}

func (s *Service) AddAttachment(ctx context.Context, attachment domain.Attachment) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	if err := s.Store.SaveAttachment(attachment); err != nil {
		return err
	}
	return s.appendEvent("attachment", attachment.RecordID, "attachment_added", "system", attachment.Name)
}

func SortRecords(records []domain.Record) []domain.Record {
	copyOf := append([]domain.Record(nil), records...)
	sort.Slice(copyOf, func(i, j int) bool {
		if copyOf[i].BatchID == copyOf[j].BatchID {
			return copyOf[i].ID < copyOf[j].ID
		}
		return copyOf[i].BatchID < copyOf[j].BatchID
	})
	return copyOf
}

func EnsureDistinctBatch(results []domain.ImportResult) error {
	seen := make(map[string]bool)
	for _, result := range results {
		if seen[result.BatchID] {
			return fmt.Errorf("duplicate batch result: %s", result.BatchID)
		}
		seen[result.BatchID] = true
	}
	return nil
}

var _ = reporting.Report{}
