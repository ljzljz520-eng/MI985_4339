package service

import (
	"context"
	"fmt"
	"sort"

	"example.com/nursery-cms/domain"
)

type BatchMetrics struct {
	BatchID       string
	Total         int
	Draft         int
	Review        int
	Approved      int
	Published     int
	Archived      int
	LatestVersion int
}

func (s *Service) Metrics(ctx context.Context, batchID string) (BatchMetrics, error) {
	items, err := s.Search(ctx, batchID, "")
	if err != nil {
		return BatchMetrics{}, err
	}
	metrics := BatchMetrics{BatchID: batchID, Total: len(items)}
	for _, item := range items {
		switch item.Status {
		case domain.StatusDraft:
			metrics.Draft++
		case domain.StatusReview:
			metrics.Review++
		case domain.StatusApproved:
			metrics.Approved++
		case domain.StatusPublished:
			metrics.Published++
		case domain.StatusArchived:
			metrics.Archived++
		}
		if item.Version > metrics.LatestVersion {
			metrics.LatestVersion = item.Version
		}
	}
	return metrics, nil
}

func (m BatchMetrics) Complete() bool {
	return m.Total > 0 && m.Archived == m.Total
}

func (m BatchMetrics) Pending() int {
	return m.Total - m.Archived
}

func (m BatchMetrics) String() string {
	return fmt.Sprintf("%s total=%d draft=%d review=%d approved=%d published=%d archived=%d", m.BatchID, m.Total, m.Draft, m.Review, m.Approved, m.Published, m.Archived)
}

func (s *Service) Owners(ctx context.Context, batchID string) ([]string, error) {
	items, err := s.Search(ctx, batchID, "")
	if err != nil {
		return nil, err
	}
	seen := make(map[string]bool)
	owners := make([]string, 0)
	for _, item := range items {
		if !seen[item.Owner] {
			seen[item.Owner] = true
			owners = append(owners, item.Owner)
		}
	}
	sort.Strings(owners)
	return owners, nil
}

func (s *Service) AuditCount(batchID, action string) (int, error) {
	events, err := s.Store.ListAuditEvents(batchID)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, event := range events {
		if action == "" || event.Action == action {
			count++
		}
	}
	return count, nil
}

func (s *Service) VerifyBatchIntegrity(ctx context.Context, batchID string) error {
	items, err := s.Search(ctx, batchID, "")
	if err != nil {
		return err
	}
	seen := make(map[string]bool)
	for _, item := range items {
		if seen[item.ID] {
			return fmt.Errorf("duplicate record id %s", item.ID)
		}
		seen[item.ID] = true
		if err := item.Validate(); err != nil {
			return err
		}
	}
	return nil
}
