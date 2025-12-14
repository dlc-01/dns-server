package upstream

import (
	"net"
	"testing"
	"time"

	"dns-server/internal/domain/model"
	"dns-server/internal/infrastructure/codec/dns"
)

func TestUDPResolver_Resolve(t *testing.T) {
	addr := "127.0.0.1:20536"

	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		t.Fatalf("resolve addr error: %v", err)
	}

	serverConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		t.Fatalf("listen udp error: %v", err)
	}
	defer serverConn.Close()

	go func() {
		buf := make([]byte, 512)
		n, client, _ := serverConn.ReadFromUDP(buf)

		req := dns.Parse(buf[:n])

		resp := model.Message{
			Header: model.Header{
				ID:     req.Header.ID,
				QR:     true,
				Opcode: req.Header.Opcode,
				RD:     req.Header.RD,
				RCode:  0,
			},
			Questions: req.Questions,
			Answers: []model.Answer{
				{
					Name:  req.Questions[0].Name,
					Type:  req.Questions[0].Type,
					Class: req.Questions[0].Class,
					TTL:   60,
					Data:  []byte{127, 0, 0, 1},
				},
			},
		}

		serverConn.WriteToUDP(dns.Write(resp), client)
	}()

	time.Sleep(50 * time.Millisecond)

	resolver := NewUDPResolver(addr)

	answers, err := resolver.Resolve(model.Question{
		Name:  "example.com",
		Type:  1,
		Class: 1,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(answers) != 1 {
		t.Fatalf("expected 1 answer, got %d", len(answers))
	}

	if answers[0].Name != "example.com" {
		t.Fatalf("unexpected answer name: %s", answers[0].Name)
	}
}
