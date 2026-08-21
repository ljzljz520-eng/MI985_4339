package service

import (
	"context"
	"testing"

	"example.com/nursery-cms/domain"
)

func TestWorkflowImportReport(t *testing.T) {
	s := testService(t)
	items := []domain.ImportedRecord{
		{ExternalID: "1", BatchID: "batch-import", Title: "语言活动", Content: "讲故事", Owner: "teacher"},
		{ExternalID: "2", BatchID: "batch-import", Title: "", Content: "缺标题", Owner: "teacher"},
	}
	result, err := s.ImportBatch(context.Background(), "batch-import", items)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Accepted) != 1 || len(result.Rejected) != 1 {
		t.Fatalf("accepted=%d rejected=%d", len(result.Accepted), len(result.Rejected))
	}
	state, err := s.BatchState("batch-import")
	if err != nil {
		t.Fatal(err)
	}
	if state != "completed" {
		t.Fatalf("state=%s", state)
	}
	reloaded, err := s.Store.GetRecord(result.Accepted[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	if reloaded.BatchID != "batch-import" {
		t.Fatal(reloaded.BatchID)
	}
}

func TestImportRejectsInvalidBatch(t *testing.T) {
	s := testService(t)
	_, err := s.ImportBatch(context.Background(), "?", nil)
	if err == nil {
		t.Fatal("expected invalid batch")
	}
}
