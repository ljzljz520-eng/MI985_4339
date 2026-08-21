package reporting

import (
	"fmt"
	"sort"
	"strings"

	"example.com/nursery-cms/domain"
)

type AuditSummary struct {
	BatchID string
	Actors  []string
	Actions []string
	Count   int
}

func SummarizeAudit(batchID string, events []domain.AuditEvent) AuditSummary {
	summary := AuditSummary{BatchID: batchID, Count: len(events)}
	actors := make(map[string]bool)
	actions := make(map[string]bool)
	for _, event := range events {
		actors[event.Actor] = true
		actions[event.Action] = true
	}
	for actor := range actors {
		summary.Actors = append(summary.Actors, actor)
	}
	for action := range actions {
		summary.Actions = append(summary.Actions, action)
	}
	sort.Strings(summary.Actors)
	sort.Strings(summary.Actions)
	return summary
}

func (s AuditSummary) String() string {
	return fmt.Sprintf("batch=%s count=%d actors=%s actions=%s", s.BatchID, s.Count, strings.Join(s.Actors, ","), strings.Join(s.Actions, ","))
}

func EventCount(events []domain.AuditEvent, action string) int {
	count := 0
	for _, event := range events {
		if event.Action == action {
			count++
		}
	}
	return count
}

func HasAction(events []domain.AuditEvent, action string) bool {
	return EventCount(events, action) > 0
}

func Timeline(events []domain.AuditEvent) []string {
	copyOf := append([]domain.AuditEvent(nil), events...)
	sort.Slice(copyOf, func(i, j int) bool { return copyOf[i].ID < copyOf[j].ID })
	lines := make([]string, 0, len(copyOf))
	for _, event := range copyOf {
		lines = append(lines, FormatEvent(event))
	}
	return lines
}
