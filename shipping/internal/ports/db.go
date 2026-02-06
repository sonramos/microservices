package ports

import "github.com/sonramos/microservices/shipping/internal/application/core/domain"

type DBPort interface {
	Get(id string) (domain.Shipping, error)
	Save(*domain.Shipping) error
}
