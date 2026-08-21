package domain

import (
	"fmt"
	"strings"
)

type ReviewPolicy struct {
	MinimumContent int
	RequireOwner   bool
	AllowArchive   bool
	AllowedKinds   []string
}

func DefaultReviewPolicy() ReviewPolicy {
	return ReviewPolicy{MinimumContent: 2, RequireOwner: true, AllowArchive: true, AllowedKinds: []string{"import", "manual", "revision"}}
}

func (p ReviewPolicy) ValidateRecord(record Record) []string {
	issues := make([]string, 0)
	if strings.TrimSpace(record.Title) == "" {
		issues = append(issues, "title")
	}
	if len([]rune(strings.TrimSpace(record.Content))) < p.MinimumContent {
		issues = append(issues, "content")
	}
	if p.RequireOwner && strings.TrimSpace(record.Owner) == "" {
		issues = append(issues, "owner")
	}
	return issues
}

func (p ReviewPolicy) CanStart(kind string) bool {
	for _, allowed := range p.AllowedKinds {
		if kind == allowed {
			return true
		}
	}
	return false
}

func (p ReviewPolicy) CanArchive(record Record) bool {
	if !p.AllowArchive {
		return false
	}
	if record.Status != StatusPublished {
		return false
	}
	return true
}

func (p ReviewPolicy) CheckWorkflow(workflow Workflow) error {
	if !p.CanStart(workflow.Kind) {
		return fmt.Errorf("workflow kind %s is not allowed", workflow.Kind)
	}
	if workflow.Owner == "" && p.RequireOwner {
		return fmt.Errorf("workflow owner is required")
	}
	return nil
}

func (p ReviewPolicy) ExplainRecord(record Record) string {
	issues := p.ValidateRecord(record)
	if len(issues) == 0 {
		return "ready"
	}
	return "needs:" + strings.Join(issues, ",")
}

func ParseStatusList(value string) []RecordStatus {
	parts := strings.Split(value, ",")
	statuses := make([]RecordStatus, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		status := RecordStatus(strings.ToLower(trimmed))
		switch status {
		case StatusDraft, StatusReview, StatusApproved, StatusPublished, StatusArchived:
			statuses = append(statuses, status)
		}
	}
	return statuses
}

func StatusMatches(status RecordStatus, accepted []RecordStatus) bool {
	if len(accepted) == 0 {
		return true
	}
	for _, candidate := range accepted {
		if status == candidate {
			return true
		}
	}
	return false
}

func NextStatus(status RecordStatus) (RecordStatus, error) {
	switch status {
	case StatusDraft:
		return StatusReview, nil
	case StatusReview:
		return StatusApproved, nil
	case StatusApproved:
		return StatusPublished, nil
	case StatusPublished:
		return StatusArchived, nil
	default:
		return "", fmt.Errorf("no next status for %s", status)
	}
}
