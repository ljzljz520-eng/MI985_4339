package domain

import (
	"errors"
	"fmt"
	"strings"
)

type RecordStatus string

const (
	StatusDraft     RecordStatus = "draft"
	StatusReview    RecordStatus = "review"
	StatusApproved  RecordStatus = "approved"
	StatusPublished RecordStatus = "published"
	StatusArchived  RecordStatus = "archived"
)

type Record struct {
	ID        string       `json:"id"`
	BatchID   string       `json:"batch_id"`
	Title     string       `json:"title"`
	Content   string       `json:"content"`
	Status    RecordStatus `json:"status"`
	Version   int          `json:"version"`
	Owner     string       `json:"owner"`
	UpdatedAt string       `json:"updated_at"`
}

type AuditEvent struct {
	ID        string `json:"id"`
	RecordID  string `json:"record_id"`
	BatchID   string `json:"batch_id"`
	Action    string `json:"action"`
	Actor     string `json:"actor"`
	Detail    string `json:"detail"`
	CreatedAt string `json:"created_at"`
}

type Workflow struct {
	ID          string `json:"id"`
	BatchID     string `json:"batch_id"`
	Kind        string `json:"kind"`
	State       string `json:"state"`
	Owner       string `json:"owner"`
	StartedAt   string `json:"started_at"`
	CompletedAt string `json:"completed_at"`
}

type Attachment struct {
	ID        string `json:"id"`
	RecordID  string `json:"record_id"`
	Name      string `json:"name"`
	MediaType string `json:"media_type"`
	Size      int64  `json:"size"`
	Digest    string `json:"digest"`
}

type ImportedRecord struct {
	ExternalID string `json:"external_id"`
	BatchID    string `json:"batch_id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Owner      string `json:"owner"`
}

type ImportResult struct {
	BatchID   string
	Accepted  []Record
	Rejected  []string
	Cancelled bool
	Message   string
}

var (
	ErrInvalidTransition = errors.New("invalid record transition")
	ErrConflict          = errors.New("record version conflict")
	ErrNotFound          = errors.New("entity not found")
)

func (r Record) Validate() error {
	if strings.TrimSpace(r.ID) == "" {
		return errors.New("record id is required")
	}
	if strings.TrimSpace(r.BatchID) == "" {
		return errors.New("batch id is required")
	}
	if strings.TrimSpace(r.Title) == "" {
		return errors.New("title is required")
	}
	if strings.TrimSpace(r.Content) == "" {
		return errors.New("content is required")
	}
	if r.Version < 1 {
		return errors.New("version must be positive")
	}
	if r.Owner == "" {
		return errors.New("owner is required")
	}
	return nil
}

func (r Record) Summary() string {
	return fmt.Sprintf("%s[%s] %s v%d", r.BatchID, r.Status, r.Title, r.Version)
}

func NewRecord(id, batchID, title, content, owner string) Record {
	return Record{ID: id, BatchID: batchID, Title: title, Content: content, Status: StatusDraft, Version: 1, Owner: owner, UpdatedAt: "fixed"}
}

func NewWorkflow(id, batchID, kind, owner string) Workflow {
	return Workflow{ID: id, BatchID: batchID, Kind: kind, State: "started", Owner: owner, StartedAt: "fixed"}
}

func NewAttachment(id, recordID, name, mediaType string, size int64, digest string) Attachment {
	return Attachment{ID: id, RecordID: recordID, Name: name, MediaType: mediaType, Size: size, Digest: digest}
}
