package service

import (
	"context"
	"fmt"

	"example.com/nursery-cms/domain"
)

type Lifecycle struct{ service *Service }

func NewLifecycle(service *Service) Lifecycle { return Lifecycle{service: service} }

func (l Lifecycle) CreateReviewConfirmArchive(ctx context.Context, record domain.Record, reviewer, archivist string) (domain.Record, error) {
	if err := l.service.CreateRecord(ctx, record); err != nil {
		return record, err
	}
	current, err := l.service.SubmitForReview(ctx, record.ID, reviewer)
	if err != nil {
		return current, err
	}
	current, err = l.service.Confirm(ctx, current.ID, reviewer)
	if err != nil {
		return current, err
	}
	current, err = l.service.Publish(ctx, current.ID, reviewer)
	if err != nil {
		return current, err
	}
	return l.service.Archive(ctx, current.ID, archivist)
}

func (l Lifecycle) ReopenableSummary(ctx context.Context, batchID string) (string, error) {
	items, err := l.service.Search(ctx, batchID, "")
	if err != nil {
		return "", err
	}
	if len(items) == 0 {
		return "empty", nil
	}
	return fmt.Sprintf("%s:%d", batchID, len(items)), nil
}

func (l Lifecycle) ValidateState(record domain.Record) error {
	if err := record.Validate(); err != nil {
		return err
	}
	if domain.IsTerminal(record.Status) && domain.CanEdit(record.Status) {
		return fmt.Errorf("terminal record editable")
	}
	return nil
}

func (l Lifecycle) EnsurePublishable(record domain.Record) error {
	if record.Status != domain.StatusApproved {
		return fmt.Errorf("record is not approved")
	}
	if record.Content == "" {
		return fmt.Errorf("content is empty")
	}
	return nil
}
