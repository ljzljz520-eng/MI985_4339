package reporting

import (
	"fmt"
	"sort"
	"strings"

	"example.com/nursery-cms/domain"
)

type Report struct {
	BatchID   string
	Accepted  int
	Rejected  int
	Cancelled bool
	Message   string
	Events    []string
}

func BuildReport(result domain.ImportResult, events []domain.AuditEvent) Report {
	report := Report{BatchID: result.BatchID, Accepted: len(result.Accepted), Rejected: len(result.Rejected), Cancelled: result.Cancelled, Message: result.Message}
	sorted := append([]domain.AuditEvent(nil), events...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].ID < sorted[j].ID })
	for _, event := range sorted {
		report.Events = append(report.Events, FormatEvent(event))
	}
	return report
}

func FormatEvent(event domain.AuditEvent) string {
	return fmt.Sprintf("%s|%s|%s|%s", event.ID, event.Action, event.Actor, event.Detail)
}

func Render(report Report) string {
	state := "completed"
	if report.Cancelled {
		state = "cancelled"
	}
	lines := []string{fmt.Sprintf("batch=%s", report.BatchID), fmt.Sprintf("state=%s", state), fmt.Sprintf("accepted=%d", report.Accepted), fmt.Sprintf("rejected=%d", report.Rejected), "message=" + report.Message}
	lines = append(lines, report.Events...)
	return strings.Join(lines, "\n")
}

func SummarizeRecords(records []domain.Record) string {
	if len(records) == 0 {
		return "no records"
	}
	parts := make([]string, 0, len(records))
	for _, record := range records {
		parts = append(parts, record.Summary())
	}
	sort.Strings(parts)
	return strings.Join(parts, ";")
}

func AcceptedIDs(result domain.ImportResult) []string {
	ids := make([]string, 0, len(result.Accepted))
	for _, record := range result.Accepted {
		ids = append(ids, record.ID)
	}
	sort.Strings(ids)
	return ids
}
