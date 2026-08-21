package reporting

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"example.com/nursery-cms/domain"
)

type ExportRow struct {
	ID      string
	BatchID string
	Title   string
	Status  string
	Owner   string
	Version int
}

func Rows(records []domain.Record) []ExportRow {
	rows := make([]ExportRow, 0, len(records))
	for _, record := range records {
		rows = append(rows, ExportRow{ID: record.ID, BatchID: record.BatchID, Title: record.Title, Status: string(record.Status), Owner: record.Owner, Version: record.Version})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].ID < rows[j].ID })
	return rows
}

func CSV(records []domain.Record) string {
	lines := []string{"id,batch_id,title,status,owner,version"}
	for _, row := range Rows(records) {
		lines = append(lines, fmt.Sprintf("%s,%s,%s,%s,%s,%d", quote(row.ID), quote(row.BatchID), quote(row.Title), quote(row.Status), quote(row.Owner), row.Version))
	}
	return strings.Join(lines, "\n")
}

func quote(value string) string {
	if strings.ContainsAny(value, ",\"\n") {
		return "\"" + strings.ReplaceAll(value, "\"", "\"\"") + "\""
	}
	return value
}

func JSON(records []domain.Record) ([]byte, error) {
	return json.Marshal(Rows(records))
}

func Markdown(report Report) string {
	lines := []string{"# " + report.BatchID, "", fmt.Sprintf("Accepted: %d", report.Accepted), fmt.Sprintf("Rejected: %d", report.Rejected)}
	if report.Cancelled {
		lines = append(lines, "State: cancelled")
	} else {
		lines = append(lines, "State: completed")
	}
	if report.Message != "" {
		lines = append(lines, "Message: "+report.Message)
	}
	if len(report.Events) > 0 {
		lines = append(lines, "", "## Events", "")
	}
	for _, event := range report.Events {
		lines = append(lines, "- "+event)
	}
	return strings.Join(lines, "\n")
}

func CompareReports(left, right Report) bool {
	return Render(left) == Render(right)
}
