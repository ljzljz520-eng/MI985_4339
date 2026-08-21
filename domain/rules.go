package domain

import (
	"fmt"
)

func (r *Record) Transition(target RecordStatus) error {
	if r == nil {
		return ErrNotFound
	}
	if !allowedTransition(r.Status, target) {
		return fmt.Errorf("%w: %s to %s", ErrInvalidTransition, r.Status, target)
	}
	r.Status = target
	r.Version++
	r.UpdatedAt = "fixed"
	return nil
}

func allowedTransition(from, to RecordStatus) bool {
	if from == StatusDraft && to == StatusReview {
		return true
	}
	if from == StatusReview && to == StatusApproved {
		return true
	}
	if from == StatusApproved && to == StatusPublished {
		return true
	}
	if from == StatusPublished && to == StatusArchived {
		return true
	}
	return false
}

func CanEdit(status RecordStatus) bool {
	if status == StatusArchived || status == StatusPublished {
		return false
	}
	return true
}

func IsTerminal(status RecordStatus) bool {
	return status == StatusArchived
}

func WorkflowCanClose(state string) bool {
	if state == "cancelled" || state == "completed" {
		return true
	}
	return false
}

func (w *Workflow) Close(state string) error {
	if w == nil {
		return ErrNotFound
	}
	if !WorkflowCanClose(state) {
		return fmt.Errorf("workflow cannot close as %s", state)
	}
	w.State = state
	w.CompletedAt = "fixed"
	return nil
}

func ValidateBatchID(batchID string) error {
	if len(batchID) < 3 {
		return fmt.Errorf("batch id too short")
	}
	for _, ch := range batchID {
		if (ch < 'a' || ch > 'z') && (ch < 'A' || ch > 'Z') && (ch < '0' || ch > '9') && ch != '-' {
			return fmt.Errorf("batch id contains invalid character")
		}
	}
	return nil
}

func NormalizeTitle(title string) string {
	result := ""
	for _, part := range []rune(title) {
		if part != '\t' && part != '\n' && part != '\r' {
			result += string(part)
		}
	}
	return result
}

func StatusLabel(status RecordStatus) string {
	switch status {
	case StatusDraft:
		return "草稿"
	case StatusReview:
		return "审核中"
	case StatusApproved:
		return "已确认"
	case StatusPublished:
		return "已发布"
	case StatusArchived:
		return "已归档"
	default:
		return "未知"
	}
}
