package main

import (
	"dns-server/internal/infrastructure/config"
	"dns-server/internal/infrastructure/resolver/upstream"
	"dns-server/internal/infrastructure/transport/udp"
	"dns-server/internal/usecase/service"
)

func main() {
	cfg := config.Load()

	resolver := upstream.NewUDPResolver(cfg.Upstream)
	handler := service.NewDNSService(resolver)
	server := udp.NewServer(handler)

	server.Listen(cfg.Listen)
}
