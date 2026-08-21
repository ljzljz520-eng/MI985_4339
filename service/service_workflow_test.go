package service

import (
	"context"
	"testing"

	"example.com/nursery-cms/domain"
	"example.com/nursery-cms/store"
)

func testService(t *testing.T) *Service {
	t.Helper()
	db, err := store.Open(t.TempDir() + "/nursery.db")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return New(db)
}

func TestWorkflowCreateReviewArchive(t *testing.T) {
	s := testService(t)
	record := domain.NewRecord("record-a", "batch-a", "春季活动", "观察植物", "teacher")
	got, err := NewLifecycle(s).CreateReviewConfirmArchive(context.Background(), record, "reviewer", "archivist")
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != domain.StatusArchived {
		t.Fatalf("status=%s", got.Status)
	}
	if got.Version != 5 {
		t.Fatalf("version=%d", got.Version)
	}
	events, err := s.Store.ListAuditEvents("batch-a")
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 5 {
		t.Fatalf("events=%d", len(events))
	}
}

func TestLifecycleRejectsInvalidTransition(t *testing.T) {
	s := testService(t)
	record := domain.NewRecord("record-b", "batch-b", "活动", "内容", "teacher")
	if err := s.CreateRecord(context.Background(), record); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Publish(context.Background(), record.ID, "reviewer"); err == nil {
		t.Fatal("expected transition error")
	}
}
