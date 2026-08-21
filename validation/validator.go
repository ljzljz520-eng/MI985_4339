package validation

import (
	"fmt"
	"strings"

	"example.com/nursery-cms/domain"
)

type Issue struct {
	ExternalID string
	Field      string
	Message    string
}

type Validator struct {
	MaxTitle   int
	MaxContent int
}

func New() Validator {
	return Validator{MaxTitle: 120, MaxContent: 5000}
}

func (v Validator) ValidateImported(item domain.ImportedRecord) []Issue {
	issues := make([]Issue, 0)
	if strings.TrimSpace(item.ExternalID) == "" {
		issues = append(issues, Issue{item.ExternalID, "external_id", "required"})
	}
	if err := domain.ValidateBatchID(item.BatchID); err != nil {
		issues = append(issues, Issue{item.ExternalID, "batch_id", err.Error()})
	}
	if strings.TrimSpace(item.Title) == "" {
		issues = append(issues, Issue{item.ExternalID, "title", "required"})
	}
	if len([]rune(item.Title)) > v.MaxTitle {
		issues = append(issues, Issue{item.ExternalID, "title", "too long"})
	}
	if strings.TrimSpace(item.Content) == "" {
		issues = append(issues, Issue{item.ExternalID, "content", "required"})
	}
	if len([]rune(item.Content)) > v.MaxContent {
		issues = append(issues, Issue{item.ExternalID, "content", "too long"})
	}
	if strings.TrimSpace(item.Owner) == "" {
		issues = append(issues, Issue{item.ExternalID, "owner", "required"})
	}
	return issues
}

func (v Validator) Valid(item domain.ImportedRecord) bool {
	return len(v.ValidateImported(item)) == 0
}

func (v Validator) Explain(item domain.ImportedRecord) string {
	issues := v.ValidateImported(item)
	if len(issues) == 0 {
		return "valid"
	}
	parts := make([]string, 0, len(issues))
	for _, issue := range issues {
		parts = append(parts, fmt.Sprintf("%s:%s", issue.Field, issue.Message))
	}
	return strings.Join(parts, ",")
}

func IsRecoverable(issue Issue) bool {
	if issue.Field == "content" && issue.Message == "too long" {
		return true
	}
	if issue.Field == "title" && issue.Message == "too long" {
		return true
	}
	return false
}

func GroupIssues(issues []Issue) map[string][]Issue {
	grouped := make(map[string][]Issue)
	for _, issue := range issues {
		grouped[issue.ExternalID] = append(grouped[issue.ExternalID], issue)
	}
	return grouped
}
