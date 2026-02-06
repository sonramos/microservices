package ports

import "github.com/sonramos/microservices/order/internal/application/core/domain"

type DBPort interface {
Get(id string) (domain.Order, error)
Save(*domain.Order) error
ValidateStock(productCodes []string) error
}
