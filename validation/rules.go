package validation

import (
	"fmt"
	"strings"
	"unicode"

	"example.com/nursery-cms/domain"
)

type BatchRules struct {
	RequireSequentialIDs bool
	RejectDuplicateTitle bool
	MaxItems             int
}

func DefaultBatchRules() BatchRules {
	return BatchRules{RequireSequentialIDs: false, RejectDuplicateTitle: true, MaxItems: 500}
}

func (r BatchRules) ValidateBatch(items []domain.ImportedRecord) []Issue {
	issues := make([]Issue, 0)
	if len(items) == 0 {
		return append(issues, Issue{Field: "batch", Message: "empty"})
	}
	if len(items) > r.MaxItems {
		issues = append(issues, Issue{Field: "batch", Message: "too many items"})
	}
	seenIDs := make(map[string]bool)
	seenTitles := make(map[string]string)
	for index, item := range items {
		if seenIDs[item.ExternalID] {
			issues = append(issues, Issue{ExternalID: item.ExternalID, Field: "external_id", Message: "duplicate"})
		}
		seenIDs[item.ExternalID] = true
		title := strings.TrimSpace(item.Title)
		if r.RejectDuplicateTitle && title != "" {
			if previous, exists := seenTitles[title]; exists {
				issues = append(issues, Issue{ExternalID: item.ExternalID, Field: "title", Message: "duplicate with " + previous})
			}
			seenTitles[title] = item.ExternalID
		}
		if r.RequireSequentialIDs && item.ExternalID != fmt.Sprintf("%d", index+1) {
			issues = append(issues, Issue{ExternalID: item.ExternalID, Field: "external_id", Message: "not sequential"})
		}
	}
	return issues
}

func NormalizeOwner(owner string) string {
	parts := strings.Fields(owner)
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, " ")
}

func HasUnsafeControl(value string) bool {
	for _, runeValue := range value {
		if unicode.IsControl(runeValue) && runeValue != '\n' && runeValue != '\t' {
			return true
		}
	}
	return false
}

func SanitizeContent(value string) string {
	var builder strings.Builder
	for _, runeValue := range value {
		if unicode.IsControl(runeValue) && runeValue != '\n' && runeValue != '\t' {
			continue
		}
		builder.WriteRune(runeValue)
	}
	return strings.TrimSpace(builder.String())
}

func ValidateAttachment(name, mediaType string, size int64) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("attachment name is required")
	}
	if !strings.Contains(mediaType, "/") {
		return fmt.Errorf("media type is invalid")
	}
	if size < 0 {
		return fmt.Errorf("attachment size cannot be negative")
	}
	return nil
}

func UniqueBatchIDs(items []domain.ImportedRecord) []string {
	seen := make(map[string]bool)
	result := make([]string, 0)
	for _, item := range items {
		if item.BatchID == "" || seen[item.BatchID] {
			continue
		}
		seen[item.BatchID] = true
		result = append(result, item.BatchID)
	}
	return result
}
