package api

import (
	"github.com/sonramos/microservices/shipping/internal/application/core/domain"
	"github.com/sonramos/microservices/shipping/internal/ports"
)

type Application struct {
	db ports.DBPort
}

func NewApplication(db ports.DBPort) *Application {
	return &Application{db: db}
}

func (a Application) CreateShipping(shipping domain.Shipping) (domain.Shipping, error) {
	err := a.db.Save(&shipping)
	if err != nil {
		return domain.Shipping{}, err
	}
	return shipping, nil
}
