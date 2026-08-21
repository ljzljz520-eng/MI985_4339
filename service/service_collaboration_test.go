package service_test

import (
	"context"
	"errors"
	"testing"

	"example.com/nursery-cms/collaboration"
	"example.com/nursery-cms/domain"
	"example.com/nursery-cms/service"
	"example.com/nursery-cms/store"
)

func TestWorkflowSearchUpdatePublish(t *testing.T) {
	db, err := store.Open(t.TempDir() + "/collaboration.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	s := service.New(db)
	record := domain.NewRecord("record-c", "batch-c", "区域活动", "积木搭建", "teacher")
	if err := s.CreateRecord(context.Background(), record); err != nil {
		t.Fatal(err)
	}
	items, err := s.Query(context.Background(), service.Query{BatchID: "batch-c", Text: "积木"})
	if err != nil || len(items) != 1 {
		t.Fatalf("query=%v len=%d", err, len(items))
	}
	editor := collaboration.New(s)
	updated, err := editor.EditContent(context.Background(), record.ID, "积木搭建与合作", "teacher", 1)
	if err != nil {
		t.Fatal(err)
	}
	if updated.Version != 2 {
		t.Fatalf("version=%d", updated.Version)
	}
	if _, err := editor.EditContent(context.Background(), record.ID, "过期内容", "other", 1); !errors.Is(err, domain.ErrConflict) {
		t.Fatalf("err=%v", err)
	}
	if _, err := editor.Publish(context.Background(), record.ID, "teacher"); err == nil {
		t.Fatal("draft cannot publish")
	}
}
