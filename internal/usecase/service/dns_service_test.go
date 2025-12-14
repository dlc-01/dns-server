package service

import (
	"errors"
	"testing"

	"dns-server/internal/domain/model"
)

type mockResolver struct {
	resolveFunc func(model.Question) ([]model.Answer, error)
}

func (m *mockResolver) Resolve(q model.Question) ([]model.Answer, error) {
	return m.resolveFunc(q)
}

func TestDNSService_Handle_Success(t *testing.T) {
	resolver := &mockResolver{
		resolveFunc: func(q model.Question) ([]model.Answer, error) {
			return []model.Answer{
				{
					Name:  q.Name,
					Type:  q.Type,
					Class: q.Class,
					TTL:   60,
					Data:  []byte{127, 0, 0, 1},
				},
			}, nil
		},
	}

	service := NewDNSService(resolver)

	req := model.Message{
		Header: model.Header{
			ID:     42,
			QR:     false,
			Opcode: 0,
			RD:     true,
		},
		Questions: []model.Question{
			{
				Name:  "example.com",
				Type:  1,
				Class: 1,
			},
		},
	}

	resp, err := service.Handle(req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Header.ID != req.Header.ID {
		t.Errorf("expected ID %d, got %d", req.Header.ID, resp.Header.ID)
	}

	if !resp.Header.QR {
		t.Errorf("expected QR=true")
	}

	if resp.Header.RCode != 0 {
		t.Errorf("expected RCode=0, got %d", resp.Header.RCode)
	}

	if len(resp.Questions) != 1 {
		t.Fatalf("expected 1 question, got %d", len(resp.Questions))
	}

	if len(resp.Answers) != 1 {
		t.Fatalf("expected 1 answer, got %d", len(resp.Answers))
	}

	ans := resp.Answers[0]
	if ans.Name != "example.com" {
		t.Errorf("unexpected answer name: %s", ans.Name)
	}
}

func TestDNSService_Handle_ResolverError(t *testing.T) {
	resolver := &mockResolver{
		resolveFunc: func(q model.Question) ([]model.Answer, error) {
			return nil, errors.New("resolver failed")
		},
	}

	service := NewDNSService(resolver)

	req := model.Message{
		Header: model.Header{
			ID:     1,
			Opcode: 0,
			RD:     true,
		},
		Questions: []model.Question{
			{
				Name:  "bad.example",
				Type:  1,
				Class: 1,
			},
		},
	}

	resp, err := service.Handle(req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Header.RCode != 2 {
		t.Errorf("expected RCode=2 (SERVFAIL), got %d", resp.Header.RCode)
	}

	if len(resp.Answers) != 0 {
		t.Errorf("expected no answers, got %d", len(resp.Answers))
	}
}
