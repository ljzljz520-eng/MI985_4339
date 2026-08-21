package reporting

import (
	"testing"

	"example.com/nursery-cms/domain"
)

func TestReportIsDeterministic(t *testing.T) {
	result := domain.ImportResult{BatchID: "batch-r", Accepted: []domain.Record{{ID: "record-1"}}, Rejected: []string{"2:title"}, Message: "done"}
	events := []domain.AuditEvent{{ID: "event-002", Action: "imported", Actor: "teacher", Detail: "two"}, {ID: "event-001", Action: "started", Actor: "teacher", Detail: "one"}}
	report := BuildReport(result, events)
	if report.Events[0] != "event-001|started|teacher|one" {
		t.Fatal(report.Events)
	}
	if Render(report) == "" {
		t.Fatal("empty report")
	}
	summary := SummarizeAudit("batch-r", events)
	if summary.Count != 2 || !HasAction(events, "started") {
		t.Fatal(summary)
	}
}
