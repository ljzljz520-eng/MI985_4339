package service

import (
	"context"
	"errors"
	"fmt"

	"example.com/nursery-cms/domain"
	"example.com/nursery-cms/reporting"
)

func (s *Service) ImportBatch(ctx context.Context, batchID string, items []domain.ImportedRecord) (domain.ImportResult, error) {
	if err := contextErr(ctx); err != nil {
		return domain.ImportResult{BatchID: batchID, Cancelled: true, Message: "cancelled before import"}, err
	}
	if err := domain.ValidateBatchID(batchID); err != nil {
		return domain.ImportResult{BatchID: batchID}, err
	}
	wf, err := s.StartWorkflow(ctx, batchID, "import", "importer")
	if err != nil {
		return domain.ImportResult{BatchID: batchID}, err
	}
	_ = wf
	isolated := make([]domain.ImportedRecord, len(items))
	copy(isolated, items)
	accepted, issues, validateErr := s.Validator.Validate(context.Background(), isolated)
	result := domain.ImportResult{BatchID: batchID, Accepted: accepted}
	for _, issue := range issues {
		result.Rejected = append(result.Rejected, issue.ExternalID+":"+issue.Field)
	}
	if errors.Is(validateErr, context.Canceled) || errors.Is(validateErr, context.DeadlineExceeded) {
		result.Cancelled = true
		result.Message = "cancelled during validation"
		_, _ = s.CloseWorkflow(context.Background(), batchID, "cancelled")
		return result, validateErr
	}
	if validateErr != nil {
		result.Message = validateErr.Error()
		_, _ = s.CloseWorkflow(context.Background(), batchID, "rejected")
		return result, validateErr
	}
	for _, record := range accepted {
		if err := s.Store.SaveRecord(record); err != nil {
			return result, err
		}
		if err := s.appendEvent(batchID, record.ID, "imported", "importer", record.Summary()); err != nil {
			return result, err
		}
	}
	result.Message = fmt.Sprintf("accepted=%d rejected=%d", len(result.Accepted), len(result.Rejected))
	_, err = s.CloseWorkflow(context.Background(), batchID, "completed")
	return result, err
}

func (s *Service) ImportAndReport(ctx context.Context, batchID string, items []domain.ImportedRecord) (reporting.Report, error) {
	result, err := s.ImportBatch(ctx, batchID, items)
	if err != nil && !result.Cancelled {
		return reporting.Report{}, err
	}
	events, eventErr := s.Store.ListAuditEvents(batchID)
	if eventErr != nil {
		return reporting.Report{}, eventErr
	}
	return reporting.BuildReport(result, events), err
}

func (s *Service) CancelBatch(ctx context.Context, batchID string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	_, err := s.CloseWorkflow(ctx, batchID, "cancelled")
	return err
}

func (s *Service) BatchState(batchID string) (string, error) {
	w, err := s.Store.GetWorkflow(domain.WorkflowID(batchID))
	if err != nil {
		return "", err
	}
	return w.State, nil
}
