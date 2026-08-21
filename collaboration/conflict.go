package collaboration

import (
	"errors"
	"fmt"

	"example.com/nursery-cms/domain"
)

type Conflict struct {
	RecordID string
	Expected int
	Actual   int
	Detail   string
}

func CheckVersion(record domain.Record, expected int) *Conflict {
	if record.Version == expected {
		return nil
	}
	return &Conflict{RecordID: record.ID, Expected: expected, Actual: record.Version, Detail: "stale editor version"}
}

func Resolve(record domain.Record, expected int, content string) (domain.Record, error) {
	conflict := CheckVersion(record, expected)
	if conflict != nil {
		return record, fmt.Errorf("%w: %s", domain.ErrConflict, conflict.Detail)
	}
	if !domain.CanEdit(record.Status) {
		return record, errors.New("record is not editable")
	}
	if content == "" {
		return record, errors.New("content is required")
	}
	record.Content = content
	record.Version++
	return record, nil
}

func MergeNotes(base, incoming []string) []string {
	seen := make(map[string]bool)
	merged := make([]string, 0, len(base)+len(incoming))
	for _, note := range append(append([]string{}, base...), incoming...) {
		if note == "" || seen[note] {
			continue
		}
		seen[note] = true
		merged = append(merged, note)
	}
	return merged
}

func ConflictMessage(conflict *Conflict) string {
	if conflict == nil {
		return "no conflict"
	}
	return fmt.Sprintf("%s expected=%d actual=%d", conflict.RecordID, conflict.Expected, conflict.Actual)
}
