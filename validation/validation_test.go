package validation

import (
	"context"
	"testing"

	"example.com/nursery-cms/domain"
)

func TestValidatorClassifiesRecords(t *testing.T) {
	v := New()
	valid := domain.ImportedRecord{ExternalID: "1", BatchID: "batch-v", Title: "标题", Content: "内容", Owner: "owner"}
	if !v.Valid(valid) {
		t.Fatal("valid record rejected")
	}
	invalid := valid
	invalid.Title = ""
	issues := v.ValidateImported(invalid)
	if len(issues) != 1 || issues[0].Field != "title" {
		t.Fatalf("issues=%v", issues)
	}
	if v.Explain(invalid) != "title:required" {
		t.Fatal(v.Explain(invalid))
	}
}

func TestBatchValidatorHonorsCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	b := NewBatchValidator()
	_, _, err := b.Validate(ctx, []domain.ImportedRecord{{ExternalID: "1", BatchID: "batch-v", Title: "标题", Content: "内容", Owner: "owner"}})
	if err == nil {
		t.Fatal("expected cancellation")
	}
}
