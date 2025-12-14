package service

import (
	"dns-server/internal/domain/model"
	"dns-server/internal/usecase/ports"
)

type DNSService struct {
	resolver ports.Resolver
}

func NewDNSService(r ports.Resolver) *DNSService {
	return &DNSService{resolver: r}
}

func (s *DNSService) Handle(req model.Message) (model.Message, error) {
	resp := model.Message{
		Header: model.Header{
			ID:     req.Header.ID,
			QR:     true,
			Opcode: req.Header.Opcode,
			RD:     req.Header.RD,
			RCode:  0,
		},
		Questions: req.Questions,
	}

	for _, q := range req.Questions {
		ans, err := s.resolver.Resolve(q)
		if err != nil {
			resp.Header.RCode = 2
			continue
		}
		resp.Answers = append(resp.Answers, ans...)
	}

	return resp, nil
}
