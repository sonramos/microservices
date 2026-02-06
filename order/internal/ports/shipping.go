package ports

import "github.com/sonramos/microservices/order/internal/application/core/domain"

type ShippingPort interface {
	Ship(order *domain.Order) error
}
