package udp

import (
	"dns-server/internal/infrastructure/codec/dns"
	"dns-server/internal/usecase/ports"
	"net"
)

type Server struct {
	handler ports.DNSHandler
}

func NewServer(h ports.DNSHandler) *Server {
	return &Server{handler: h}
}

func (s *Server) Listen(addr string) error {
	udpAddr, _ := net.ResolveUDPAddr("udp", addr)
	conn, _ := net.ListenUDP("udp", udpAddr)
	defer conn.Close()

	buf := make([]byte, 512)

	for {
		n, src, _ := conn.ReadFromUDP(buf)
		req := dns.Parse(buf[:n])
		resp, _ := s.handler.Handle(req)
		conn.WriteToUDP(dns.Write(resp), src)
	}
}
