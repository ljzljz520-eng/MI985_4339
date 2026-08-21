package collaboration

import (
	"fmt"
	"sort"
	"strings"
)

type Note struct {
	ID       string
	RecordID string
	Author   string
	Body     string
	Resolved bool
}

type Notes struct{ values map[string]Note }

func NewNotes() *Notes { return &Notes{values: make(map[string]Note)} }

func (n *Notes) Add(note Note) error {
	if n == nil {
		return fmt.Errorf("notes collection is nil")
	}
	if note.ID == "" || note.RecordID == "" || note.Author == "" {
		return fmt.Errorf("note identity is incomplete")
	}
	if strings.TrimSpace(note.Body) == "" {
		return fmt.Errorf("note body is empty")
	}
	if _, exists := n.values[note.ID]; exists {
		return fmt.Errorf("note already exists")
	}
	n.values[note.ID] = note
	return nil
}

func (n *Notes) Resolve(id string) error {
	if n == nil {
		return fmt.Errorf("notes collection is nil")
	}
	note, ok := n.values[id]
	if !ok {
		return fmt.Errorf("note not found")
	}
	note.Resolved = true
	n.values[id] = note
	return nil
}

func (n *Notes) List(recordID string, includeResolved bool) []Note {
	if n == nil {
		return nil
	}
	result := make([]Note, 0)
	for _, note := range n.values {
		if recordID != "" && note.RecordID != recordID {
			continue
		}
		if !includeResolved && note.Resolved {
			continue
		}
		result = append(result, note)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result
}

func (n *Notes) CountOpen(recordID string) int { return len(n.List(recordID, false)) }

func FormatNotes(notes []Note) string {
	parts := make([]string, 0, len(notes))
	for _, note := range notes {
		state := "open"
		if note.Resolved {
			state = "resolved"
		}
		parts = append(parts, fmt.Sprintf("%s:%s:%s", note.ID, state, note.Body))
	}
	return strings.Join(parts, "\n")
}
