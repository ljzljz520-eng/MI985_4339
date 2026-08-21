package store

import (
	"testing"

	"example.com/nursery-cms/domain"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := t.TempDir() + "/reopen.db"
	db, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	record := domain.NewRecord("record-reopen", "batch-reopen", "重开课程", "内容", "teacher")
	if err := db.SaveRecord(record); err != nil {
		t.Fatal(err)
	}
	if err := db.SaveAuditEvent(domain.AuditEvent{ID: "event-reopen", BatchID: "batch-reopen", RecordID: record.ID, Action: "created", Actor: "teacher", Detail: "saved", CreatedAt: "fixed"}); err != nil {
		t.Fatal(err)
	}
	if err := db.SaveWorkflow(domain.NewWorkflow("workflow-reopen", "batch-reopen", "import", "teacher")); err != nil {
		t.Fatal(err)
	}
	if err := db.SaveAttachment(domain.NewAttachment("attachment-reopen", record.ID, "plan.pdf", "application/pdf", 12, "d000012")); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	db, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := db.GetRecord(record.ID); err != nil {
		t.Fatal(err)
	}
	events, err := db.ListAuditEvents("batch-reopen")
	if err != nil || len(events) != 1 {
		t.Fatalf("events=%d err=%v", len(events), err)
	}
	if _, err := db.GetWorkflow("workflow-reopen"); err != nil {
		t.Fatal(err)
	}
	attachments, err := db.ListAttachments(record.ID)
	if err != nil || len(attachments) != 1 {
		t.Fatalf("attachments=%d err=%v", len(attachments), err)
	}
}
