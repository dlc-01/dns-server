package dns

import (
	"testing"

	"dns-server/internal/domain/model"
)

func TestWrite_ResponseWithQuestionAndAnswer(t *testing.T) {
	msg := model.Message{
		Header: model.Header{
			ID:     123,
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
				TTL:   60,
				Data:  []byte{127, 0, 0, 1},
			},
		},
	}

	buf := Write(msg)

	if len(buf) < 12 {
		t.Fatalf("buffer too short: %d", len(buf))
	}

	resp := Parse(buf)

	if resp.Header.ID != msg.Header.ID {
		t.Fatalf("ID mismatch: expected %d, got %d", msg.Header.ID, resp.Header.ID)
	}

	if !resp.Header.QR {
		t.Fatalf("QR flag not set")
	}

	if resp.Header.RCode != 0 {
		t.Fatalf("unexpected RCode: %d", resp.Header.RCode)
	}

	if len(resp.Questions) != 1 {
		t.Fatalf("expected 1 question, got %d", len(resp.Questions))
	}

	if resp.Questions[0].Name != "example.com" {
		t.Fatalf("unexpected question name: %s", resp.Questions[0].Name)
	}

	if len(resp.Answers) != 1 {
		t.Fatalf("expected 1 answer, got %d", len(resp.Answers))
	}

	ans := resp.Answers[0]
	if ans.Name != "example.com" {
		t.Fatalf("unexpected answer name: %s", ans.Name)
	}

	if len(ans.Data) != 4 {
		t.Fatalf("unexpected RDATA length: %d", len(ans.Data))
	}
}
