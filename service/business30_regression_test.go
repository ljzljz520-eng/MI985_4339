package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"example.com/nursery-cms/domain"
)

type stagedCancelContext struct {
	calls  int
	closed chan struct{}
}

func (c *stagedCancelContext) Deadline() (time.Time, bool) { return time.Time{}, false }
func (c *stagedCancelContext) Done() <-chan struct{} {
	c.calls++
	if c.calls > 2 {
		select {
		case <-c.closed:
		default:
			close(c.closed)
		}
	}
	return c.closed
}
func (c *stagedCancelContext) Err() error {
	if c.calls > 2 {
		return context.Canceled
	}
	return nil
}
func (c *stagedCancelContext) Value(any) any { return nil }

func TestBusiness30Regression(t *testing.T) {
	s := testService(t)
	ctx := &stagedCancelContext{closed: make(chan struct{})}
	items := []domain.ImportedRecord{{ExternalID: "one", BatchID: "985-30", Title: "交接课程", Content: "真实结果", Owner: "teacher"}}
	result, err := s.ImportBatch(ctx, "985-30", items)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected cancellation, got err=%v result=%+v", err, result)
	}
	if !result.Cancelled {
		t.Fatalf("expected cancelled batch, got %+v", result)
	}
	if len(result.Accepted) != 0 {
		t.Fatalf("cancelled batch accepted=%d", len(result.Accepted))
	}
	if _, err := s.Store.GetRecord(domain.RecordID("985-30", "one")); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("cancelled record persisted: %v", err)
	}
	nextItems := []domain.ImportedRecord{{ExternalID: "one", BatchID: "985-31", Title: "交接课程", Content: "真实结果", Owner: "teacher"}}
	next, err := s.ImportBatch(context.Background(), "985-31", nextItems)
	if err != nil {
		t.Fatal(err)
	}
	if next.Cancelled || len(next.Accepted) != 1 {
		t.Fatalf("unexpected next batch result: %+v", next)
	}
}
