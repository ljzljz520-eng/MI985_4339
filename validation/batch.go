package validation

import (
	"context"
	"fmt"

	"example.com/nursery-cms/domain"
)

type BatchValidator struct {
	Record Validator
}

func NewBatchValidator() BatchValidator {
	return BatchValidator{Record: New()}
}

func (b BatchValidator) Validate(ctx context.Context, items []domain.ImportedRecord) ([]domain.Record, []Issue, error) {
	accepted := make([]domain.Record, 0, len(items))
	issues := make([]Issue, 0)
	for _, item := range items {
		select {
		case <-ctx.Done():
			return accepted, issues, ctx.Err()
		default:
		}
		found := b.Record.ValidateImported(item)
		if len(found) > 0 {
			issues = append(issues, found...)
			continue
		}
		id := domain.RecordID(item.BatchID, item.ExternalID)
		accepted = append(accepted, domain.NewRecord(id, item.BatchID, domain.NormalizeTitle(item.Title), item.Content, item.Owner))
	}
	if len(accepted) == 0 && len(issues) > 0 {
		return accepted, issues, fmt.Errorf("batch has no valid records")
	}
	return accepted, issues, nil
}

func (b BatchValidator) ValidateOne(ctx context.Context, item domain.ImportedRecord) (domain.Record, []Issue, error) {
	select {
	case <-ctx.Done():
		return domain.Record{}, nil, ctx.Err()
	default:
	}
	issues := b.Record.ValidateImported(item)
	if len(issues) != 0 {
		return domain.Record{}, issues, nil
	}
	return domain.NewRecord(domain.RecordID(item.BatchID, item.ExternalID), item.BatchID, item.Title, item.Content, item.Owner), nil, nil
}

func (b BatchValidator) CountValid(items []domain.ImportedRecord) int {
	count := 0
	for _, item := range items {
		if b.Record.Valid(item) {
			count++
		}
	}
	return count
}
