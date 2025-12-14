package ports

import "dns-server/internal/domain/model"

type Resolver interface {
	Resolve(model.Question) ([]model.Answer, error)
}
