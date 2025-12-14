package dns

import (
	"testing"

	"dns-server/internal/domain/model"
)

func TestParse_RequestWithQuestion(t *testing.T) {
	msg := model.Message{
		Header: model.Header{
			ID:     10,
			QR:     false,
			Opcode: 0,
			RD:     true,
			RCode:  0,
		},
		Questions: []model.Question{
			{
				Name:  "example.com",
				Type:  1,
				Class: 1,
			},
		},
	}

	buf := Write(msg)
	parsed := Parse(buf)

	if parsed.Header.ID != msg.Header.ID {
		t.Fatalf("ID mismatch")
	}

	if parsed.Header.QR {
		t.Fatalf("QR should be false")
	}

	if len(parsed.Questions) != 1 {
		t.Fatalf("expected 1 question, got %d", len(parsed.Questions))
	}

	q := parsed.Questions[0]
	if q.Name != "example.com" {
		t.Fatalf("unexpected question name: %s", q.Name)
	}
}

func TestParse_ResponseWithAnswer(t *testing.T) {
	msg := model.Message{
		Header: model.Header{
			ID:     20,
			QR:     true,
			Opcode: 0,
			RD:     true,
			RCode:  0,
		},
		Questions: []model.Question{
			{
				Name:  "example.com",
				Type:  1,
				Class: 1,
			},
		},
		Answers: []model.Answer{
			{
				Name:  "example.com",
				Type:  1,
				Class: 1,
				TTL:   30,
				Data:  []byte{127, 0, 0, 1},
			},
		},
	}

	buf := Write(msg)
	parsed := Parse(buf)

	if !parsed.Header.QR {
		t.Fatalf("QR flag not set")
	}

	if len(parsed.Answers) != 1 {
		t.Fatalf("expected 1 answer, got %d", len(parsed.Answers))
	}

	a := parsed.Answers[0]
	if a.Name != "example.com" {
		t.Fatalf("unexpected answer name: %s", a.Name)
	}

	if len(a.Data) != 4 {
		t.Fatalf("unexpected rdata length")
	}
}
