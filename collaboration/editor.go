package collaboration

import (
	"context"
	"fmt"

	"example.com/nursery-cms/domain"
	"example.com/nursery-cms/service"
)

type Editor struct{ Service *service.Service }

func New(s *service.Service) Editor { return Editor{Service: s} }

func (e Editor) EditTitle(ctx context.Context, recordID, title, actor string, version int) (domain.Record, error) {
	record, err := e.Service.UpdateTitle(ctx, recordID, title, version)
	if err != nil {
		return record, err
	}
	if actor == "" {
		return record, fmt.Errorf("actor is required")
	}
	return record, nil
}

func (e Editor) EditContent(ctx context.Context, recordID, content, actor string, version int) (domain.Record, error) {
	record, err := e.Service.UpdateContent(ctx, recordID, content, version)
	if err != nil {
		return record, err
	}
	if actor == "" {
		return record, fmt.Errorf("actor is required")
	}
	return record, nil
}

func (e Editor) Publish(ctx context.Context, recordID, actor string) (domain.Record, error) {
	if actor == "" {
		return domain.Record{}, fmt.Errorf("actor is required")
	}
	return e.Service.Publish(ctx, recordID, actor)
}

func (e Editor) CanCollaborate(record domain.Record) bool {
	return domain.CanEdit(record.Status)
}
