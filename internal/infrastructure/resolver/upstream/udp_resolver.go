package upstream

import (
	"dns-server/internal/domain/model"
	"dns-server/internal/infrastructure/codec/dns"
	"net"
)

type UDPResolver struct {
	addr *net.UDPAddr
}

func NewUDPResolver(addr string) *UDPResolver {
	a, _ := net.ResolveUDPAddr("udp", addr)
	return &UDPResolver{addr: a}
}

func (r *UDPResolver) Resolve(q model.Question) ([]model.Answer, error) {
	conn, _ := net.DialUDP("udp", nil, r.addr)
	defer conn.Close()

	req := model.Message{
		Header:    model.Header{ID: 1, RD: true},
		Questions: []model.Question{q},
	}

	conn.Write(dns.Write(req))

	buf := make([]byte, 512)
	n, _ := conn.Read(buf)

	resp := dns.Parse(buf[:n])
	return resp.Answers, nil
}
