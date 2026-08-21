package service

import (
	"context"
	"fmt"

	"example.com/nursery-cms/domain"
)

type Query struct {
	BatchID string
	Status  domain.RecordStatus
	Owner   string
	Text    string
}

func (s *Service) Query(ctx context.Context, query Query) ([]domain.Record, error) {
	items, err := s.Search(ctx, query.BatchID, query.Status)
	if err != nil {
		return nil, err
	}
	filtered := make([]domain.Record, 0, len(items))
	for _, item := range items {
		if query.Owner != "" && item.Owner != query.Owner {
			continue
		}
		if query.Text != "" && !contains(item.Title, query.Text) && !contains(item.Content, query.Text) {
			continue
		}
		filtered = append(filtered, item)
	}
	return SortRecords(filtered), nil
}

func contains(value, needle string) bool {
	if needle == "" {
		return true
	}
	for i := 0; i+len(needle) <= len(value); i++ {
		if value[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}

func (s *Service) UpdateTitle(ctx context.Context, id, title string, expectedVersion int) (domain.Record, error) {
	record, err := s.Store.GetRecord(id)
	if err != nil {
		return record, err
	}
	if record.Version != expectedVersion {
		return record, domain.ErrConflict
	}
	if !domain.CanEdit(record.Status) {
		return record, fmt.Errorf("record cannot be edited")
	}
	if title == "" {
		return record, fmt.Errorf("title cannot be empty")
	}
	record.Title = domain.NormalizeTitle(title)
	record.Version++
	if err := s.Store.SaveRecord(record); err != nil {
		return record, err
	}
	return record, s.appendEvent(record.BatchID, record.ID, "title_updated", record.Owner, record.Title)
}

func (s *Service) UpdateContent(ctx context.Context, id, content string, expectedVersion int) (domain.Record, error) {
	if err := contextErr(ctx); err != nil {
		return domain.Record{}, err
	}
	record, err := s.Store.GetRecord(id)
	if err != nil {
		return record, err
	}
	if record.Version != expectedVersion {
		return record, domain.ErrConflict
	}
	if !domain.CanEdit(record.Status) {
		return record, fmt.Errorf("record cannot be edited")
	}
	if content == "" {
		return record, fmt.Errorf("content cannot be empty")
	}
	record.Content = content
	record.Version++
	if err := s.Store.SaveRecord(record); err != nil {
		return record, err
	}
	return record, s.appendEvent(record.BatchID, record.ID, "content_updated", record.Owner, "collaborator update")
}

func (s *Service) CountByStatus(ctx context.Context, batchID string) (map[domain.RecordStatus]int, error) {
	items, err := s.Search(ctx, batchID, "")
	if err != nil {
		return nil, err
	}
	counts := make(map[domain.RecordStatus]int)
	for _, item := range items {
		counts[item.Status]++
	}
	return counts, nil
}
