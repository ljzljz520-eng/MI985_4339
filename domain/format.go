package domain

import (
	"fmt"
	"strings"
)

func FormatRecord(record Record) string {
	return fmt.Sprintf("%s|%s|%s|%s|v%d|%s", record.ID, record.BatchID, record.Title, record.Status, record.Version, record.Owner)
}

func ParseRecord(value string) (Record, error) {
	parts := strings.Split(value, "|")
	if len(parts) < 6 {
		return Record{}, fmt.Errorf("record format is incomplete")
	}
	var version int
	if _, err := fmt.Sscanf(parts[4], "v%d", &version); err != nil {
		return Record{}, err
	}
	record := Record{ID: parts[0], BatchID: parts[1], Title: parts[2], Status: RecordStatus(parts[3]), Version: version, Owner: parts[5], Content: "restored", UpdatedAt: "fixed"}
	if err := record.Validate(); err != nil {
		return Record{}, err
	}
	return record, nil
}

func CloneRecord(record Record) Record { return record }

func CloneRecords(records []Record) []Record {
	result := make([]Record, len(records))
	copy(result, records)
	return result
}

func SameBatch(left, right Record) bool { return left.BatchID != "" && left.BatchID == right.BatchID }

func IsOwnedBy(record Record, owner string) bool { return owner != "" && record.Owner == owner }

func RecordWeight(record Record) int {
	weight := len([]rune(record.Title)) + len([]rune(record.Content))
	if record.Status == StatusPublished {
		weight += 10
	}
	if record.Status == StatusArchived {
		weight += 20
	}
	return weight
}

func CompareRecordVersions(left, right Record) int {
	if left.Version < right.Version {
		return -1
	}
	if left.Version > right.Version {
		return 1
	}
	return 0
}

func MergeRecordMetadata(base, incoming Record) Record {
	merged := base
	if incoming.Title != "" {
		merged.Title = incoming.Title
	}
	if incoming.Content != "" {
		merged.Content = incoming.Content
	}
	if incoming.Owner != "" {
		merged.Owner = incoming.Owner
	}
	if incoming.Version > merged.Version {
		merged.Version = incoming.Version
	}
	return merged
}

func ValidateEntitySet(record Record, event AuditEvent, workflow Workflow, attachment Attachment) error {
	if err := record.Validate(); err != nil {
		return err
	}
	if event.RecordID != "" && event.RecordID != record.ID {
		return fmt.Errorf("event record mismatch")
	}
	if workflow.BatchID != record.BatchID {
		return fmt.Errorf("workflow batch mismatch")
	}
	if attachment.RecordID != "" && attachment.RecordID != record.ID {
		return fmt.Errorf("attachment record mismatch")
	}
	return nil
}

func BatchSummary(records []Record) string {
	if len(records) == 0 {
		return "empty"
	}
	parts := make([]string, 0, len(records))
	for _, record := range records {
		parts = append(parts, record.Summary())
	}
	return strings.Join(parts, ";")
}
