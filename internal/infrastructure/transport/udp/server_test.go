package udp

import (
	"net"
	"testing"
	"time"

	"dns-server/internal/domain/model"
)

type mockDNSHandler struct {
	handleFunc func(model.Message) (model.Message, error)
}

func (m *mockDNSHandler) Handle(req model.Message) (model.Message, error) {
	return m.handleFunc(req)
}

func TestUDPServer_Listen_HandleRequest(t *testing.T) {
	handler := &mockDNSHandler{
		handleFunc: func(req model.Message) (model.Message, error) {
			return model.Message{
				Header: model.Header{
					ID:     req.Header.ID,
					QR:     true,
					Opcode: req.Header.Opcode,
					RD:     req.Header.RD,
					RCode:  0,
				},
				Questions: req.Questions,
			}, nil
		},
	}

	server := NewServer(handler)

	addr := "127.0.0.1:20535"

	go func() {
		_ = server.Listen(addr)
	}()

	time.Sleep(100 * time.Millisecond)

	conn, err := net.Dial("udp", addr)
	if err != nil {
		t.Fatalf("dial error: %v", err)
	}
	defer conn.Close()

	req := []byte{
		0x00, 0x01,
		0x01, 0x00,
		0x00, 0x00,
		0x00, 0x00,
		0x00, 0x00,
		0x00, 0x00,
	}

	_, err = conn.Write(req)
	if err != nil {
		t.Fatalf("write error: %v", err)
	}

	buf := make([]byte, 512)
	conn.SetReadDeadline(time.Now().Add(1 * time.Second))

	n, err := conn.Read(buf)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}

	if n < 12 {
		t.Fatalf("response too short")
	}

	flags := uint16(buf[2])<<8 | uint16(buf[3])
	if flags&(1<<15) == 0 {
		t.Fatalf("QR flag not set")
	}
}
