package service

import (
	"context"
	"fmt"

	"example.com/nursery-cms/domain"
)

func (s *Service) ArchiveBatch(ctx context.Context, batchID, actor string) ([]domain.Record, error) {
	items, err := s.Search(ctx, batchID, domain.StatusPublished)
	if err != nil {
		return nil, err
	}
	archived := make([]domain.Record, 0, len(items))
	for _, item := range items {
		updated, transitionErr := s.Archive(ctx, item.ID, actor)
		if transitionErr != nil {
			return archived, transitionErr
		}
		archived = append(archived, updated)
	}
	return archived, nil
}

func (s *Service) Reassign(ctx context.Context, batchID, from, to string) (int, error) {
	if to == "" {
		return 0, fmt.Errorf("new owner is required")
	}
	items, err := s.Search(ctx, batchID, "")
	if err != nil {
		return 0, err
	}
	changed := 0
	for _, item := range items {
		if from != "" && item.Owner != from {
			continue
		}
		if !domain.CanEdit(item.Status) {
			continue
		}
		item.Owner = to
		item.Version++
		if err := s.Store.SaveRecord(item); err != nil {
			return changed, err
		}
		if err := s.appendEvent(batchID, item.ID, "reassigned", to, from+"->"+to); err != nil {
			return changed, err
		}
		changed++
	}
	return changed, nil
}

func (s *Service) RemoveDrafts(ctx context.Context, batchID string) (int, error) {
	items, err := s.Search(ctx, batchID, domain.StatusDraft)
	if err != nil {
		return 0, err
	}
	removed := 0
	for _, item := range items {
		if err := s.Store.DeleteRecord(item.ID); err != nil {
			return removed, err
		}
		removed++
	}
	return removed, nil
}

func (s *Service) EnsureWorkflowClosed(batchID string) error {
	state, err := s.BatchState(batchID)
	if err != nil {
		return err
	}
	if !domain.WorkflowCanClose(state) {
		return fmt.Errorf("workflow %s remains %s", batchID, state)
	}
	return nil
}
