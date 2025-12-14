package ports

import "dns-server/internal/domain/model"

type DNSHandler interface {
	Handle(model.Message) (model.Message, error)
}
