package service

import (
	"context"
	"fmt"

	"example.com/nursery-cms/domain"
)

func (s *Service) ValidateForReview(ctx context.Context, recordID string, policy domain.ReviewPolicy) ([]string, error) {
	if err := contextErr(ctx); err != nil {
		return nil, err
	}
	record, err := s.Store.GetRecord(recordID)
	if err != nil {
		return nil, err
	}
	return policy.ValidateRecord(record), nil
}

func (s *Service) StartPolicyWorkflow(ctx context.Context, batchID, kind, owner string, policy domain.ReviewPolicy) (domain.Workflow, error) {
	if !policy.CanStart(kind) {
		return domain.Workflow{}, fmt.Errorf("workflow kind %s is not allowed", kind)
	}
	workflow, err := s.StartWorkflow(ctx, batchID, kind, owner)
	if err != nil {
		return workflow, err
	}
	if err := policy.CheckWorkflow(workflow); err != nil {
		return workflow, err
	}
	return workflow, nil
}

func (s *Service) ApplyNextStatus(ctx context.Context, recordID, actor string) (domain.Record, error) {
	if err := contextErr(ctx); err != nil {
		return domain.Record{}, err
	}
	record, err := s.Store.GetRecord(recordID)
	if err != nil {
		return record, err
	}
	next, err := domain.NextStatus(record.Status)
	if err != nil {
		return record, err
	}
	if err := record.Transition(next); err != nil {
		return record, err
	}
	if err := s.Store.SaveRecord(record); err != nil {
		return record, err
	}
	if err := s.appendEvent(record.BatchID, record.ID, "status_"+string(next), actor, "policy transition"); err != nil {
		return record, err
	}
	return record, nil
}

func (s *Service) ValidateAll(ctx context.Context, batchID string, policy domain.ReviewPolicy) (map[string][]string, error) {
	items, err := s.Search(ctx, batchID, "")
	if err != nil {
		return nil, err
	}
	result := make(map[string][]string)
	for _, item := range items {
		issues := policy.ValidateRecord(item)
		if len(issues) > 0 {
			result[item.ID] = issues
		}
	}
	return result, nil
}

func (s *Service) RequireReady(ctx context.Context, recordID string) error {
	issues, err := s.ValidateForReview(ctx, recordID, domain.DefaultReviewPolicy())
	if err != nil {
		return err
	}
	if len(issues) > 0 {
		return fmt.Errorf("record not ready: %v", issues)
	}
	return nil
}
